import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";

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
  const navigate = useNavigate();
  return (
    <div className="flex flex-col gap-5">
      <div className="flex flex-col gap-1">
        <Link
          to="/admin/members"
          className="text-secondary inline-flex min-h-9 w-fit items-center text-sm font-medium whitespace-nowrap underline-offset-4 transition-colors hover:underline active:opacity-70"
        >
          ← สมาชิก
        </Link>
        <h1 className="text-primary text-2xl font-extrabold">
          เพิ่มสมาชิก (ไม่มี LINE)
        </h1>
      </div>
      <ProfileForm
        nicknameOnly
        submitLabel="เพิ่มสมาชิก"
        pending={create.isPending}
        error={create.error ? errorText(create.error) : undefined}
        onSubmit={(input) =>
          create.mutate(input, {
            onSuccess: (m) =>
              navigate({
                to: "/admin/members/$memberId",
                params: { memberId: String(m.id) },
              }),
          })
        }
      />
    </div>
  );
}
