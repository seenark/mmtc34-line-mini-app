import { Card, CardContent } from "@mmtc34/ui/components/card";

import { DeliveryShip, type Order } from "@/api/types.gen";

import { baht, balanceLabel } from "./order-status";

const itemDetails = (i: Order["items"][number]) =>
  [i.size, i.printName && `ชื่อ ${i.printName}`, i.number && `เบอร์ ${i.number}`]
    .filter(Boolean)
    .join(" · ");

export function OrderSummary({ order }: { order: Order }) {
  return (
    <>
      <Card>
        <CardContent className="flex flex-col gap-3">
          <h2 className="text-sm font-medium">รายการ</h2>
          <table className="w-full text-left">
            <thead className="text-muted-foreground">
              <tr>
                <th className="py-1 font-normal">สินค้า</th>
                <th className="py-1 text-right font-normal">จำนวน</th>
                <th className="py-1 text-right font-normal">รวม</th>
              </tr>
            </thead>
            <tbody>
              {order.items.map((i) => (
                <tr key={i.id} className="border-t align-top">
                  <td className="py-1.5 pr-2">
                    <div>{i.optionName}</div>
                    <div className="text-muted-foreground">
                      {itemDetails(i)}
                    </div>
                  </td>
                  <td className="py-1.5 text-right tabular-nums">
                    {i.quantity} × {i.unitPrice.toLocaleString("th-TH")}
                  </td>
                  <td className="py-1.5 text-right tabular-nums">
                    {(i.quantity * i.unitPrice).toLocaleString("th-TH")}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          <dl className="grid grid-cols-[1fr_auto] gap-x-2 gap-y-1 border-t pt-2 tabular-nums">
            <dt>ยอดรวม</dt>
            <dd className="text-right">{baht(order.total)}</dd>
            <dt>ชำระแล้ว</dt>
            <dd className="text-right">{baht(order.paid)}</dd>
            <dt className="font-semibold">คงเหลือ</dt>
            <dd className="text-right font-semibold">
              {balanceLabel(order.balance)}
            </dd>
          </dl>
        </CardContent>
      </Card>
      <Card>
        <CardContent className="flex flex-col gap-1">
          <h2 className="text-sm font-medium">การรับสินค้า</h2>
          {order.delivery === DeliveryShip ? (
            <>
              <p>จัดส่ง</p>
              <p>
                {order.shipName}
                {order.shipPhone && ` · ${order.shipPhone}`}
              </p>
              <p className="whitespace-pre-wrap">{order.shipAddress}</p>
            </>
          ) : (
            <p>รับเอง</p>
          )}
        </CardContent>
      </Card>
    </>
  );
}
