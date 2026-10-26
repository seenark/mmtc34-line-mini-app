import { createFileRoute } from "@tanstack/react-router";

import { enforceStage, errorText, useOnboard } from "@/api/me";
import { AuthShell } from "@/components/auth-shell";
import { ProfileForm } from "@/components/profile-form";

export const Route = createFileRoute("/onboarding")({
  beforeLoad: ({ context }) => enforceStage(context.queryClient, "/onboarding"),
  component: Onboarding,
});

function Onboarding() {
  const onboard = useOnboard();
  return (
    <AuthShell heading="ลงทะเบียนสมาชิก" lede="กรอกข้อมูลของคุณเพื่อเริ่มสั่งซื้อได้ทันที">
      <ProfileForm
        submitLabel="เริ่มใช้งาน"
        pending={onboard.isPending}
        error={onboard.error ? errorText(onboard.error) : undefined}
        onSubmit={(input) => onboard.mutate(input)}
      />
    </AuthShell>
  );
}
