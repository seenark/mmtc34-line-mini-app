import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { Link } from "@tanstack/react-router";
import { useState } from "react";

import { errorText } from "@/api/me";
import { useConfirmSlip, useRejectSlip } from "@/api/payments";
import { type Order, PaymentPending } from "@/api/types.gen";
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
      <CardContent className="flex flex-col gap-3">
        {slip && <SlipImage name={slip.slip} />}
        <div className="flex flex-col gap-0.5">
          <Link
            to="/orders/$orderId"
            params={{ orderId: String(order.id) }}
            className="font-medium underline"
          >
            คำสั่งซื้อ #{order.id}
          </Link>
          <span>{order.campaignTitle}</span>
          <span className="text-muted-foreground text-sm">
            {order.memberNickname}
          </span>
        </div>
        <ul className="text-sm">
          {order.items.map((i) => (
            <li key={i.id}>
              {i.optionName}
              {i.size && ` ${i.size}`} × {i.quantity}
            </li>
          ))}
        </ul>
        <p className="font-semibold tabular-nums">ยอดรวม {baht(order.total)}</p>
        <div className="flex flex-col gap-1.5">
          <label htmlFor={`amount-${order.id}`} className="text-sm">
            ยอดที่ยืนยัน (บาท)
          </label>
          <Input
            id={`amount-${order.id}`}
            type="number"
            inputMode="numeric"
            min={1}
            step={1}
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
          />
          <Button
            size="lg"
            disabled={busy || !Number.isInteger(value) || value <= 0}
            onClick={() => confirm.mutate(value)}
          >
            ยืนยัน
          </Button>
        </div>
        <div className="flex flex-col gap-1.5">
          <label htmlFor={`reason-${order.id}`} className="text-sm">
            เหตุผลที่ปฏิเสธ
          </label>
          <Input
            id={`reason-${order.id}`}
            value={reason}
            onChange={(e) => setReason(e.target.value)}
          />
          <Button
            size="lg"
            variant="destructive"
            disabled={busy || reason.trim() === ""}
            onClick={() => reject.mutate(reason.trim())}
          >
            ปฏิเสธ
          </Button>
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
