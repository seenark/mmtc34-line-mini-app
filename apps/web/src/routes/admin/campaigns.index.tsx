import { buttonVariants } from "@mmtc34/ui/components/button";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { adminCampaignsQuery } from "@/api/campaigns";
import { CampaignList } from "@/components/campaign/campaign-list";

export const Route = createFileRoute("/admin/campaigns/")({
  staticData: { nav: { slot: "admin", label: "อีเวนต์", order: 20 } },
  loader: ({ context }) =>
    context.queryClient.ensureQueryData(adminCampaignsQuery),
  component: AdminCampaigns,
});

function AdminCampaigns() {
  const { data: campaigns } = useSuspenseQuery(adminCampaignsQuery);
  return (
    <div className="flex flex-col gap-5">
      <div className="flex items-center justify-between gap-3">
        <h1 className="text-primary min-w-0 text-2xl font-extrabold">อีเวนต์</h1>
        <Link to="/admin/campaigns/new" className={buttonVariants()}>
          สร้างอีเวนต์
        </Link>
      </div>
      <CampaignList campaigns={campaigns} admin />
    </div>
  );
}
