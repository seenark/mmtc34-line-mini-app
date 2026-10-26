import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { ImagePlus, Loader2, X } from "lucide-react";
import { useRef } from "react";

import { errorText } from "@/api/me";
import { slipUrl, useAddSlipImages, useRemoveSlipImage } from "@/api/payments";
import {
  type Order,
  type Payment,
  PaymentConfirmed,
  PaymentPending,
  PaymentRejected,
} from "@/api/types.gen";
import { ToneBadge } from "@/components/campaign/status-badge";

import { baht } from "./order-status";

const paymentLabels: Record<string, string> = {
  [PaymentPending]: "รอตรวจ",
  [PaymentConfirmed]: "ยืนยันแล้ว",
  [PaymentRejected]: "ไม่ผ่าน",
};

const paymentTones: Record<string, { tone: string; dot?: string }> = {
  [PaymentPending]: {
    tone: "bg-accent text-accent-foreground",
    dot: "bg-secondary",
  },
  [PaymentConfirmed]: { tone: "bg-primary text-primary-foreground" },
  [PaymentRejected]: { tone: "bg-destructive/10 text-destructive" },
};

const when = (iso: string) =>
  new Date(iso).toLocaleString("th-TH", {
    dateStyle: "medium",
    timeStyle: "short",
  });

function PaymentSlips({ order, p }: { order: Order; p: Payment }) {
  const orderId = String(order.id);
  const add = useAddSlipImages(orderId, p.id);
  const remove = useRemoveSlipImage(orderId, p.id);
  const input = useRef<HTMLInputElement>(null);
  const editable = order.actions.editSlips;
  const busy = add.isPending || remove.isPending;
  const error = add.error ?? remove.error;
  return (
    <div className="mt-1 flex flex-col gap-2">
      <ul className="flex flex-wrap gap-2">
        {p.slips.map((name) => (
          <li key={name} className="relative">
            <a
              href={slipUrl(name)}
              target="_blank"
              rel="noreferrer"
              className="ring-foreground/10 block w-20 overflow-hidden rounded-md ring-1 transition-opacity duration-(--dur-micro) ease-(--ease-out) hover:opacity-80 active:opacity-70"
            >
              <img
                src={slipUrl(name)}
                alt="สลิป"
                loading="lazy"
                className="bg-muted size-20 object-cover"
              />
            </a>
            {editable && (
              <Button
                type="button"
                variant="secondary"
                size="icon-sm"
                aria-label="ลบรูปสลิป"
                className="absolute top-0.5 right-0.5"
                disabled={busy || p.slips.length === 1}
                onClick={() => {
                  if (confirm("ลบรูปสลิปนี้?")) remove.mutate(name);
                }}
              >
                <X aria-hidden />
              </Button>
            )}
          </li>
        ))}
      </ul>
      {editable && (
        <>
          <input
            ref={input}
            type="file"
            accept="image/*"
            multiple
            className="sr-only"
            tabIndex={-1}
            aria-label="เพิ่มรูปสลิป"
            disabled={busy}
            onChange={(e) => {
              const picked = Array.from(e.target.files ?? []);
              e.target.value = "";
              if (picked.length > 0) add.mutate(picked);
            }}
          />
          <Button
            type="button"
            variant="outline"
            className="w-fit"
            disabled={busy}
            onClick={() => input.current?.click()}
          >
            {add.isPending ? (
              <Loader2 aria-hidden className="motion-safe:animate-spin" />
            ) : (
              <ImagePlus aria-hidden />
            )}
            เพิ่มรูปสลิป
          </Button>
        </>
      )}
      {error && (
        <p
          role="alert"
          className="bg-destructive/10 text-destructive rounded-md px-3 py-2 text-sm"
        >
          {errorText(error)}
        </p>
      )}
    </div>
  );
}

export function PaymentHistory({ order }: { order: Order }) {
  if (order.payments.length === 0 && order.adjustments.length === 0) {
    return null;
  }
  return (
    <Card>
      <CardContent className="flex flex-col gap-4">
        {order.payments.length > 0 && (
          <section className="flex flex-col gap-2">
            <h2 className="text-primary text-base font-semibold">การชำระเงิน</h2>
            <ul className="flex flex-col gap-3">
              {order.payments.map((p) => {
                const badge = paymentTones[p.status] ?? {
                  tone: "bg-muted text-muted-foreground",
                };
                return (
                  <li
                    key={p.id}
                    className="flex flex-col gap-1 border-t pt-3 text-sm"
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1">
                        <ToneBadge tone={badge.tone} dot={badge.dot}>
                          {paymentLabels[p.status] ?? p.status}
                        </ToneBadge>
                        {p.slips.length === 0 && (
                          <span className="text-muted-foreground">เงินสด</span>
                        )}
                      </span>
                      {p.amount !== null && (
                        <span className="font-semibold tabular-nums">
                          {baht(p.amount)}
                        </span>
                      )}
                    </div>
                    <span className="text-muted-foreground">
                      {when(p.createdAt)}
                    </span>
                    {p.slips.length > 0 && <PaymentSlips order={order} p={p} />}
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
                );
              })}
            </ul>
          </section>
        )}
        {order.adjustments.length > 0 && (
          <section className="flex flex-col gap-2">
            <h2 className="text-primary text-base font-semibold">ปรับยอด</h2>
            <ul className="flex flex-col gap-3">
              {order.adjustments.map((a) => (
                <li
                  key={a.id}
                  className="flex flex-col gap-1 border-t pt-3 text-sm"
                >
                  <div className="flex justify-between gap-2">
                    <span className="min-w-0 [overflow-wrap:anywhere]">
                      {a.note}
                    </span>
                    <span className="shrink-0 font-semibold tabular-nums">
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
