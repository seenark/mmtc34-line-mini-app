import { createFileRoute, Link } from "@tanstack/react-router";

import { errorText } from "@/api/me";
import { useCreateUnlinkedMember } from "@/api/unlinked";
import { ProfileForm } from "@/components/profile-form";

export const Route = createFileRoute("/admin/members/new")({
  staticData: {
    nav: { slot: "admin", label: "เพิ่มสมาชิก (ไม่มี LINE)", order: 12 },
  },
  component: NewUnlinkedMember,
});

function NewUnlinkedMember() {
  const create = useCreateUnlinkedMember();
  return (
    <div className="flex flex-col gap-4">
      <Link
        to="/admin/members"
        search={{ status: "approved" }}
        className="text-sm underline"
      >
        ← สมาชิก
      </Link>
      <h1 className="text-2xl font-semibold">เพิ่มสมาชิก (ไม่มี LINE)</h1>
      <ProfileForm
        nicknameOnly
        submitLabel="เพิ่มสมาชิก"
        pending={create.isPending}
        error={create.error ? errorText(create.error) : undefined}
        onSubmit={(input) => create.mutate(input)}
      />
    </div>
  );
}
