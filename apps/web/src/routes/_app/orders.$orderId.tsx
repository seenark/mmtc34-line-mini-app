import { Button } from "@mmtc34/ui/components/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

import { errorText } from "@/api/me";
import { orderQuery, useCancelOrder } from "@/api/orders";
import { OrderAwaitingPayment, OrderCancelled } from "@/api/types.gen";
import { OrderStatusBadge } from "@/components/order/order-status";
import { OrderSummary } from "@/components/order/order-summary";
import { PaymentHistory } from "@/components/order/payment-history";
import { PaymentInfo } from "@/components/order/payment-info";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/_app/orders/$orderId")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(orderQuery(params.orderId)),
  component: OrderPage,
});

function OrderPage() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const cancel = useCancelOrder(orderId);
  const { actions } = order;
  return (
    <div className="flex flex-col gap-4">
      <header className="flex flex-col gap-1">
        <div className="flex items-start justify-between gap-2">
          <h1 className="text-2xl font-semibold">คำสั่งซื้อ #{order.id}</h1>
          <OrderStatusBadge status={order.status} />
        </div>
        <Link
          to="/campaigns/$campaignId"
          params={{ campaignId: String(order.campaignId) }}
          className="text-muted-foreground text-sm underline"
        >
          {order.campaignTitle}
        </Link>
        <p className="text-muted-foreground text-sm">
          {order.memberNickname}
          {order.placedByAdmin && " · แอดมินสั่งแทน"}
        </p>
      </header>
      <OrderSummary order={order} />
      {order.status === OrderAwaitingPayment && (
        <PaymentInfo campaignId={order.campaignId} />
      )}
      <PaymentHistory order={order} />
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
      {!actions.cancel && !actions.edit && order.status !== OrderCancelled && (
        <p role="status" className="bg-muted p-3 text-sm">
          หากต้องการแก้ไขหรือยกเลิกคำสั่งซื้อนี้ กรุณาติดต่อแอดมิน
        </p>
      )}
      <NavSlot slot="orderAction" params={{ orderId }} actions={actions} />
      <Outlet />
    </div>
  );
}
