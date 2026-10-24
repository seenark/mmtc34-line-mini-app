import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { adminCampaignQuery, useUpdateCampaign } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { CampaignForm } from "@/components/campaign/campaign-form";
import { OptionsEditor } from "@/components/campaign/option-editor";

export const Route = createFileRoute("/admin/campaigns/$campaignId/")({
  staticData: {
    nav: { slot: "adminCampaign", label: "ตั้งค่า", order: 10 },
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(adminCampaignQuery(params.campaignId)),
  component: CampaignSettings,
});

function CampaignSettings() {
  const { campaignId } = Route.useParams();
  const {
    data: { campaign, options },
  } = useSuspenseQuery(adminCampaignQuery(campaignId));
  const update = useUpdateCampaign(campaignId);
  return (
    <div className="flex flex-col gap-6">
      <CampaignForm
        initial={campaign}
        submitLabel="บันทึก"
        pending={update.isPending}
        error={update.error ? errorText(update.error) : undefined}
        onSubmit={(input) => update.mutate(input)}
      />
      {update.isSuccess && (
        <p className="text-muted-foreground text-sm">บันทึกแล้ว</p>
      )}
      <OptionsEditor campaignId={campaignId} options={options} />
    </div>
  );
}
