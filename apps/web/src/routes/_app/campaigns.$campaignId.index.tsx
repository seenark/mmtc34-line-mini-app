import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { campaignQuery } from "@/api/campaigns";
import { CampaignClosed } from "@/api/types.gen";
import { BankDetails } from "@/components/campaign/bank-details";
import { OptionCard } from "@/components/campaign/option-card";

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
  return (
    <div className="flex flex-col gap-4">
      {campaign.status === CampaignClosed && (
        <p role="status" className="bg-muted p-3 text-sm font-medium">
          ปิดรับแล้ว ติดต่อแอดมิน
        </p>
      )}
      {campaign.description && (
        <p className="text-sm whitespace-pre-wrap">{campaign.description}</p>
      )}
      <ul className="flex flex-col gap-3">
        {options.map((o) => (
          <li key={o.id}>
            <OptionCard option={o} />
          </li>
        ))}
      </ul>
      <BankDetails campaign={campaign} />
    </div>
  );
}
