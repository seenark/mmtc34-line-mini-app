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
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-1">
        <Link
          to="/admin/campaigns"
          className="text-secondary inline-flex min-h-9 w-fit items-center text-sm font-medium whitespace-nowrap underline-offset-4 transition-colors hover:underline active:opacity-70"
        >
          ← อีเวนต์
        </Link>
        <div className="flex items-start justify-between gap-3">
          <h1 className="text-primary min-w-0 text-2xl font-extrabold">
            {campaign.title}
          </h1>
          <StatusBadge status={campaign.status} />
        </div>
      </div>
      <NavSlot slot="adminCampaign" params={{ campaignId }} />
      <Outlet />
    </div>
  );
}
