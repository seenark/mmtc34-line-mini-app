package server_test

import (
	"fmt"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// UploadSlip uploads a real JPEG as this client's Slip for the Order and asserts 200.
func (c *Client) UploadSlip(orderID int64) api.Order {
	return As[api.Order](c.POST(fmt.Sprintf("/api/orders/%d/slip", orderID), JPEG()).Status(200))
}
