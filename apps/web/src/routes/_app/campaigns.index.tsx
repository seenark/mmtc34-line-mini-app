import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { campaignsQuery } from "@/api/campaigns";
import { CampaignList } from "@/components/campaign/campaign-list";

export const Route = createFileRoute("/_app/campaigns/")({
  staticData: { nav: { slot: "member", label: "แคมเปญ", order: 10 } },
  loader: ({ context }) => context.queryClient.ensureQueryData(campaignsQuery),
  component: Campaigns,
});

function Campaigns() {
  const { data: campaigns } = useSuspenseQuery(campaignsQuery);
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">แคมเปญ</h1>
      <CampaignList campaigns={campaigns} />
    </div>
  );
}
