package server

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/db"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

const (
	cookieName = "mmtc34_session"
	sessionTTL = 30 * 24 * time.Hour
)

// The cookie names the LINE userId, not the Member id: linking an Applicant to an Unlinked
// Member moves the userId, and deleting an account clears it, so both take effect on the
// next request without touching sessions.
//
// Value: base64url(lineUserID) "." expUnix "." base64url(HMAC-SHA256(secret, lineUserID "." expUnix)).
func signSession(secret []byte, lineUserID string, exp time.Time) string {
	panic("not implemented")
}

// parseSession checks the MAC with hmac.Equal and expiry against now.
func parseSession(secret []byte, value string, now time.Time) (lineUserID string, ok bool) {
	panic("not implemented")
}

// setSession writes the cookie: HttpOnly, SameSite=Lax, Path=BasePath, Max-Age 30 days,
// Secure unless DevLogin.
func (s *Server) setSession(c *gin.Context, lineUserID string) {
	panic("not implemented")
}

func (s *Server) clearSession(c *gin.Context) {
	panic("not implemented")
}

// viewer is the per-request identity: LINE userId from the cookie, Member reloaded from the DB.
type viewer struct {
	LineUserID string
	Member     *db.Member // nil before onboarding or after deletion
}

// requireSession: valid cookie or 401 "กรุณาเข้าสู่ระบบ". Loads the Member by line_user_id
// every request and stores the viewer in the Gin context.
func (s *Server) requireSession(c *gin.Context) { panic("not implemented") }

// requireApproved: Member exists, approved, not deleted; else 403 "รอแอดมินอนุมัติ".
func (s *Server) requireApproved(c *gin.Context) { panic("not implemented") }

// requireAdmin: role admin; else 403 "สำหรับแอดมินเท่านั้น".
func (s *Server) requireAdmin(c *gin.Context) { panic("not implemented") }

func viewerOf(c *gin.Context) viewer { panic("not implemented") }

// actorOf is the domain.Actor for the approved Member on this request.
func actorOf(c *gin.Context) domain.Actor { panic("not implemented") }
