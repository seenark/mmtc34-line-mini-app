import { Button, buttonVariants } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useId, useState } from "react";

import { campaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { useCheckout } from "@/api/orders";
import { myAddressesQuery } from "@/api/profile";
import { CampaignOpen, type Delivery, DeliveryShip } from "@/api/types.gen";
import { priceText } from "@/components/campaign/option-card";
import { useShipAddress } from "@/components/order/address-picker";
import {
  DeliverySection,
  deliveryInput,
  isDonationOnly,
} from "@/components/order/delivery-section";
import { ItemFields } from "@/components/order/item-fields";
import { estimateTotal, useCart } from "@/lib/cart";

export const Route = createFileRoute("/_app/campaigns/$campaignId/checkout")({
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(campaignQuery(params.campaignId)),
      context.queryClient.ensureQueryData(myAddressesQuery),
    ]),
  component: Checkout,
});

function Checkout() {
  const { campaignId } = Route.useParams();
  const {
    data: { campaign, options },
  } = useSuspenseQuery(campaignQuery(campaignId));
  const { data: addresses } = useSuspenseQuery(myAddressesQuery);
  const cart = useCart(campaignId);
  const place = useCheckout(campaignId, cart.clear);
  const id = useId();
  const [method, setMethod] = useState<Delivery>(DeliveryShip);
  const address = useShipAddress(addresses);
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
        <p
          role="status"
          className="bg-muted rounded-md p-3 text-sm font-medium"
        >
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
  const donationOnly = isDonationOnly(cart.lines, byId);

  return (
    <form
      className="flex flex-col gap-6"
      onSubmit={(e) => {
        e.preventDefault();
        place.mutate({
          ...deliveryInput(donationOnly, method, address),
          items: cart.lines.map((l) => ({ id: null, ...l })),
        });
      }}
    >
      <h1 className="text-primary text-2xl font-extrabold">ตะกร้าสินค้า</h1>
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
                    <span className="min-w-0 text-sm font-semibold [overflow-wrap:anywhere]">
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

      <DeliverySection
        id={id}
        donationOnly={donationOnly}
        method={method}
        onMethod={setMethod}
        address={address}
        autofill
      />

      <p className="border-t pt-4 text-sm">
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
