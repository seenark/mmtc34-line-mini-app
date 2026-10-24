import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";

import { campaignQuery } from "@/api/campaigns";
import { closeAtText, StatusBadge } from "@/components/campaign/status-badge";
import { imageUrl } from "@/lib/image";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/_app/campaigns/$campaignId")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(campaignQuery(params.campaignId)),
  component: CampaignLayout,
});

function CampaignLayout() {
  const { campaignId } = Route.useParams();
  const {
    data: { campaign },
  } = useSuspenseQuery(campaignQuery(campaignId));
  return (
    <div className="flex flex-col gap-4">
      {campaign.coverImage && (
        <img
          src={imageUrl(campaign.coverImage)}
          alt=""
          className="aspect-video w-full object-cover"
        />
      )}
      <header className="flex flex-col gap-1">
        <div className="flex items-start justify-between gap-2">
          <h1 className="text-2xl font-semibold">{campaign.title}</h1>
          <StatusBadge status={campaign.status} />
        </div>
        <p className="text-muted-foreground text-sm">
          {closeAtText(campaign.closeAt)}
        </p>
      </header>
      <NavSlot slot="memberCampaign" params={{ campaignId }} />
      <Outlet />
    </div>
  );
}
