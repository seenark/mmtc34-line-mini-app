import {
  type Order,
  OrderAwaitingPayment,
  OrderCancelled,
  OrderPaid,
  OrderPendingReview,
} from "@/api/types.gen";

const labels: Record<string, string> = {
  [OrderAwaitingPayment]: "รอชำระ",
  [OrderPendingReview]: "รอตรวจ",
  [OrderPaid]: "ชำระแล้ว",
  [OrderCancelled]: "ยกเลิก",
};

export function OrderStatusBadge({ status }: { status: Order["status"] }) {
  const tone =
    status === OrderPaid
      ? "bg-primary text-primary-foreground"
      : status === OrderCancelled
        ? "bg-muted text-muted-foreground"
        : "border";
  return (
    <span className={`shrink-0 px-2 py-0.5 text-xs ${tone}`}>
      {labels[status] ?? status}
    </span>
  );
}

export const baht = (n: number) => `${n.toLocaleString("th-TH")} บาท`;

/** `balance` is the server's Order.Money() result; this only words it. */
export function balanceLabel(balance: number): string {
  if (balance > 0) return `ค้างจ่าย ${baht(balance)}`;
  if (balance < 0) return `จ่ายเกิน ${baht(-balance)}`;
  return "จ่ายครบแล้ว";
}
