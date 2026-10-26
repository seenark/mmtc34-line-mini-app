import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useId, useRef, useState } from "react";

import {
  type AddressChoice,
  type Delivery,
  DeliveryPickup,
  DeliveryShip,
  type ItemInput,
  type Option,
  type OrderInput,
} from "@/api/types.gen";
import { selectClass } from "@/components/campaign/field";
import { useShipAddress } from "@/components/order/address-picker";
import {
  DeliverySection,
  deliveryInput,
  isDonationOnly,
} from "@/components/order/delivery-section";
import { blankDetails, ItemFields } from "@/components/order/item-fields";
import type { ItemDetails } from "@/lib/cart";

type Line = ItemInput & { key: number };

/**
 * Admin editor for an Order's lines and delivery, used to place an Order on a Member's behalf and to
 * edit one. Lines keep their `id` so the server keeps the unit price captured at checkout.
 * `addresses` are the owner's saved shipping triples; the Order's own snapshot is offered as well.
 */
export function OrderForm(props: {
  options: Option[];
  initial: OrderInput;
  addresses: AddressChoice[];
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
  const [method, setMethod] = useState<Delivery>(
    props.initial.delivery === DeliveryPickup ? DeliveryPickup : DeliveryShip,
  );
  const address = useShipAddress(
    props.addresses,
    props.initial.delivery === DeliveryShip
      ? {
          name: props.initial.shipName,
          address: props.initial.shipAddress,
          phone: props.initial.shipPhone,
        }
      : undefined,
  );
  const [pick, setPick] = useState("");
  const byId = new Map(props.options.map((o) => [o.id, o]));
  const donationOnly = isDonationOnly(lines, byId);

  const update = (key: number, details: ItemDetails) =>
    setLines(lines.map((l) => (l.key === key ? { ...l, ...details } : l)));

  return (
    <form
      className="flex flex-col gap-6"
      onSubmit={(e) => {
        e.preventDefault();
        props.onSubmit({
          ...deliveryInput(donationOnly, method, address),
          items: lines.map(({ key: _, ...item }) => item),
        });
      }}
    >
      <section className="flex flex-col gap-3">
        <h2 className="text-base font-semibold">รายการสินค้า</h2>
        <ul className="flex flex-col gap-3">
          {lines.map((l) => {
            const o = byId.get(l.optionId);
            return (
              <li key={l.key}>
                <Card>
                  <CardContent className="flex flex-col gap-3">
                    <span className="text-sm font-semibold [overflow-wrap:anywhere]">
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
            className={`${selectClass} min-w-0 flex-1`}
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
      </section>

      <DeliverySection
        id={id}
        donationOnly={donationOnly}
        method={method}
        onMethod={setMethod}
        address={address}
      />

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
