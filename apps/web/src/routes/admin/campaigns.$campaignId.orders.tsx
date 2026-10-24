import { buttonVariants } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { campaignOrdersQuery } from "@/api/ledger";
import {
  OrderAwaitingPayment,
  OrderCancelled,
  OrderPaid,
  OrderPendingReview,
} from "@/api/types.gen";
import {
  baht,
  balanceLabel,
  OrderStatusBadge,
} from "@/components/order/order-status";

const filters: { status: string | undefined; label: string }[] = [
  { status: undefined, label: "ทั้งหมด" },
  { status: OrderAwaitingPayment, label: "รอชำระ" },
  { status: OrderPendingReview, label: "รอตรวจ" },
  { status: OrderPaid, label: "ชำระแล้ว" },
  { status: OrderCancelled, label: "ยกเลิก" },
];

export const Route = createFileRoute("/admin/campaigns/$campaignId/orders")({
  staticData: {
    nav: { slot: "adminCampaign", label: "คำสั่งซื้อ", order: 15 },
  },
  validateSearch: (search: { status?: string }) => ({
    status: filters.find((f) => f.status === search.status)?.status,
  }),
  loaderDeps: ({ search }) => ({ status: search.status }),
  loader: ({ context, params, deps }) =>
    context.queryClient.ensureQueryData(
      campaignOrdersQuery(params.campaignId, deps.status),
    ),
  component: CampaignOrders,
});

function CampaignOrders() {
  const { campaignId } = Route.useParams();
  const { status } = Route.useSearch();
  const { data: orders } = useSuspenseQuery(
    campaignOrdersQuery(campaignId, status),
  );
  // balances are the server's Order.Money(); only the display total is summed here
  const owed = orders.reduce((sum, o) => sum + Math.max(o.balance, 0), 0);

  return (
    <div className="flex flex-col gap-4">
      <div className="flex gap-2 overflow-x-auto">
        {filters.map((f) => (
          <Link
            key={f.label}
            from={Route.fullPath}
            search={{ status: f.status }}
            className={`${buttonVariants({ variant: f.status === status ? "default" : "outline" })} shrink-0 whitespace-nowrap`}
          >
            {f.label}
          </Link>
        ))}
      </div>
      <p className="text-sm tabular-nums">
        {orders.length} คำสั่งซื้อ · ยอดค้างรวม {baht(owed)}
      </p>
      {orders.length === 0 && (
        <p className="text-muted-foreground text-sm">ไม่มีคำสั่งซื้อ</p>
      )}
      <ul className="flex flex-col gap-2">
        {orders.map((o) => (
          <li key={o.id}>
            <Link to="/orders/$orderId" params={{ orderId: String(o.id) }}>
              <Card>
                <CardContent className="flex flex-col gap-1.5">
                  <div className="flex items-center justify-between gap-2">
                    <span className="font-medium">
                      #{o.id} {o.memberNickname}
                    </span>
                    <span className="flex shrink-0 items-center gap-1">
                      {o.placedByAdmin && (
                        <span className="border px-2 py-0.5 text-xs">
                          แอดมินสั่งแทน
                        </span>
                      )}
                      <OrderStatusBadge status={o.status} />
                    </span>
                  </div>
                  <p className="text-muted-foreground text-sm">
                    {o.items
                      .map(
                        (i) =>
                          `${i.optionName}${i.size && ` ${i.size}`} × ${i.quantity}`,
                      )
                      .join(", ")}
                  </p>
                  <dl className="grid grid-cols-[1fr_auto] gap-x-2 text-sm tabular-nums">
                    <dt>ยอดรวม</dt>
                    <dd className="text-right">{baht(o.total)}</dd>
                    <dt>ชำระแล้ว</dt>
                    <dd className="text-right">{baht(o.paid)}</dd>
                    <dt className="font-semibold">คงเหลือ</dt>
                    <dd className="text-right font-semibold">
                      {balanceLabel(o.balance)}
                    </dd>
                  </dl>
                </CardContent>
              </Card>
            </Link>
          </li>
        ))}
      </ul>
    </div>
  );
}
