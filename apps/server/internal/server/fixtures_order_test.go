package server_test

import (
	"fmt"

	"github.com/seenark/mmtc34-line-mini-app/apps/server/internal/api"
)

// Shirt is one jersey Item with a valid print name and number.
func Shirt(j Jersey, size string, qty int64) api.ItemInput {
	return api.ItemInput{OptionID: j.Option.ID, Size: size, PrintName: "ALICE", Number: "7", Quantity: qty}
}

func Pickup(items ...api.ItemInput) api.OrderInput {
	return api.OrderInput{Delivery: api.DeliveryPickup, Items: items}
}

func Ship(items ...api.ItemInput) api.OrderInput {
	return api.OrderInput{
		Delivery:    api.DeliveryShip,
		ShipName:    "สมชาย ใจดี",
		ShipAddress: "1 ถนนสุขุมวิท กรุงเทพฯ 10110",
		ShipPhone:   "0812345678",
		Items:       items,
	}
}

// Checkout places an Order as this client's Member and asserts 201.
func (c *Client) Checkout(campaignID int64, in api.OrderInput) api.Order {
	return As[api.Order](c.POST(fmt.Sprintf("/api/campaigns/%d/orders", campaignID), in).Status(201))
}

// Order reads one Order as this client and asserts 200.
func (c *Client) Order(id int64) api.Order {
	return As[api.Order](c.GET(fmt.Sprintf("/api/orders/%d", id)).Status(200))
}
