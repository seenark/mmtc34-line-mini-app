package server_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func uploadJPEG(t *testing.T, admin *Client) string {
	t.Helper()
	return As[api.ImageUpload](admin.POST("/api/admin/images", JPEG()).Status(201)).Name
}

func TestUploadedImageIsReadableByMembersWithSafeHeaders(t *testing.T) {
	h := New(t)
	name := uploadJPEG(t, h.Admin())

	r := h.Member("bob").GET("/api/images/" + name).Status(200)

	Equal(t, r.Header.Get("Cache-Control"), "private, max-age=31536000, immutable")
	Equal(t, r.Header.Get("X-Content-Type-Options"), "nosniff")
	Equal(t, r.Header.Get("Content-Type"), "image/jpeg")
	Equal(t, bytes.Equal(r.Body, JPEG().Bytes), true)
	Equal(t, len(name), 32+len(".jpg"))
	h.Anonymous().GET("/api/images/" + name).Status(401)
	h.Visitor("nok").GET("/api/images/" + name).Status(403)
}

func TestEachUploadGetsItsOwnName(t *testing.T) {
	h := New(t)
	admin := h.Admin()

	Equal(t, uploadJPEG(t, admin) != uploadJPEG(t, admin), true)
}

func TestUploadAcceptsPNGAndWebPBySniffedContent(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	png := Upload{Filename: "x.bin", Bytes: append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 32)...)}
	webp := Upload{Filename: "x.bin", Bytes: append([]byte("RIFF\x00\x00\x00\x00WEBPVP8 "), make([]byte, 32)...)}

	Equal(t, As[api.ImageUpload](admin.POST("/api/admin/images", png).Status(201)).Name[32:], ".png")
	Equal(t, As[api.ImageUpload](admin.POST("/api/admin/images", webp).Status(201)).Name[32:], ".webp")
}

func TestUploadRejectsNonImageEvenWithImageFilename(t *testing.T) {
	h := New(t)

	r := h.Admin().POST("/api/admin/images", Upload{Filename: "photo.png", Bytes: []byte("hello, not a picture")}).Status(400)

	Equal(t, r.Error(), "รองรับเฉพาะไฟล์ JPEG, PNG หรือ WebP")
}

func TestUploadWithoutFileIs400(t *testing.T) {
	h := New(t)

	h.Admin().POST("/api/admin/images", nil).Status(400)
}

func TestUploadOverFiveMegabytesIs413(t *testing.T) {
	h := New(t)
	big := append(JPEG().Bytes, make([]byte, 5<<20)...)

	r := h.Admin().POST("/api/admin/images", Upload{Filename: "big.jpg", Bytes: big}).Status(413)

	Equal(t, r.Error(), "ไฟล์ใหญ่เกิน 5 MB")
}

func TestImageWithBadOrUnknownNameIs404(t *testing.T) {
	h := New(t)
	bob := h.Member("bob")

	for _, name := range []string{"x.jpg", "..%2F..%2Fapp.db", "0123456789abcdef0123456789abcdef.gif", "0123456789abcdef0123456789abcdef.jpg"} {
		bob.GET("/api/images/" + name).Status(404)
	}
}

func TestSavingWithUnknownImageNameIs400(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	j := h.OpenJersey()
	unknown := "0123456789abcdef0123456789abcdef.jpg"

	in := draft("ชื่อ")
	in.CoverImage = unknown
	Equal(t, admin.POST("/api/admin/campaigns", in).Status(400).Error(), "ไม่พบรูปภาพ กรุณาอัปโหลดใหม่")
	in.CoverImage, in.QRImage = "", "../etc/passwd"
	admin.PUT(fmt.Sprintf("/api/admin/campaigns/%d", j.Campaign.ID), in).Status(400)
	admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", j.Campaign.ID), api.OptionInput{Name: "x", Image: unknown}).Status(400)
	admin.PUT(fmt.Sprintf("/api/admin/options/%d", j.Option.ID), api.OptionInput{Name: "x", Image: unknown}).Status(400)
}

func TestSavingWithUploadedImageKeepsTheName(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	name := uploadJPEG(t, admin)
	in := draft("ชื่อ")
	in.CoverImage, in.QRImage = name, name

	camp := As[api.Campaign](admin.POST("/api/admin/campaigns", in).Status(201))
	opt := As[api.Option](admin.POST(fmt.Sprintf("/api/admin/campaigns/%d/options", camp.ID), api.OptionInput{Name: "x", Image: name}).Status(201))

	Equal(t, camp.CoverImage, name)
	Equal(t, camp.QRImage, name)
	Equal(t, opt.Image, name)
}
