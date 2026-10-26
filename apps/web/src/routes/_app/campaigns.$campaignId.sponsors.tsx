import { Empty, EmptyDescription } from "@mmtc34/ui/components/empty";
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
    return (
      <Empty className="flex-none rounded-xl border">
        <EmptyDescription>ยังไม่มีตำแหน่งสปอนเซอร์</EmptyDescription>
      </Empty>
    );
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
