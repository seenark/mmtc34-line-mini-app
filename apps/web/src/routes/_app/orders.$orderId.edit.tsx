import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useId, useState } from "react";

import { campaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { orderQuery, useEditOrderDetails } from "@/api/orders";
import { myAddressesQuery } from "@/api/profile";
import { DeliveryShip } from "@/api/types.gen";
import {
  AddressPicker,
  useShipAddress,
} from "@/components/order/address-picker";
import { ItemFields } from "@/components/order/item-fields";
import type { ItemDetails } from "@/lib/cart";

export const Route = createFileRoute("/_app/orders/$orderId/edit")({
  loader: async ({ context, params }) => {
    const order = await context.queryClient.ensureQueryData(
      orderQuery(params.orderId),
    );
    await Promise.all([
      context.queryClient.ensureQueryData(
        campaignQuery(String(order.campaignId)),
      ),
      context.queryClient.ensureQueryData(myAddressesQuery),
    ]);
  },
  component: EditOrderDetails,
});

function EditOrderDetails() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const {
    data: { options },
  } = useSuspenseQuery(campaignQuery(String(order.campaignId)));
  const { data: addresses } = useSuspenseQuery(myAddressesQuery);
  const save = useEditOrderDetails(orderId);
  const id = useId();
  const [items, setItems] = useState<Record<number, ItemDetails>>(() =>
    Object.fromEntries(
      order.items.map((i) => [
        i.id,
        {
          size: i.size,
          printName: i.printName,
          number: i.number,
          quantity: i.quantity,
          amount: null,
        },
      ]),
    ),
  );
  const ships = order.delivery === DeliveryShip;
  const address = useShipAddress(
    addresses,
    ships
      ? {
          name: order.shipName,
          address: order.shipAddress,
          phone: order.shipPhone,
        }
      : undefined,
  );
  const byId = new Map(options.map((o) => [o.id, o]));
  // donation Items have no details to edit
  const editable = order.items.filter((i) => {
    const o = byId.get(i.optionId);
    return o && (o.sizes.length > 0 || o.needsPrintName || o.needsNumber);
  });

  return (
    <div className="flex flex-col gap-4">
      <Link
        to="/orders/$orderId"
        params={{ orderId }}
        className="text-secondary inline-flex min-h-9 w-fit items-center text-sm font-medium whitespace-nowrap underline-offset-4 transition-colors hover:underline active:opacity-70"
      >
        ← คำสั่งซื้อ #{order.id}
      </Link>
      <h1 className="text-primary text-2xl font-extrabold">แก้ไขรายละเอียด</h1>
      <p role="note" className="bg-muted rounded-md p-3 text-sm">
        แก้จำนวนสินค้าไม่ได้ หากต้องการแก้จำนวนกรุณาติดต่อแอดมิน
      </p>
      <form
        className="flex flex-col gap-6"
        onSubmit={(e) => {
          e.preventDefault();
          save.mutate({
            shipName: ships ? address.ship.name.trim() : "",
            shipAddress: ships ? address.ship.address.trim() : "",
            shipPhone: ships ? address.ship.phone.trim() : "",
            items: order.items.map((i) => ({
              id: i.id,
              size: items[i.id].size,
              printName: items[i.id].printName,
              number: items[i.id].number,
            })),
          });
        }}
      >
        <ul className="flex flex-col gap-3">
          {editable.map((i) => {
            const o = byId.get(i.optionId);
            if (!o) return null;
            return (
              <li key={i.id}>
                <Card>
                  <CardContent className="flex flex-col gap-3">
                    <span className="text-sm font-semibold [overflow-wrap:anywhere]">
                      {o.name} × {i.quantity}
                    </span>
                    <ItemFields
                      id={`${id}-${i.id}`}
                      option={o}
                      value={items[i.id]}
                      detailsOnly
                      onChange={(d) => setItems({ ...items, [i.id]: d })}
                    />
                  </CardContent>
                </Card>
              </li>
            );
          })}
        </ul>
        {ships && (
          <section className="flex flex-col gap-3">
            <h2 className="text-base font-semibold">ที่อยู่จัดส่ง</h2>
            <AddressPicker id={id} address={address} autofill />
          </section>
        )}
        {save.error && (
          <p role="alert" className="text-destructive text-sm">
            {errorText(save.error)}
          </p>
        )}
        <Button type="submit" size="lg" disabled={save.isPending}>
          บันทึก
        </Button>
      </form>
    </div>
  );
}
