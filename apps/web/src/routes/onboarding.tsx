import { createFileRoute } from "@tanstack/react-router";

import { enforceStage, errorText, useOnboard } from "@/api/me";
import { ProfileForm } from "@/components/profile-form";

export const Route = createFileRoute("/onboarding")({
  beforeLoad: ({ context }) => enforceStage(context.queryClient, "/onboarding"),
  component: Onboarding,
});

function Onboarding() {
  const onboard = useOnboard();
  return (
    <main className="mx-auto flex min-h-svh max-w-md flex-col gap-6 p-6">
      <header>
        <h1 className="text-2xl font-semibold">ลงทะเบียนสมาชิก</h1>
        <p className="text-muted-foreground mt-1 text-sm">
          กรอกข้อมูลของคุณ แอดมินจะตรวจสอบและอนุมัติ
        </p>
      </header>
      <ProfileForm
        submitLabel="ส่งใบสมัคร"
        pending={onboard.isPending}
        error={onboard.error ? errorText(onboard.error) : undefined}
        onSubmit={(input) => onboard.mutate(input)}
      />
    </main>
  );
}
