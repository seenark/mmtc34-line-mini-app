package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

const (
	cookieName = "mmtc34_session"
	sessionTTL = 30 * 24 * time.Hour
	viewerKey  = "viewer"
)

var b64 = base64.RawURLEncoding

// The cookie names the LINE userId, not the Member id: deleting an account clears the userId,
// which takes effect on the next request without touching sessions.
//
// Value: base64url(lineUserID) "." expUnix "." base64url(HMAC-SHA256(secret, lineUserID "." expUnix)).
func signSession(secret []byte, lineUserID string, exp time.Time) string {
	id, unix := b64.EncodeToString([]byte(lineUserID)), strconv.FormatInt(exp.Unix(), 10)
	return id + "." + unix + "." + b64.EncodeToString(mac(secret, lineUserID, unix))
}

func mac(secret []byte, lineUserID, unix string) []byte {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(lineUserID + "." + unix))
	return h.Sum(nil)
}

// parseSession checks the MAC with hmac.Equal and expiry against now.
func parseSession(secret []byte, value string, now time.Time) (lineUserID string, ok bool) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", false
	}
	id, err := b64.DecodeString(parts[0])
	if err != nil || len(id) == 0 {
		return "", false
	}
	got, err := b64.DecodeString(parts[2])
	if err != nil || !hmac.Equal(got, mac(secret, string(id), parts[1])) {
		return "", false
	}
	unix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || !now.Before(time.Unix(unix, 0)) {
		return "", false
	}
	return string(id), true
}

// cookiePath is the Path attribute: the prefix, or "/" at the root (an empty Path would make
// browsers scope the cookie to the login endpoint's directory).
func (s *Server) cookiePath() string {
	if s.cfg.BasePath == "" {
		return "/"
	}
	return s.cfg.BasePath
}

// setSession writes the cookie: HttpOnly, SameSite=Lax, Path=BasePath (or "/"), Max-Age 30 days,
// Secure unless DevLogin.
func (s *Server) setSession(c *gin.Context, lineUserID string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    signSession(s.cfg.SessionSecret, lineUserID, s.now().Add(sessionTTL)),
		Path:     s.cookiePath(),
		MaxAge:   int(sessionTTL / time.Second),
		HttpOnly: true,
		Secure:   !s.cfg.DevLogin,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *Server) clearSession(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Path:     s.cookiePath(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   !s.cfg.DevLogin,
		SameSite: http.SameSiteLaxMode,
	})
}

// viewer is the per-request identity: LINE userId from the cookie, Member reloaded from the DB.
type viewer struct {
	LineUserID string
	Member     *db.Member // nil before onboarding or after deletion
}

// requireSession: valid cookie or 401 "กรุณาเข้าสู่ระบบ". Loads the Member by line_user_id
// every request and stores the viewer in the Gin context.
func (s *Server) requireSession(c *gin.Context) {
	value, err := c.Cookie(cookieName)
	lineUserID, ok := parseSession(s.cfg.SessionSecret, value, s.now())
	if err != nil || !ok {
		abort(c, http.StatusUnauthorized, "กรุณาเข้าสู่ระบบ")
		return
	}
	v := viewer{LineUserID: lineUserID}
	m, err := s.q.GetMemberByLineUserID(c, &lineUserID)
	switch {
	case err == nil:
		v.Member = &m
	case !errors.Is(err, sql.ErrNoRows):
		fail(c, err)
		return
	}
	c.Set(viewerKey, v)
}

// requireMember: a Member exists and is not deleted; else 403 "กรุณาลงทะเบียนก่อน".
func (s *Server) requireMember(c *gin.Context) {
	if !isActive(viewerOf(c).Member) {
		abort(c, http.StatusForbidden, "กรุณาลงทะเบียนก่อน")
	}
}

// requireAdmin: active Member with role admin; else 403 "สำหรับแอดมินเท่านั้น".
func (s *Server) requireAdmin(c *gin.Context) {
	if m := viewerOf(c).Member; !isActive(m) || m.Role != string(api.RoleAdmin) {
		abort(c, http.StatusForbidden, "สำหรับแอดมินเท่านั้น")
	}
}

func isActive(m *db.Member) bool {
	return m != nil && m.DeletedAt == nil
}

func viewerOf(c *gin.Context) viewer { return c.MustGet(viewerKey).(viewer) }

// actorOf is the domain.Actor for the Member on this request. Only valid behind requireMember.
func actorOf(c *gin.Context) domain.Actor {
	m := viewerOf(c).Member
	return domain.Actor{MemberID: m.ID, Admin: m.Role == string(api.RoleAdmin)}
}

// rfc3339 is the one stored-time format (UTC), always fed from s.now().
func rfc3339(t time.Time) string { return t.UTC().Format(time.RFC3339) }
