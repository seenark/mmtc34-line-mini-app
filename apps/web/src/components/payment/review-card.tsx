import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { errorText } from "@/api/me";
import { useConfirmSlip, useRejectSlip } from "@/api/payments";
import { type Order, PaymentPending } from "@/api/types.gen";
import { Field } from "@/components/campaign/field";
import { baht } from "@/components/order/order-status";

import { SlipImage } from "./slip-image";

export function ReviewCard({ order }: { order: Order }) {
  const slip = order.payments.findLast(
    (p) => p.status === PaymentPending && p.slip !== "",
  );
  const confirm = useConfirmSlip(String(order.id));
  const reject = useRejectSlip(String(order.id));
  const [amount, setAmount] = useState(String(order.total));
  const [reason, setReason] = useState("");
  const busy = confirm.isPending || reject.isPending;
  const error = confirm.error ?? reject.error;
  const value = Number(amount);

  return (
    <Card>
      <CardContent className="flex flex-col gap-4">
        {slip && <SlipImage name={slip.slip} />}
        <div className="flex min-w-0 flex-col gap-0.5">
          <Link
            to="/admin/orders/$orderId"
            params={{ orderId: String(order.id) }}
            className="text-secondary w-fit text-base font-semibold underline-offset-4 hover:underline active:opacity-70"
          >
            คำสั่งซื้อ #{order.id}
          </Link>
          <span className="text-sm">{order.campaignTitle}</span>
          <span className="text-muted-foreground text-sm">
            {order.memberNickname}
          </span>
        </div>
        <div className="flex flex-col gap-2 border-y py-3">
          <ul className="flex flex-col gap-0.5 text-sm">
            {order.items.map((i) => (
              <li key={i.id} className="min-w-0">
                {i.optionName}
                {i.size && ` ${i.size}`} × {i.quantity}
              </li>
            ))}
          </ul>
          <p className="text-primary font-semibold tabular-nums">
            ยอดรวม {baht(order.total)}
          </p>
        </div>
        <div className="grid gap-4 sm:grid-cols-2">
          <div className="flex min-w-0 flex-col gap-2">
            <Field label="ยอดที่ยืนยัน (บาท)">
              <Input
                id={`amount-${order.id}`}
                type="number"
                inputMode="numeric"
                min={1}
                step={1}
                value={amount}
                onChange={(e) => setAmount(e.target.value)}
              />
            </Field>
            <Button
              size="lg"
              disabled={busy || !Number.isInteger(value) || value <= 0}
              onClick={() => confirm.mutate(value)}
            >
              {confirm.isPending ? "กำลังยืนยัน…" : "ยืนยัน"}
            </Button>
          </div>
          <div className="flex min-w-0 flex-col gap-2">
            <Field label="เหตุผลที่ปฏิเสธ">
              <Input
                id={`reason-${order.id}`}
                value={reason}
                onChange={(e) => setReason(e.target.value)}
              />
            </Field>
            <Button
              size="lg"
              variant="destructive"
              disabled={busy || reason.trim() === ""}
              onClick={() => reject.mutate(reason.trim())}
            >
              {reject.isPending ? "กำลังปฏิเสธ…" : "ปฏิเสธ"}
            </Button>
          </div>
        </div>
        {error && (
          <p role="alert" className="text-destructive text-sm">
            {errorText(error)}
          </p>
        )}
      </CardContent>
    </Card>
  );
}
