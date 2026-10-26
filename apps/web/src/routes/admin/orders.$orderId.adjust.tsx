import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";

import { useAddAdjustment } from "@/api/ledger";
import { errorText } from "@/api/me";
import { orderQuery } from "@/api/orders";
import { Field } from "@/components/campaign/field";
import { baht, balanceLabel } from "@/components/order/order-status";

export const Route = createFileRoute("/admin/orders/$orderId/adjust")({
  staticData: {
    nav: {
      slot: "adminOrder",
      label: "ปรับยอด",
      order: 50,
      when: "addAdjust",
    },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(orderQuery(params.orderId)),
  component: AddAdjustment,
});

const directions = [
  {
    sign: 1,
    label: "เพิ่มยอด",
    effect: "คงเหลือเพิ่มขึ้น สมาชิกต้องจ่ายเพิ่มตามจำนวนนี้",
  },
  {
    sign: -1,
    label: "ลดยอดหรือคืนเงิน",
    effect: "คงเหลือลดลง ใช้กับส่วนลด หรือคืนเงินที่สมาชิกจ่ายเกิน",
  },
];

function AddAdjustment() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const save = useAddAdjustment(orderId);
  const [sign, setSign] = useState(1);
  const [amount, setAmount] = useState("");
  const [note, setNote] = useState("");
  const value = Number(amount);
  const effect = directions.find((d) => d.sign === sign)?.effect;

  return (
    <form
      className="flex flex-col gap-5"
      onSubmit={(e) => {
        e.preventDefault();
        save.mutate({ amount: sign * value, note: note.trim() });
      }}
    >
      <div className="flex flex-col gap-1">
        <Link
          to="/admin/orders/$orderId"
          params={{ orderId }}
          className="text-secondary inline-flex min-h-9 w-fit items-center text-sm font-medium whitespace-nowrap underline-offset-4 transition-colors hover:underline active:opacity-70"
        >
          ← คำสั่งซื้อ #{order.id}
        </Link>
        <h1 className="text-primary text-2xl font-extrabold">
          ปรับยอด #{order.id}
        </h1>
      </div>
      <Card size="sm">
        <CardContent>
          <dl className="grid grid-cols-[1fr_auto] gap-x-3 gap-y-1 text-sm tabular-nums">
            <dt className="text-muted-foreground">ยอดรวม</dt>
            <dd className="text-right">{baht(order.total)}</dd>
            <dt className="text-muted-foreground">ชำระแล้ว</dt>
            <dd className="text-right">{baht(order.paid)}</dd>
            <dt className="mt-1 border-t pt-2 font-semibold">คงเหลือ</dt>
            <dd className="mt-1 border-t pt-2 text-right font-semibold">
              {balanceLabel(order.balance)}
            </dd>
          </dl>
        </CardContent>
      </Card>
      <div className="flex flex-col gap-2">
        <div className="grid grid-cols-1 gap-2 min-[380px]:grid-cols-2">
          {directions.map((d) => (
            <Button
              key={d.sign}
              type="button"
              variant={d.sign === sign ? "default" : "outline"}
              aria-pressed={d.sign === sign}
              onClick={() => setSign(d.sign)}
            >
              {d.label}
            </Button>
          ))}
        </div>
        <p className="text-muted-foreground text-sm">{effect}</p>
      </div>
      <Field label="จำนวนเงิน (บาท)">
        <Input
          id="amount"
          type="number"
          inputMode="numeric"
          min={1}
          step={1}
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
      </Field>
      <Field label="เหตุผล">
        <Input
          id="note"
          value={note}
          onChange={(e) => setNote(e.target.value)}
        />
      </Field>
      {save.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(save.error)}
        </p>
      )}
      <Button
        type="submit"
        size="lg"
        disabled={
          save.isPending ||
          !Number.isInteger(value) ||
          value <= 0 ||
          note.trim() === ""
        }
      >
        {save.isPending ? "กำลังบันทึก…" : "บันทึกการปรับยอด"}
      </Button>
    </form>
  );
}
