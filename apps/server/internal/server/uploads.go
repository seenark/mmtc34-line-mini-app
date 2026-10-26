package server

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

// fileKind is a subdirectory of DATA_DIR/files. The directory, not a table, decides who may
// read a file: images -> any approved Member; slips -> owning Member or Admin (payments.go).
type fileKind string

const (
	kindImages fileKind = "images"
	kindSlips  fileKind = "slips"

	maxFile  = 5 << 20
	maxSlips = 6 // images per Slip upload
)

// fileName is the only accepted shape of a stored file name; it rules out path traversal.
var fileName = regexp.MustCompile(`^[0-9a-f]{32}\.(jpg|png|webp)$`)

var extOf = map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"}

func (s *Server) mountUploads(r routes) {
	r.admin.POST("/images", s.uploadImage)
	r.member.GET("/images/:name", s.readImage)
}

func (s *Server) filePath(kind fileKind, name string) string {
	return filepath.Join(s.cfg.DataDir, "files", string(kind), name)
}

// imageExists reports whether name is a stored image. The empty name means "no image" and is allowed.
func (s *Server) imageExists(name string) bool {
	if name == "" {
		return true
	}
	if !fileName.MatchString(name) {
		return false
	}
	_, err := os.Stat(s.filePath(kindImages, name))
	return err == nil
}

// requireImages is the 400 every Campaign / Option save answers for an unknown image name.
func (s *Server) requireImages(names ...string) error {
	for _, n := range names {
		if !s.imageExists(n) {
			return &domain.Error{Kind: domain.Invalid, Msg: "ไม่พบรูปภาพ กรุณาอัปโหลดใหม่"}
		}
	}
	return nil
}

func (s *Server) uploadImage(c *gin.Context) {
	name, err := s.saveUpload(c, kindImages)
	if err != nil {
		fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, api.ImageUpload{Name: name})
}

func (s *Server) readImage(c *gin.Context) {
	name := c.Param("name")
	if !s.imageExists(name) {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return
	}
	path := s.filePath(kindImages, name)
	c.Header("Cache-Control", "private, max-age=31536000, immutable") // names are random, never reused
	c.Header("X-Content-Type-Options", "nosniff")
	c.File(path)
}

// saveUpload stores the one file of an image upload.
func (s *Server) saveUpload(c *gin.Context, kind fileKind) (string, error) {
	names, err := s.saveUploads(c, kind, 1)
	if err != nil {
		return "", err
	}
	return names[0], nil
}

// saveUploads stores 1..max files sent under the repeated form field "file". The body is capped
// at max files plus 1 MB of multipart overhead and each file at 5 MB; only JPEG, PNG and WebP,
// judged by sniffing the bytes rather than trusting the client's name or type, are accepted.
// On any error the files already stored by this call are deleted.
// shortcut: a file whose DB row later fails to commit is orphaned on disk; add a sweep if disk use matters.
func (s *Server) saveUploads(c *gin.Context, kind fileKind, max int) ([]string, error) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(max)*maxFile+1<<20)
	form, err := c.MultipartForm()
	var tooBig *http.MaxBytesError
	switch {
	case errors.As(err, &tooBig):
		return nil, &domain.Error{Kind: domain.TooLarge, Msg: "ไฟล์ใหญ่เกิน 5 MB"}
	case err != nil || len(form.File["file"]) == 0:
		return nil, &domain.Error{Kind: domain.Invalid, Msg: "กรุณาเลือกไฟล์"}
	case len(form.File["file"]) > max:
		return nil, &domain.Error{Kind: domain.Invalid, Msg: fmt.Sprintf("เลือกได้ไม่เกิน %d รูป", max)}
	}
	names := make([]string, 0, len(form.File["file"]))
	for _, fh := range form.File["file"] {
		name, err := s.saveFile(fh, kind)
		if err != nil {
			s.removeFiles(kind, names)
			return nil, err
		}
		names = append(names, name)
	}
	return names, nil
}

func (s *Server) removeFiles(kind fileKind, names []string) {
	for _, n := range names {
		os.Remove(s.filePath(kind, n))
	}
}

// saveFile writes a temp file in the target directory and renames it, so a reader never sees a partial file.
func (s *Server) saveFile(fh *multipart.FileHeader, kind fileKind) (string, error) {
	if fh.Size > maxFile {
		return "", &domain.Error{Kind: domain.TooLarge, Msg: "ไฟล์ใหญ่เกิน 5 MB"}
	}
	src, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	head := make([]byte, 512)
	n, err := io.ReadFull(src, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return "", err
	}
	ext, ok := extOf[http.DetectContentType(head[:n])] // Go sniffs the WebP RIFF magic itself
	if !ok {
		return "", &domain.Error{Kind: domain.Invalid, Msg: "รองรับเฉพาะไฟล์ JPEG, PNG หรือ WebP"}
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	dir := filepath.Dir(s.filePath(kind, "x"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp(dir, ".upload-*")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmp.Name()) // no-op after a successful rename
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}
	name := hex.EncodeToString(raw[:]) + ext
	if err := os.Rename(tmp.Name(), filepath.Join(dir, name)); err != nil {
		return "", err
	}
	return name, nil
}
