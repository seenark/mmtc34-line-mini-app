import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { useId, useRef, useState } from "react";

import {
  DeliveryPickup,
  DeliveryShip,
  type ItemInput,
  type Option,
  type OrderInput,
} from "@/api/types.gen";
import { Field, selectClass } from "@/components/campaign/field";
import { blankDetails, ItemFields } from "@/components/order/item-fields";
import type { ItemDetails } from "@/lib/cart";

type Line = ItemInput & { key: number };

/**
 * Admin editor for an Order's lines and delivery, used to place an Order on a Member's behalf and to
 * edit one. Lines keep their `id` so the server keeps the unit price captured at checkout.
 */
export function OrderForm(props: {
  options: Option[];
  initial: OrderInput;
  submitLabel: string;
  pending: boolean;
  error?: string;
  onSubmit: (input: OrderInput) => void;
}) {
  const id = useId();
  const nextKey = useRef(props.initial.items.length);
  const [lines, setLines] = useState<Line[]>(() =>
    props.initial.items.map((item, key) => ({ ...item, key })),
  );
  const [delivery, setDelivery] = useState(props.initial.delivery);
  const [ship, setShip] = useState({
    name: props.initial.shipName,
    address: props.initial.shipAddress,
    phone: props.initial.shipPhone,
  });
  const [pick, setPick] = useState("");
  const byId = new Map(props.options.map((o) => [o.id, o]));
  const shipping = delivery === DeliveryShip;

  const update = (key: number, details: ItemDetails) =>
    setLines(lines.map((l) => (l.key === key ? { ...l, ...details } : l)));

  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        props.onSubmit({
          delivery,
          shipName: shipping ? ship.name.trim() : "",
          shipAddress: shipping ? ship.address.trim() : "",
          shipPhone: shipping ? ship.phone.trim() : "",
          items: lines.map(({ key: _, ...item }) => item),
        });
      }}
    >
      <ul className="flex flex-col gap-3">
        {lines.map((l) => {
          const o = byId.get(l.optionId);
          return (
            <li key={l.key}>
              <Card>
                <CardContent className="flex flex-col gap-3">
                  <span className="text-sm font-medium">
                    {o?.name ?? "ตัวเลือกนี้ไม่มีแล้ว"}
                  </span>
                  {o && (
                    <ItemFields
                      id={`${id}-${l.key}`}
                      option={o}
                      value={l}
                      onChange={(d) => update(l.key, d)}
                    />
                  )}
                  <Button
                    type="button"
                    variant="destructive"
                    onClick={() => setLines(lines.filter((x) => x !== l))}
                  >
                    ลบรายการ
                  </Button>
                </CardContent>
              </Card>
            </li>
          );
        })}
      </ul>

      <div className="flex gap-2">
        <select
          aria-label="เลือกสินค้า"
          className={selectClass}
          value={pick}
          onChange={(e) => setPick(e.target.value)}
        >
          <option value="">เลือกสินค้าที่จะเพิ่ม</option>
          {props.options.map((o) => (
            <option key={o.id} value={o.id}>
              {o.name}
            </option>
          ))}
        </select>
        <Button
          type="button"
          variant="outline"
          disabled={pick === ""}
          onClick={() => {
            setLines([
              ...lines,
              {
                ...blankDetails,
                key: nextKey.current++,
                id: null,
                optionId: Number(pick),
              },
            ]);
            setPick("");
          }}
        >
          เพิ่มรายการ
        </Button>
      </div>

      <fieldset className="flex flex-col gap-3">
        <legend className="mb-1 text-sm font-medium">การรับสินค้า</legend>
        <div className="grid grid-cols-2 gap-2">
          {(
            [
              [DeliveryShip, "จัดส่ง"],
              [DeliveryPickup, "รับเอง"],
            ] as const
          ).map(([value, label]) => (
            <Button
              key={value}
              type="button"
              variant={delivery === value ? "default" : "outline"}
              aria-pressed={delivery === value}
              onClick={() => setDelivery(value)}
            >
              {label}
            </Button>
          ))}
        </div>
        {shipping && (
          <>
            <Field label="ชื่อผู้รับ" required>
              <Input
                id={`${id}-shipName`}
                required
                autoComplete="off"
                value={ship.name}
                onChange={(e) => setShip({ ...ship, name: e.target.value })}
              />
            </Field>
            <Field label="ที่อยู่จัดส่ง" required>
              <Textarea
                id={`${id}-shipAddress`}
                required
                autoComplete="off"
                value={ship.address}
                onChange={(e) => setShip({ ...ship, address: e.target.value })}
              />
            </Field>
            <Field label="เบอร์โทร" required>
              <Input
                id={`${id}-shipPhone`}
                type="tel"
                required
                autoComplete="off"
                value={ship.phone}
                onChange={(e) => setShip({ ...ship, phone: e.target.value })}
              />
            </Field>
          </>
        )}
      </fieldset>

      {props.error && (
        <p role="alert" className="text-destructive text-sm">
          {props.error}
        </p>
      )}
      <Button
        type="submit"
        size="lg"
        disabled={props.pending || lines.length === 0}
      >
        {props.submitLabel}
      </Button>
    </form>
  );
}
