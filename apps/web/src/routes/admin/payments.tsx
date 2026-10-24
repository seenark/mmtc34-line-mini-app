import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { pendingPaymentsQuery, pendingSlipCount } from "@/api/payments";
import { ReviewCard } from "@/components/payment/review-card";

export const Route = createFileRoute("/admin/payments")({
  staticData: {
    nav: {
      slot: "admin",
      label: "ตรวจสลิป",
      order: 15,
      badge: pendingSlipCount,
    },
  },
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(pendingPaymentsQuery),
  component: Payments,
});

function Payments() {
  const { data: orders } = useSuspenseQuery(pendingPaymentsQuery);
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">ตรวจสลิป</h1>
      {orders.length === 0 && (
        <p className="text-muted-foreground text-sm">ไม่มีสลิปรอตรวจ</p>
      )}
      <ul className="flex flex-col gap-3">
        {orders.map((o) => (
          <li key={o.id}>
            <ReviewCard order={o} />
          </li>
        ))}
      </ul>
    </div>
  );
}
