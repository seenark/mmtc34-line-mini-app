import { Button, buttonVariants } from "@mmtc34/ui/components/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { errorText } from "@/api/me";
import { orderQuery, useCancelOrder } from "@/api/orders";
import { OrderAwaitingPayment, OrderCancelled } from "@/api/types.gen";
import { OrderHeader } from "@/components/order/order-header";
import { OrderSummary } from "@/components/order/order-summary";
import { PaymentHistory } from "@/components/order/payment-history";
import { PaymentInfo } from "@/components/order/payment-info";
import { SlipUpload } from "@/components/payment/slip-upload";

export const Route = createFileRoute("/_app/orders/$orderId/")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(orderQuery(params.orderId)),
  component: OrderPage,
});

function OrderPage() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const cancel = useCancelOrder(orderId);
  const { actions } = order;
  // actions merges the owner's and an Admin's rights; in member mode an Admin has only the owner's,
  // which is cancelling while awaiting payment.
  const canCancel = actions.cancel && order.status === OrderAwaitingPayment;
  return (
    <div className="flex flex-col gap-4">
      <OrderHeader order={order} />
      <OrderSummary order={order} />
      {order.status === OrderAwaitingPayment && <PaymentInfo order={order} />}
      <SlipUpload order={order} />
      <PaymentHistory order={order} />
      {actions.editDetails && (
        <Link
          to="/orders/$orderId/edit"
          params={{ orderId }}
          className={buttonVariants({ variant: "outline", size: "lg" })}
        >
          แก้ไขไซซ์ ชื่อ และที่อยู่
        </Link>
      )}
      {canCancel && (
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
      {!canCancel && order.status !== OrderCancelled && (
        <p role="status" className="bg-muted rounded-md p-3 text-sm">
          หากต้องการแก้จำนวนสินค้าหรือยกเลิกคำสั่งซื้อนี้ กรุณาติดต่อแอดมิน
        </p>
      )}
    </div>
  );
}
