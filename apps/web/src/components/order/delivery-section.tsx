import { Button } from "@mmtc34/ui/components/button";

import {
  type Delivery,
  DeliveryNone,
  DeliveryPickup,
  DeliveryShip,
  type Option,
  type OrderInput,
} from "@/api/types.gen";
import {
  AddressPicker,
  type ShipAddress,
} from "@/components/order/address-picker";

/** True when every Item is a donation (a no-price Option). Such an Order has no Delivery Method at all. */
export function isDonationOnly(
  lines: { optionId: number }[],
  options: Map<number, Option>,
) {
  return (
    lines.length > 0 &&
    lines.every((l) => options.get(l.optionId)?.price === null)
  );
}

/** `method` is the Member's choice between ship and pickup; it is ignored for a donation-only Order. */
export function deliveryInput(
  donationOnly: boolean,
  method: Delivery,
  address: ShipAddress,
): Pick<OrderInput, "delivery" | "shipName" | "shipAddress" | "shipPhone"> {
  const delivery = donationOnly ? DeliveryNone : method;
  const ship = delivery === DeliveryShip ? address.ship : null;
  return {
    delivery,
    shipName: ship?.name.trim() ?? "",
    shipAddress: ship?.address.trim() ?? "",
    shipPhone: ship?.phone.trim() ?? "",
  };
}

export function DeliverySection(props: {
  id: string;
  donationOnly: boolean;
  method: Delivery;
  onMethod: (method: Delivery) => void;
  address: ShipAddress;
  autofill?: boolean;
}) {
  if (props.donationOnly) return null;
  return (
    <section className="flex flex-col gap-3">
      <h2 className="text-base font-semibold">การรับสินค้า</h2>
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
            variant={props.method === value ? "default" : "outline"}
            aria-pressed={props.method === value}
            onClick={() => props.onMethod(value)}
          >
            {label}
          </Button>
        ))}
      </div>
      {props.method === DeliveryShip && (
        <AddressPicker
          id={props.id}
          address={props.address}
          autofill={props.autofill}
        />
      )}
    </section>
  );
}
