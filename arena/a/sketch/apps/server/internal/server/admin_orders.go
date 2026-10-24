package server

import "github.com/gin-gonic/gin"

// Ticket .7. Queries: db/queries/admin_orders.sql.
func (s *Server) mountAdminOrders(r routes) {
	r.admin.POST("/campaigns/:campaignId/orders", s.orderOnBehalf) // api.AdminOrderInput
	r.admin.PUT("/orders/:orderId", s.editOrder)                   // api.OrderInput
}

// orderOnBehalf: member must exist, be approved, not deleted;
// s.checkout(campaignID, in.MemberID, &actor.MemberID, in.Order).
func (s *Server) orderOnBehalf(c *gin.Context) { panic("not implemented") }

// editOrder: apply(EvEdit) with fx: ValidateOrder(in, options, keep = current Items by id),
// update delivery/ship snapshot, delete Items not kept, upsert the rest.
func (s *Server) editOrder(c *gin.Context) { panic("not implemented") }
