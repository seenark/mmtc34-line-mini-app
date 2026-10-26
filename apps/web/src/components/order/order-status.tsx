import {
  type Order,
  OrderAwaitingPayment,
  OrderCancelled,
  OrderPaid,
  OrderPendingReview,
} from "@/api/types.gen";
import { ToneBadge } from "@/components/campaign/status-badge";

const labels: Record<string, string> = {
  [OrderAwaitingPayment]: "รอชำระ",
  [OrderPendingReview]: "รอตรวจ",
  [OrderPaid]: "ชำระแล้ว",
  [OrderCancelled]: "ยกเลิก",
};

export function OrderStatusBadge({ status }: { status: Order["status"] }) {
  const label = labels[status] ?? status;
  if (status === OrderAwaitingPayment) {
    return (
      <ToneBadge
        tone="text-secondary ring-secondary ring-1 ring-inset"
        dot="bg-secondary"
      >
        {label}
      </ToneBadge>
    );
  }
  if (status === OrderPendingReview) {
    return (
      <ToneBadge tone="bg-accent text-accent-foreground" dot="bg-secondary">
        {label}
      </ToneBadge>
    );
  }
  if (status === OrderPaid) {
    return (
      <ToneBadge tone="bg-primary text-primary-foreground">{label}</ToneBadge>
    );
  }
  return <ToneBadge tone="bg-muted text-muted-foreground">{label}</ToneBadge>;
}

export const baht = (n: number) => `${n.toLocaleString("th-TH")} บาท`;

/** `balance` is the server's Order.Money() result; this only words it. */
export function balanceLabel(balance: number): string {
  if (balance > 0) return `ค้างจ่าย ${baht(balance)}`;
  if (balance < 0) return `จ่ายเกิน ${baht(-balance)}`;
  return "จ่ายครบแล้ว";
}
