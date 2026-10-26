package server_test

import (
	"os"
	"strings"
	"testing"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/server"
)

func TestConfigEndpointServesLiffIDAndLoginModes(t *testing.T) {
	h := New(t)

	cfg := As[api.Config](h.Anonymous().GET("/api/config").Status(200))

	Equal(t, cfg.LiffID, "test-liff-id")
	Equal(t, cfg.DevLogin, true)
	Equal(t, cfg.AdminLogin, true)
}

func TestConfigListsBanksInRegistryOrder(t *testing.T) {
	h := New(t)

	banks := As[api.Config](h.Anonymous().GET("/api/config").Status(200)).Banks

	Equal(t, len(banks), 16)
	Equal(t, banks[0], api.Bank{Code: "002", Name: "ธนาคารกรุงเทพ"})
	Equal(t, banks[1], api.Bank{Code: "004", Name: "ธนาคารกสิกรไทย"})
	Equal(t, banks[15], api.Bank{Code: "073", Name: "ธนาคารแลนด์ แอนด์ เฮ้าส์"})
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

func TestServerAtRootServesAPIWithoutPrefix(t *testing.T) {
	h := New(t, WithBasePath(""))
	c := h.Anonymous()

	cfg := As[api.Config](c.GET("/api/config").Status(200))

	Equal(t, cfg.LiffID, "test-liff-id")
}

func TestUnknownAPIPathAtRootIsJSON404NotIndex(t *testing.T) {
	h := New(t, WithBasePath(""))

	r := h.Anonymous().GET("/api/nope").Status(404)

	Equal(t, r.Error(), "ไม่พบข้อมูล")
}

func TestSessionCookieAtRootIsScopedToRoot(t *testing.T) {
	h := New(t, WithBasePath(""))
	c := h.Anonymous()

	r := c.POST("/api/auth/line", api.LineLoginInput{IDToken: "ok:U-nok"}).Status(204)

	if cookie := r.Header.Get("Set-Cookie"); !strings.Contains(cookie, "Path=/;") {
		t.Fatalf("Set-Cookie %q lacks Path=/", cookie)
	}
	Equal(t, c.Me().LineUserID, "U-nok")
}

func TestBasePathFromEnv(t *testing.T) {
	cases := []struct {
		name  string
		unset bool
		value string
		want  string
	}{
		{"unset keeps the default prefix", true, "", "/mmtc34"},
		{"explicit prefix", false, "/x", "/x"},
		{"prefix without leading slash", false, "x/", "/x"},
		{"slash means root", false, "/", ""},
		{"empty means root", false, "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("SESSION_SECRET", strings.Repeat("k", 32))
			t.Setenv("DEV_LOGIN", "true")
			t.Setenv("BASE_PATH", tc.value)
			if tc.unset {
				os.Unsetenv("BASE_PATH")
			}

			cfg, err := server.ConfigFromEnv()

			if err != nil {
				t.Fatal(err)
			}
			Equal(t, cfg.BasePath, tc.want)
		})
	}
}
