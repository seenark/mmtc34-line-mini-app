package server_test

import (
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func TestConfigEndpointServesLiffIDAndLoginModes(t *testing.T) {
	h := New(t)

	cfg := As[api.Config](h.Anonymous().GET("/api/config").Status(200))

	Equal(t, cfg.LiffID, "test-liff-id")
	Equal(t, cfg.DevLogin, true)
	Equal(t, cfg.AdminLogin, true)
}

func TestUnknownAPIPathIsJSON404(t *testing.T) {
	h := New(t)

	Equal(t, h.Anonymous().GET("/api/nope").Status(404).Error(), "ไม่พบข้อมูล")
}

func TestServerHonoursNonDefaultBasePath(t *testing.T) {
	h := New(t, WithBasePath("/x"))
	c := h.Anonymous()

	c.GET("/api/config").Status(200)
	c.GETAt("/mmtc34/api/config").Status(404)
}

func TestBasePathWithoutSlashRedirects(t *testing.T) {
	h := New(t)

	r := h.Anonymous().GETAt("/mmtc34").Status(301)

	Equal(t, r.Header.Get("Location"), "/mmtc34/")
}
