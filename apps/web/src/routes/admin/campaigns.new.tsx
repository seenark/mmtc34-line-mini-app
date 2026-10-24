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
    <div className="flex flex-col gap-4">
      <Link to="/admin/campaigns" className="text-sm underline">
        ← แคมเปญ
      </Link>
      <h1 className="text-2xl font-semibold">สร้างแคมเปญ</h1>
      <CampaignForm
        submitLabel="สร้างแคมเปญ"
        pending={create.isPending}
        error={create.error ? errorText(create.error) : undefined}
        onSubmit={(input) => create.mutate(input)}
      />
    </div>
  );
}
