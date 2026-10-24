package server

import "github.com/gin-gonic/gin"

// fail is the only place errors become responses: *domain.Error -> its Kind's status + Thai Msg;
// sqlite FK RESTRICT -> 409 "ลบไม่ได้ เพราะมีรายการสั่งซื้อใช้อยู่"; sql.ErrNoRows -> 404 "ไม่พบข้อมูล";
// anything else -> slog.Error + 500 "เกิดข้อผิดพลาด กรุณาลองใหม่".
func fail(c *gin.Context, err error) {
	panic("not implemented")
}

// bind decodes JSON into T or fails 400 "ข้อมูลไม่ถูกต้อง". ok=false means the response is written.
func bind[T any](c *gin.Context) (in T, ok bool) {
	panic("not implemented")
}

// idParam parses a positive int64 path param or fails 404. ok=false means the response is written.
func idParam(c *gin.Context, name string) (id int64, ok bool) {
	panic("not implemented")
}
