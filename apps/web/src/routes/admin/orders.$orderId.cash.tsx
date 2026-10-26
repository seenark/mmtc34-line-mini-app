import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";

import { useRecordCash } from "@/api/ledger";
import { errorText } from "@/api/me";
import { orderQuery } from "@/api/orders";
import { Field } from "@/components/campaign/field";
import { baht, balanceLabel } from "@/components/order/order-status";

export const Route = createFileRoute("/admin/orders/$orderId/cash")({
  staticData: {
    nav: {
      slot: "adminOrder",
      label: "บันทึกเงินสด",
      order: 40,
      when: "recordCash",
    },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(orderQuery(params.orderId)),
  component: RecordCash,
});

function RecordCash() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const save = useRecordCash(orderId);
  const [amount, setAmount] = useState(
    String(order.balance > 0 ? order.balance : order.total),
  );
  const value = Number(amount);

  return (
    <form
      className="flex flex-col gap-5"
      onSubmit={(e) => {
        e.preventDefault();
        save.mutate(value);
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
          บันทึกเงินสด #{order.id}
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
      <Field label="จำนวนเงินสดที่รับ (บาท)">
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
      {save.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(save.error)}
        </p>
      )}
      <Button
        type="submit"
        size="lg"
        disabled={save.isPending || !Number.isInteger(value) || value <= 0}
      >
        {save.isPending ? "กำลังบันทึก…" : "บันทึกเงินสด"}
      </Button>
    </form>
  );
}
