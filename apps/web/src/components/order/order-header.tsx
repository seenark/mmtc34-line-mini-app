import { Link } from "@tanstack/react-router";

import type { Order } from "@/api/types.gen";

import { OrderStatusBadge } from "./order-status";

const CAMPAIGN_LINK =
  "text-secondary focus-visible:ring-ring w-fit rounded-md text-sm underline underline-offset-4 outline-none focus-visible:ring-2";

/** `admin` points the Campaign link at the admin shell so an Admin never crosses into the member shell. */
export function OrderHeader({
  order,
  admin = false,
}: {
  order: Order;
  admin?: boolean;
}) {
  const params = { campaignId: String(order.campaignId) };
  return (
    <header className="flex flex-col gap-1">
      <div className="flex items-start justify-between gap-3">
        <h1 className="text-primary min-w-0 text-2xl font-extrabold break-words">
          คำสั่งซื้อ #{order.id}
        </h1>
        <OrderStatusBadge status={order.status} />
      </div>
      {admin ? (
        <Link
          to="/admin/campaigns/$campaignId"
          params={params}
          className={CAMPAIGN_LINK}
        >
          {order.campaignTitle}
        </Link>
      ) : (
        <Link
          to="/campaigns/$campaignId"
          params={params}
          className={CAMPAIGN_LINK}
        >
          {order.campaignTitle}
        </Link>
      )}
      <p className="text-muted-foreground text-sm">
        {order.memberNickname}
        {order.placedByAdmin && " · แอดมินสั่งแทน"}
      </p>
    </header>
  );
}
