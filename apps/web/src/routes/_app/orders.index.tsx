import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { ordersQuery } from "@/api/orders";
import {
  balanceLabel,
  OrderStatusBadge,
} from "@/components/order/order-status";

export const Route = createFileRoute("/_app/orders/")({
  staticData: {
    nav: { slot: "member", label: "คำสั่งซื้อของฉัน", order: 20 },
  },
  loader: ({ context }) => context.queryClient.ensureQueryData(ordersQuery),
  component: Orders,
});

function Orders() {
  const { data: orders } = useSuspenseQuery(ordersQuery);
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">คำสั่งซื้อของฉัน</h1>
      {orders.length === 0 ? (
        <p className="text-muted-foreground text-sm">ยังไม่มีคำสั่งซื้อ</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {orders.map((o) => (
            <li key={o.id}>
              <Link to="/orders/$orderId" params={{ orderId: String(o.id) }}>
                <Card>
                  <CardContent className="flex flex-col gap-1">
                    <div className="flex items-start justify-between gap-2">
                      <span className="text-sm font-medium">
                        {o.campaignTitle}
                      </span>
                      <OrderStatusBadge status={o.status} />
                    </div>
                    <span className="text-muted-foreground">
                      คำสั่งซื้อ #{o.id}
                    </span>
                    <span>{balanceLabel(o.balance)}</span>
                  </CardContent>
                </Card>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
