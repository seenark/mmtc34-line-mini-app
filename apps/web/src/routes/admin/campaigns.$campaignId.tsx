import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

import { adminCampaignQuery } from "@/api/campaigns";
import { StatusBadge } from "@/components/campaign/status-badge";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/admin/campaigns/$campaignId")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(adminCampaignQuery(params.campaignId)),
  component: AdminCampaignLayout,
});

function AdminCampaignLayout() {
  const { campaignId } = Route.useParams();
  const {
    data: { campaign },
  } = useSuspenseQuery(adminCampaignQuery(campaignId));
  return (
    <div className="flex flex-col gap-4">
      <Link to="/admin/campaigns" className="text-sm underline">
        ← แคมเปญ
      </Link>
      <div className="flex items-start justify-between gap-2">
        <h1 className="text-2xl font-semibold">{campaign.title}</h1>
        <StatusBadge status={campaign.status} />
      </div>
      <NavSlot slot="adminCampaign" params={{ campaignId }} />
      <Outlet />
    </div>
  );
}
