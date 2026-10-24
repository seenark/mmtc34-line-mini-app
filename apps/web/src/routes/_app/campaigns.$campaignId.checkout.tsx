import { Button, buttonVariants } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useId, useState } from "react";

import { campaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { useCheckout } from "@/api/orders";
import {
  CampaignOpen,
  type Delivery,
  DeliveryPickup,
  DeliveryShip,
} from "@/api/types.gen";
import { Field } from "@/components/campaign/field";
import { priceText } from "@/components/campaign/option-card";
import { ItemFields } from "@/components/order/item-fields";
import { estimateTotal, useCart } from "@/lib/cart";

export const Route = createFileRoute("/_app/campaigns/$campaignId/checkout")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(campaignQuery(params.campaignId)),
  component: Checkout,
});

function Checkout() {
  const { campaignId } = Route.useParams();
  const { member } = Route.useRouteContext();
  const {
    data: { campaign, options },
  } = useSuspenseQuery(campaignQuery(campaignId));
  const cart = useCart(campaignId);
  const place = useCheckout(campaignId, cart.clear);
  const id = useId();
  const [delivery, setDelivery] = useState<Delivery>(DeliveryShip);
  const [ship, setShip] = useState({
    name: member.realName,
    address: member.address,
    phone: member.phone,
  });
  const back = (
    <Link
      to="/campaigns/$campaignId"
      params={{ campaignId }}
      className={buttonVariants({ variant: "outline" })}
    >
      กลับไปเลือกสินค้า
    </Link>
  );

  if (campaign.status !== CampaignOpen) {
    return (
      <div className="flex flex-col gap-4">
        <p role="status" className="bg-muted p-3 text-sm font-medium">
          ปิดรับแล้ว ติดต่อแอดมิน
        </p>
        {back}
      </div>
    );
  }
  if (cart.lines.length === 0) {
    return (
      <div className="flex flex-col gap-4">
        <p className="text-muted-foreground text-sm">ตะกร้าว่างเปล่า</p>
        {back}
      </div>
    );
  }

  const byId = new Map(options.map((o) => [o.id, o]));
  const missing = cart.lines.some((l) => !byId.has(l.optionId));
  const shipping = delivery === DeliveryShip;

  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        place.mutate({
          delivery,
          shipName: shipping ? ship.name.trim() : "",
          shipAddress: shipping ? ship.address.trim() : "",
          shipPhone: shipping ? ship.phone.trim() : "",
          items: cart.lines.map((l) => ({ id: null, ...l })),
        });
      }}
    >
      <h2 className="text-xl font-semibold">ตะกร้าสินค้า</h2>
      <ul className="flex flex-col gap-3">
        {cart.lines.map((l, i) => {
          const o = byId.get(l.optionId);
          return (
            // lines have no id of their own; the cart is edited by position
            // biome-ignore lint/suspicious/noArrayIndexKey: see above
            <li key={i}>
              <Card>
                <CardContent className="flex flex-col gap-3">
                  <div className="flex items-start justify-between gap-2">
                    <span className="text-sm font-medium">
                      {o?.name ?? "ตัวเลือกนี้ไม่มีแล้ว"}
                    </span>
                    {o && (
                      <span className="shrink-0 text-sm">
                        {priceText(l.amount ?? o.price)}
                      </span>
                    )}
                  </div>
                  {o && (
                    <ItemFields
                      id={`${id}-${i}`}
                      option={o}
                      value={l}
                      onChange={(d) => cart.update(i, d)}
                    />
                  )}
                  <Button
                    type="button"
                    variant="destructive"
                    onClick={() => cart.remove(i)}
                  >
                    ลบออกจากตะกร้า
                  </Button>
                </CardContent>
              </Card>
            </li>
          );
        })}
      </ul>

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
                autoComplete="name"
                value={ship.name}
                onChange={(e) => setShip({ ...ship, name: e.target.value })}
              />
            </Field>
            <Field label="ที่อยู่จัดส่ง" required>
              <Textarea
                id={`${id}-shipAddress`}
                required
                autoComplete="street-address"
                value={ship.address}
                onChange={(e) => setShip({ ...ship, address: e.target.value })}
              />
            </Field>
            <Field label="เบอร์โทร" required>
              <Input
                id={`${id}-shipPhone`}
                type="tel"
                required
                autoComplete="tel"
                value={ship.phone}
                onChange={(e) => setShip({ ...ship, phone: e.target.value })}
              />
            </Field>
          </>
        )}
      </fieldset>

      <p className="text-sm">
        ยอดประมาณ{" "}
        <span className="font-semibold">
          {priceText(estimateTotal(options, cart.lines))}
        </span>
        <span className="text-muted-foreground">
          {" "}
          (ยอดจริงคำนวณโดยระบบเมื่อสั่งซื้อ)
        </span>
      </p>
      {missing && (
        <p role="alert" className="text-destructive text-sm">
          มีสินค้าที่ไม่มีให้สั่งแล้ว กรุณาลบออกจากตะกร้า
        </p>
      )}
      {place.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(place.error)}
        </p>
      )}
      <Button type="submit" size="lg" disabled={place.isPending || missing}>
        สั่งซื้อ
      </Button>
    </form>
  );
}
