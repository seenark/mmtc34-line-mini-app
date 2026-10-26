import { createFileRoute, Link } from "@tanstack/react-router";

import { useCreateCampaign } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { CampaignForm } from "@/components/campaign/campaign-form";

export const Route = createFileRoute("/admin/campaigns/new")({
  component: NewCampaign,
});

function NewCampaign() {
  const create = useCreateCampaign();
  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-1">
        <Link
          to="/admin/campaigns"
          className="text-secondary inline-flex min-h-9 w-fit items-center text-sm font-medium whitespace-nowrap underline-offset-4 transition-colors hover:underline active:opacity-70"
        >
          ← อีเวนต์
        </Link>
        <h1 className="text-primary text-2xl font-extrabold">สร้างอีเวนต์</h1>
      </div>
      <CampaignForm
        submitLabel="สร้างอีเวนต์"
        pending={create.isPending}
        error={create.error ? errorText(create.error) : undefined}
        onSubmit={(input) => create.mutate(input)}
      />
    </div>
  );
}
