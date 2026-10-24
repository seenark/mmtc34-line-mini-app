import { Card, CardContent } from "@mmtc34/ui/components/card";

import { slipUrl } from "@/api/payments";
import {
  type Order,
  PaymentConfirmed,
  PaymentPending,
  PaymentRejected,
} from "@/api/types.gen";

import { baht } from "./order-status";

const paymentLabels: Record<string, string> = {
  [PaymentPending]: "รอตรวจ",
  [PaymentConfirmed]: "ยืนยันแล้ว",
  [PaymentRejected]: "ไม่ผ่าน",
};

const when = (iso: string) =>
  new Date(iso).toLocaleString("th-TH", {
    dateStyle: "medium",
    timeStyle: "short",
  });

export function PaymentHistory({ order }: { order: Order }) {
  if (order.payments.length === 0 && order.adjustments.length === 0) {
    return null;
  }
  return (
    <Card>
      <CardContent className="flex flex-col gap-3">
        {order.payments.length > 0 && (
          <section className="flex flex-col gap-2">
            <h2 className="text-sm font-medium">การชำระเงิน</h2>
            <ul className="flex flex-col gap-2">
              {order.payments.map((p) => (
                <li key={p.id} className="flex flex-col border-t pt-2">
                  <div className="flex justify-between gap-2">
                    <span>
                      {paymentLabels[p.status] ?? p.status}
                      {p.slip === "" && " · เงินสด"}
                    </span>
                    {p.amount !== null && (
                      <span className="tabular-nums">{baht(p.amount)}</span>
                    )}
                  </div>
                  <span className="text-muted-foreground">
                    {when(p.createdAt)}
                  </span>
                  {p.slip !== "" && (
                    <a
                      href={slipUrl(p.slip)}
                      target="_blank"
                      rel="noreferrer"
                      className="mt-1 w-20"
                    >
                      <img
                        src={slipUrl(p.slip)}
                        alt="สลิป"
                        className="h-20 w-20 object-cover"
                      />
                    </a>
                  )}
                  {p.reviewedBy !== "" && p.reviewedAt && (
                    <span className="text-muted-foreground">
                      ตรวจโดย {p.reviewedBy} · {when(p.reviewedAt)}
                    </span>
                  )}
                  {p.status === PaymentRejected && p.rejectReason && (
                    <span className="text-destructive">
                      เหตุผล: {p.rejectReason}
                    </span>
                  )}
                </li>
              ))}
            </ul>
          </section>
        )}
        {order.adjustments.length > 0 && (
          <section className="flex flex-col gap-2">
            <h2 className="text-sm font-medium">ปรับยอด</h2>
            <ul className="flex flex-col gap-2">
              {order.adjustments.map((a) => (
                <li key={a.id} className="flex flex-col border-t pt-2">
                  <div className="flex justify-between gap-2">
                    <span>{a.note}</span>
                    <span className="tabular-nums">
                      {a.amount > 0 && "+"}
                      {baht(a.amount)}
                    </span>
                  </div>
                  <span className="text-muted-foreground">
                    {a.createdBy} · {when(a.createdAt)}
                  </span>
                </li>
              ))}
            </ul>
          </section>
        )}
      </CardContent>
    </Card>
  );
}
