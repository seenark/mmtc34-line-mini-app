import { Empty, EmptyDescription } from "@mmtc34/ui/components/empty";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { rosterQuery } from "@/api/roster";
import type { RosterItem } from "@/api/types.gen";

export const Route = createFileRoute("/_app/campaigns/$campaignId/roster")({
  staticData: {
    nav: { slot: "memberCampaign", label: "รายชื่อ", order: 30 },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(rosterQuery(params.campaignId)),
  component: Roster,
});

const itemLabel = (i: RosterItem) =>
  `${[i.optionName, i.size].filter(Boolean).join(" ")} ×${i.quantity}`;

function Roster() {
  const { campaignId } = Route.useParams();
  const { data: roster } = useSuspenseQuery(rosterQuery(campaignId));
  if (roster.length === 0) {
    return (
      <Empty className="flex-none rounded-xl border">
        <EmptyDescription>ยังไม่มีคนที่ชำระเงิน</EmptyDescription>
        <EmptyDescription>รายชื่อจะแสดงคำสั่งซื้อหลังชำระเงินแล้ว</EmptyDescription>
      </Empty>
    );
  }
  const pieces = roster.reduce(
    (sum, e) => sum + e.items.reduce((s, i) => s + i.quantity, 0),
    0,
  );
  return (
    <div className="flex flex-col gap-3">
      <h2 className="text-primary text-base font-semibold">
        สั่งแล้ว {roster.length} คน · {pieces} ชิ้น
      </h2>
      <ol className="divide-border bg-card ring-foreground/10 flex flex-col divide-y overflow-hidden rounded-xl text-sm ring-1">
        {roster.map((e, idx) => (
          // biome-ignore lint/suspicious/noArrayIndexKey: nicknames repeat; list is never reordered
          <li key={idx} className="flex gap-3 px-4 py-3">
            <span className="text-muted-foreground w-6 shrink-0 text-right tabular-nums">
              {idx + 1}
            </span>
            <span className="flex min-w-0 flex-col [overflow-wrap:anywhere]">
              <span className="font-semibold">{e.nickname}</span>
              <span className="text-muted-foreground">
                {e.items.map(itemLabel).join(", ")}
              </span>
            </span>
          </li>
        ))}
      </ol>
    </div>
  );
}
