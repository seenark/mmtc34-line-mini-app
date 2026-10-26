import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Lock } from "lucide-react";

import { campaignQuery } from "@/api/campaigns";
import { CampaignOpen } from "@/api/types.gen";
import { BankDetails } from "@/components/campaign/bank-details";
import { OptionCard } from "@/components/campaign/option-card";
import { AddToCart } from "@/components/order/add-to-cart";
import { CartBar } from "@/components/order/cart-bar";
import { useCart } from "@/lib/cart";

export const Route = createFileRoute("/_app/campaigns/$campaignId/")({
  staticData: {
    nav: { slot: "memberCampaign", label: "รายละเอียด", order: 10 },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(campaignQuery(params.campaignId)),
  component: CampaignDetail,
});

function CampaignDetail() {
  const { campaignId } = Route.useParams();
  const {
    data: { campaign, options },
  } = useSuspenseQuery(campaignQuery(campaignId));
  const cart = useCart(campaignId);
  const open = campaign.status === CampaignOpen;
  return (
    <div className="flex flex-col gap-4">
      {!open && (
        <p
          role="status"
          className="bg-muted flex items-center gap-2 rounded-md px-3 py-2.5 text-sm font-medium"
        >
          <Lock aria-hidden className="text-muted-foreground size-4 shrink-0" />
          ปิดรับแล้ว ติดต่อแอดมิน
        </p>
      )}
      {campaign.description && (
        <p className="text-sm [overflow-wrap:anywhere] whitespace-pre-wrap">
          {campaign.description}
        </p>
      )}
      <ul className="flex flex-col gap-3">
        {options.map((o) => (
          <li key={o.id}>
            <OptionCard option={o}>
              {open && <AddToCart option={o} onAdd={cart.add} />}
            </OptionCard>
          </li>
        ))}
      </ul>
      <BankDetails campaign={campaign} />
      {open && (
        <CartBar campaignId={campaignId} options={options} lines={cart.lines} />
      )}
    </div>
  );
}
