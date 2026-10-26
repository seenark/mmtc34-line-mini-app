import { Button } from "@mmtc34/ui/components/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { useCancelOrderAsAdmin } from "@/api/admin-orders";
import { errorText } from "@/api/me";
import { memberQuery } from "@/api/members";
import { orderQuery } from "@/api/orders";
import { LinkOrderDialog } from "@/components/order/link-order-dialog";
import { OrderHeader } from "@/components/order/order-header";
import { OrderSummary } from "@/components/order/order-summary";
import { PaymentHistory } from "@/components/order/payment-history";
import { ReviewCard } from "@/components/payment/review-card";
import { SlipUpload } from "@/components/payment/slip-upload";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/admin/orders/$orderId/")({
  loader: async ({ context, params }) => {
    const order = await context.queryClient.ensureQueryData(
      orderQuery(params.orderId),
    );
    await context.queryClient.ensureQueryData(
      memberQuery(String(order.memberId)),
    );
  },
  component: AdminOrderPage,
});

function AdminOrderPage() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const { data: owner } = useSuspenseQuery(memberQuery(String(order.memberId)));
  const cancel = useCancelOrderAsAdmin(orderId);
  const { actions } = order;
  return (
    <div className="flex flex-col gap-4">
      <OrderHeader order={order} admin />
      <OrderSummary order={order} />
      <PaymentHistory order={order} />
      <SlipUpload order={order} />
      {actions.reviewSlip && <ReviewCard order={order} />}
      <NavSlot slot="adminOrder" params={{ orderId }} actions={actions} />
      {!owner.linked && !owner.deleted && (
        <LinkOrderDialog orderId={order.id} />
      )}
      {actions.cancel && (
        <Button
          variant="destructive"
          size="lg"
          disabled={cancel.isPending}
          onClick={() => {
            if (confirm("ยกเลิกคำสั่งซื้อนี้?")) cancel.mutate();
          }}
        >
          ยกเลิกคำสั่งซื้อ
        </Button>
      )}
      {cancel.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(cancel.error)}
        </p>
      )}
    </div>
  );
}
