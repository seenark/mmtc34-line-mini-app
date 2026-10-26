import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { campaignsQuery } from "@/api/campaigns";
import { CampaignList } from "@/components/campaign/campaign-list";

export const Route = createFileRoute("/_app/")({
  loader: ({ context }) => context.queryClient.ensureQueryData(campaignsQuery),
  component: Home,
});

function Home() {
  const { member } = Route.useRouteContext();
  const { data: campaigns } = useSuspenseQuery(campaignsQuery);
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-primary text-2xl font-extrabold [overflow-wrap:anywhere]">
        สวัสดี {member.nickname}
      </h1>
      <CampaignList campaigns={campaigns} />
    </div>
  );
}
