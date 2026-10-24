import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { membersQuery } from "@/api/members";
import { sponsorsQuery } from "@/api/sponsors";
import { SponsorsEditor } from "@/components/sponsor/sponsor-editor";

export const Route = createFileRoute("/admin/campaigns/$campaignId/sponsors")({
  staticData: {
    nav: { slot: "adminCampaign", label: "สปอนเซอร์", order: 20 },
  },
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(sponsorsQuery(params.campaignId)),
      context.queryClient.ensureQueryData(membersQuery("approved")),
    ]),
  component: AdminSponsors,
});

function AdminSponsors() {
  const { campaignId } = Route.useParams();
  const { data: positions } = useSuspenseQuery(sponsorsQuery(campaignId));
  return <SponsorsEditor campaignId={campaignId} positions={positions} />;
}
