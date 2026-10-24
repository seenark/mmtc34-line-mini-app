package server

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"modernc.org/sqlite"
	sqlitelib "modernc.org/sqlite/lib"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/domain"
)

var statusOfKind = map[domain.Kind]int{
	domain.Invalid:   http.StatusBadRequest,
	domain.Forbidden: http.StatusForbidden,
	domain.NotFound:  http.StatusNotFound,
	domain.Conflict:  http.StatusConflict,
	domain.TooLarge:  http.StatusRequestEntityTooLarge,
}

// fail is the only place errors become responses: *domain.Error -> its Kind's status + Thai Msg;
// sqlite FK violation -> 409 "ลบไม่ได้ เพราะมีรายการสั่งซื้อใช้อยู่"; sql.ErrNoRows -> 404 "ไม่พบข้อมูล";
// anything else -> slog.Error + 500 "เกิดข้อผิดพลาด กรุณาลองใหม่".
func fail(c *gin.Context, err error) {
	var de *domain.Error
	var se *sqlite.Error
	switch {
	case errors.As(err, &de):
		abort(c, statusOfKind[de.Kind], de.Msg)
	case errors.As(err, &se) && se.Code() == sqlitelib.SQLITE_CONSTRAINT_FOREIGNKEY:
		abort(c, http.StatusConflict, "ลบไม่ได้ เพราะมีรายการสั่งซื้อใช้อยู่")
	case errors.Is(err, sql.ErrNoRows):
		abort(c, http.StatusNotFound, "ไม่พบข้อมูล")
	default:
		slog.Error("request failed", "method", c.Request.Method, "path", c.Request.URL.Path, "err", err)
		abort(c, http.StatusInternalServerError, "เกิดข้อผิดพลาด กรุณาลองใหม่")
	}
}

func abort(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, api.ErrorResponse{Error: msg})
}

const maxJSONBody = 64 << 10

// bind decodes JSON (body capped at maxJSONBody) into T or fails 400 "ข้อมูลไม่ถูกต้อง", or 413
// when the body is over the cap. ok=false means the response is written.
func bind[T any](c *gin.Context) (in T, ok bool) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxJSONBody)
	if err := c.ShouldBindJSON(&in); err != nil {
		var tooBig *http.MaxBytesError
		if errors.As(err, &tooBig) {
			fail(c, &domain.Error{Kind: domain.TooLarge, Msg: "ข้อมูลใหญ่เกินกำหนด"})
		} else {
			fail(c, &domain.Error{Kind: domain.Invalid, Msg: "ข้อมูลไม่ถูกต้อง"})
		}
		return in, false
	}
	return in, true
}

// idParam parses a positive int64 path param or fails 404. ok=false means the response is written.
func idParam(c *gin.Context, name string) (id int64, ok bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		fail(c, &domain.Error{Kind: domain.NotFound, Msg: "ไม่พบข้อมูล"})
		return 0, false
	}
	return id, true
}
