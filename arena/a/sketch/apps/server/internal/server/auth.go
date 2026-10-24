package server

import "github.com/gin-gonic/gin"

// Ticket .2.
func (s *Server) mountAuth(r routes) {
	r.public.POST("/auth/line", s.lineLogin)
	r.public.POST("/auth/dev", s.devLogin) // 404 unless DEV_LOGIN
	r.public.POST("/auth/logout", s.logout)
	r.self.GET("/me", s.me)
	r.self.POST("/me", s.onboard)
}

// lineLogin: Verify(idToken) -> 401 "ยืนยันตัวตนกับ LINE ไม่สำเร็จ" on failure; promote a
// bootstrap Admin's existing row (approved + admin); setSession; return api.Me.
func (s *Server) lineLogin(c *gin.Context) { panic("not implemented") }

// devLogin: same as lineLogin with the userId taken on trust. Route returns 404 when !DevLogin.
func (s *Server) devLogin(c *gin.Context) { panic("not implemented") }

// login is the shared tail of lineLogin and devLogin.
func (s *Server) login(c *gin.Context, lineUserID string) { panic("not implemented") }

func (s *Server) logout(c *gin.Context) { panic("not implemented") }

// me returns api.Me; Member null means "show onboarding".
func (s *Server) me(c *gin.Context) { panic("not implemented") }

// onboard creates the Applicant (or approved Admin for a bootstrap id) from ProfileInput
// (RequireIdentity). 409 if this LINE userId already has a Member.
func (s *Server) onboard(c *gin.Context) { panic("not implemented") }
