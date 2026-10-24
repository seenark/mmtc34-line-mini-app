package server_test

import (
	"fmt"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

func campaignOrdersPath(campaignID int64) string {
	return fmt.Sprintf("/api/admin/campaigns/%d/orders", campaignID)
}
func cashPath(orderID int64) string { return fmt.Sprintf("/api/admin/orders/%d/cash", orderID) }
func adjustPath(orderID int64) string {
	return fmt.Sprintf("/api/admin/orders/%d/adjustments", orderID)
}

// CampaignOrders lists a Campaign's Orders as this Admin; status "" means all.
func (c *Client) CampaignOrders(campaignID int64, status string) []api.Order {
	path := campaignOrdersPath(campaignID)
	if status != "" {
		path += "?status=" + status
	}
	return As[[]api.Order](c.GET(path).Status(200))
}

// Cash records a cash Payment as this Admin and asserts 200.
func (c *Client) Cash(orderID, amount int64) api.Order {
	return As[api.Order](c.POST(cashPath(orderID), api.AmountInput{Amount: amount}).Status(200))
}

// Adjust adds an Adjustment as this Admin and asserts 200.
func (c *Client) Adjust(orderID, amount int64, note string) api.Order {
	return As[api.Order](c.POST(adjustPath(orderID), api.AdjustmentInput{Amount: amount, Note: note}).Status(200))
}
