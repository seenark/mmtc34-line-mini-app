import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Empty, EmptyDescription } from "@mmtc34/ui/components/empty";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { ChevronRight } from "lucide-react";

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
      <h1 className="text-primary text-2xl font-extrabold">คำสั่งซื้อของฉัน</h1>
      {orders.length === 0 ? (
        <Empty className="flex-none rounded-xl border">
          <EmptyDescription>ยังไม่มีคำสั่งซื้อ</EmptyDescription>
        </Empty>
      ) : (
        <ul className="flex flex-col gap-3">
          {orders.map((o) => (
            <li key={o.id}>
              <Link
                to="/orders/$orderId"
                params={{ orderId: String(o.id) }}
                className="group block rounded-xl"
              >
                <Card className="group-hover:bg-accent group-active:bg-accent transition-colors duration-(--dur-micro) ease-(--ease-out)">
                  <CardContent className="flex items-center gap-3">
                    <div className="flex min-w-0 flex-1 flex-col gap-1.5">
                      <div className="flex items-start justify-between gap-2">
                        <span className="text-primary min-w-0 text-base font-semibold [overflow-wrap:anywhere]">
                          {o.campaignTitle}
                        </span>
                        <OrderStatusBadge status={o.status} />
                      </div>
                      <span className="text-muted-foreground text-sm">
                        คำสั่งซื้อ #{o.id}
                      </span>
                      <span
                        className={
                          o.balance > 0
                            ? "text-sm font-semibold tabular-nums"
                            : "text-muted-foreground text-sm tabular-nums"
                        }
                      >
                        {balanceLabel(o.balance)}
                      </span>
                    </div>
                    <ChevronRight
                      aria-hidden
                      className="text-muted-foreground size-4 shrink-0"
                    />
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
