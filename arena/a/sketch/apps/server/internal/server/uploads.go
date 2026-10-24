package server

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

// fileKind is a subdirectory of DATA_DIR/files. The directory, not a table, decides who may
// read a file: images -> any approved Member; slips -> owning Member or Admin (payments.go).
type fileKind string

const (
	kindImages fileKind = "images"
	kindSlips  fileKind = "slips"
	maxUpload           = 5 << 20
)

// fileName is the only accepted shape of a stored file name; it rules out path traversal.
var fileName = regexp.MustCompile(`^[0-9a-f]{32}\.(jpg|png|webp)$`)

// Ticket .3.
func (s *Server) mountUploads(r routes) {
	r.admin.POST("/images", s.uploadImage)   // multipart field "file" -> api.ImageUpload
	r.member.GET("/images/:name", s.readImage)
}

func (s *Server) uploadImage(c *gin.Context) { panic("not implemented") }
func (s *Server) readImage(c *gin.Context)   { panic("not implemented") }

// saveUpload reads multipart field "file" through http.MaxBytesReader(5 MB) -> 413
// "ไฟล์ใหญ่เกิน 5 MB"; sniffs the first 512 bytes with http.DetectContentType and accepts only
// image/jpeg, image/png, image/webp -> else 400 "รองรับเฉพาะไฟล์ JPEG, PNG หรือ WebP";
// writes to a temp file then renames to <32 hex random>.<ext> under DATA_DIR/files/<kind>.
// shortcut: a file whose DB row later fails to commit is orphaned on disk; add a sweep if disk use matters.
func (s *Server) saveUpload(c *gin.Context, kind fileKind) (name string, err error) {
	panic("not implemented")
}

// sendFile validates name against fileName (else 404) and serves it with
// Cache-Control: private, max-age=31536000, immutable (names are random and never reused).
func (s *Server) sendFile(c *gin.Context, kind fileKind, name string) { panic("not implemented") }
