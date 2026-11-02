import { useSuspenseQuery } from "@tanstack/react-query";
import {
  createFileRoute,
  Outlet,
  useChildMatches,
} from "@tanstack/react-router";
import { CalendarClock } from "lucide-react";

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
  const onDetail = useChildMatches({
    select: (matches) =>
      matches.at(-1)?.routeId === "/_app/campaigns/$campaignId/",
  });
  const closeAt = closeAtText(campaign.closeAt);
  const cover = campaign.photos[0];
  return (
    <div className="flex flex-col gap-4">
      {onDetail ? (
        <>
          {cover && (
            // The jersey cover is a square product photo; other shapes letterbox inside the same box.
            <img
              src={imageUrl(cover)}
              alt=""
              className="bg-muted aspect-square max-h-80 w-full rounded-xl object-contain"
            />
          )}
          <header className="flex flex-col gap-2">
            <div className="flex items-start justify-between gap-3">
              <h1 className="text-primary min-w-0 text-2xl font-extrabold [overflow-wrap:anywhere]">
                {campaign.title}
              </h1>
              <span className="mt-2 shrink-0">
                <StatusBadge status={campaign.status} />
              </span>
            </div>
            {closeAt && (
              <p className="text-muted-foreground flex items-center gap-1.5 text-sm">
                <CalendarClock aria-hidden className="size-4 shrink-0" />
                {closeAt}
              </p>
            )}
          </header>
        </>
      ) : (
        <header className="flex items-center gap-3">
          {cover && (
            <img
              src={imageUrl(cover)}
              alt=""
              className="bg-muted size-12 shrink-0 rounded-md object-contain"
            />
          )}
          <h1 className="text-primary min-w-0 flex-1 text-base font-semibold [overflow-wrap:anywhere]">
            {campaign.title}
          </h1>
          <StatusBadge status={campaign.status} />
        </header>
      )}
      <NavSlot slot="memberCampaign" params={{ campaignId }} />
      <Outlet />
    </div>
  );
}
