package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/server"
)

func adminLogin(password string) api.AdminLoginInput {
	return api.AdminLoginInput{Username: "admin", Password: password}
}

func TestSuperAdminLoginSignsInAsApprovedAdmin(t *testing.T) {
	h := New(t)

	m := h.SuperAdmin().Me().Member

	Equal(t, m.Status, api.MemberApproved)
	Equal(t, m.Role, api.RoleAdmin)
	Equal(t, m.SuperAdmin, true)
	Equal(t, m.Nickname, "ผู้ดูแลระบบ")
	h.SuperAdmin().GET("/api/admin/members").Status(200)
}

func TestSuperAdminLoginTwiceKeepsOneRow(t *testing.T) {
	h := New(t)

	h.SuperAdmin()
	boss := h.SuperAdmin()

	supers := 0
	for _, m := range As[[]api.Member](boss.GET("/api/admin/members").Status(200)) {
		if m.SuperAdmin {
			supers++
		}
	}
	Equal(t, supers, 1)
}

func TestSuperAdminWrongCredentialsAre401(t *testing.T) {
	h := New(t)
	c := h.Anonymous()

	Equal(t, c.POST("/api/auth/admin", adminLogin("wrong-password-123")).Status(401).Error(), "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง")
	c.POST("/api/auth/admin", api.AdminLoginInput{Username: "root", Password: superAdminPassword}).Status(401)
	c.GET("/api/me").Status(401)
}

func TestSuperAdminLoginIsThrottledAfterFiveFailures(t *testing.T) {
	h := New(t)
	c := h.Anonymous()

	for range 5 {
		c.POST("/api/auth/admin", adminLogin("wrong-password-123")).Status(401)
	}

	Equal(t, c.POST("/api/auth/admin", adminLogin(superAdminPassword)).Status(429).Error(), "ลองใหม่อีกครั้งในอีกสักครู่")
	h.Advance(61 * time.Second)
	c.POST("/api/auth/admin", adminLogin(superAdminPassword)).Status(204)
}

func TestSuperAdminLoginIsDisabledWithoutPassword(t *testing.T) {
	h := New(t, func(c *server.Config) { c.SuperAdminPassword = "" })
	c := h.Anonymous()

	c.POST("/api/auth/admin", adminLogin("")).Status(404)
	Equal(t, As[api.Config](c.GET("/api/config").Status(200)).AdminLogin, false)
}

func TestReservedSuperAdminIDIsRejectedByDevLogin(t *testing.T) {
	h := New(t)

	h.Anonymous().POST("/api/auth/dev", api.DevLoginInput{LineUserID: "superadmin"}).Status(400)
}

func TestReservedSuperAdminIDIsRejectedByLineLogin(t *testing.T) {
	h := New(t, func(c *server.Config) { c.DevLogin = false })

	h.Anonymous().POST("/api/auth/line", api.LineLoginInput{IDToken: "ok:superadmin"}).Status(401)
}

func TestSuperAdminApprovesAndPromotesAnApplicant(t *testing.T) {
	h := New(t)
	boss, nok := h.SuperAdmin(), h.Applicant("nok")
	id := nok.Me().Member.ID

	boss.POST(fmt.Sprintf("/api/admin/members/%d/approve", id), nil).Status(200)
	boss.PUT(memberPath(id)+"/role", api.RoleInput{Role: api.RoleAdmin}).Status(200)

	nok.GET("/api/admin/members").Status(200)
}

func TestPromotedAdminCannotDemoteSuperAdmin(t *testing.T) {
	h := New(t)
	admin := h.Admin()
	superID := h.SuperAdmin().Me().Member.ID

	Equal(t, admin.PUT(memberPath(superID)+"/role", api.RoleInput{Role: api.RoleMember}).Status(409).Error(), "เปลี่ยนบทบาทผู้ดูแลระบบหลักไม่ได้")
}

func TestSuperAdminCannotDemoteSelfOrDeleteAccount(t *testing.T) {
	h := New(t)
	boss := h.SuperAdmin()

	boss.PUT(memberPath(boss.Me().Member.ID)+"/role", api.RoleInput{Role: api.RoleMember}).Status(409)
	boss.DELETE("/api/me").Status(409)
	boss.GET("/api/admin/members").Status(200)
}
