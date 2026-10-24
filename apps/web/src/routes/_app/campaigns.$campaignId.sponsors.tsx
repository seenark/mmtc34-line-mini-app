import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { sponsorsQuery } from "@/api/sponsors";
import { SponsorCard } from "@/components/sponsor/sponsor-card";

export const Route = createFileRoute("/_app/campaigns/$campaignId/sponsors")({
  staticData: {
    nav: { slot: "memberCampaign", label: "สปอนเซอร์", order: 20 },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(sponsorsQuery(params.campaignId)),
  component: Sponsors,
});

function Sponsors() {
  const { campaignId } = Route.useParams();
  const { data: positions } = useSuspenseQuery(sponsorsQuery(campaignId));
  if (positions.length === 0) {
    return <p className="text-muted-foreground text-sm">ยังไม่มีตำแหน่งสปอนเซอร์</p>;
  }
  return (
    <ul className="flex flex-col gap-3">
      {positions.map((p) => (
        <li key={p.id}>
          <SponsorCard position={p} />
        </li>
      ))}
    </ul>
  );
}
