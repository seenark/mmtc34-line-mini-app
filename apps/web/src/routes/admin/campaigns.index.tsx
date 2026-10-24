import { buttonVariants } from "@mmtc34/ui/components/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { adminCampaignsQuery } from "@/api/campaigns";
import { CampaignList } from "@/components/campaign/campaign-list";

export const Route = createFileRoute("/admin/campaigns/")({
  staticData: { nav: { slot: "admin", label: "แคมเปญ", order: 20 } },
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(adminCampaignsQuery),
  component: AdminCampaigns,
});

function AdminCampaigns() {
  const { data: campaigns } = useSuspenseQuery(adminCampaignsQuery);
  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between gap-2">
        <h1 className="text-2xl font-semibold">แคมเปญ</h1>
        <Link to="/admin/campaigns/new" className={buttonVariants()}>
          สร้างแคมเปญ
        </Link>
      </div>
      <CampaignList campaigns={campaigns} admin />
    </div>
  );
}
