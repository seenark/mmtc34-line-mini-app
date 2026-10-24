import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";

import { useRecordCash } from "@/api/ledger";
import { errorText } from "@/api/me";
import { orderQuery } from "@/api/orders";
import { baht, balanceLabel } from "@/components/order/order-status";

export const Route = createFileRoute("/admin/orders/$orderId/cash")({
  staticData: {
    nav: {
      slot: "orderAction",
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
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        save.mutate(value);
      }}
    >
      <Link
        to="/orders/$orderId"
        params={{ orderId }}
        className="text-sm underline"
      >
        ← คำสั่งซื้อ #{order.id}
      </Link>
      <h1 className="text-2xl font-semibold">บันทึกเงินสด #{order.id}</h1>
      <dl className="grid grid-cols-[1fr_auto] gap-x-2 gap-y-1 tabular-nums">
        <dt>ยอดรวม</dt>
        <dd className="text-right">{baht(order.total)}</dd>
        <dt>ชำระแล้ว</dt>
        <dd className="text-right">{baht(order.paid)}</dd>
        <dt className="font-semibold">คงเหลือ</dt>
        <dd className="text-right font-semibold">
          {balanceLabel(order.balance)}
        </dd>
      </dl>
      <div className="flex flex-col gap-1.5">
        <label htmlFor="amount" className="text-sm">
          จำนวนเงินสดที่รับ (บาท)
        </label>
        <Input
          id="amount"
          type="number"
          inputMode="numeric"
          min={1}
          step={1}
          value={amount}
          onChange={(e) => setAmount(e.target.value)}
        />
      </div>
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
        บันทึกเงินสด
      </Button>
    </form>
  );
}
