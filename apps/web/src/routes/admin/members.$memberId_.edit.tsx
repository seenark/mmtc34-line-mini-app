import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { useState } from "react";

import { errorText } from "@/api/me";
import { memberQuery } from "@/api/members";
import { useAdminEditMember, useSetRole } from "@/api/profile";
import { RoleAdmin, RoleMember } from "@/api/types.gen";
import { ProfileForm } from "@/components/profile-form";

export const Route = createFileRoute("/admin/members/$memberId_/edit")({
  staticData: { nav: { slot: "adminMember", label: "แก้ไข", order: 10 } },
  component: EditMember,
});

function EditMember() {
  const { memberId } = Route.useParams();
  const { data: m } = useSuspenseQuery(memberQuery(memberId));
  const edit = useAdminEditMember(memberId);
  return (
    <div className="flex flex-col gap-8">
      <div className="flex flex-col gap-1">
        <Link
          to="/admin/members/$memberId"
          params={{ memberId }}
          className="text-secondary inline-flex min-h-9 w-fit max-w-full items-center text-sm font-medium underline-offset-4 transition-colors hover:underline active:opacity-70"
        >
          <span className="truncate">← {m.nickname}</span>
        </Link>
        <h1 className="text-primary text-2xl font-extrabold">แก้ไขสมาชิก</h1>
      </div>
      <section className="flex flex-col gap-3">
        <h2 className="text-base font-semibold">ข้อมูลสมาชิก</h2>
        {/* keyed by the saved record so the fields reset to what the server stored */}
        <ProfileForm
          key={`${m.id}:${m.nickname}:${m.cohort}:${m.realName}:${m.address}:${m.phone}`}
          submitLabel="บันทึก"
          pending={edit.isPending}
          error={edit.error ? errorText(edit.error) : undefined}
          initial={m}
          nicknameOnly={!m.linked}
          onSubmit={(input) => edit.mutate(input)}
        />
        {edit.isSuccess && (
          <p role="status" className="text-secondary text-sm font-medium">
            บันทึกข้อมูลแล้ว
          </p>
        )}
      </section>
      {!m.superAdmin && (
        <RoleToggle memberId={memberId} role={m.role} nickname={m.nickname} />
      )}
    </div>
  );
}

function RoleToggle(props: {
  memberId: string;
  role: string;
  nickname: string;
}) {
  const setRole = useSetRole(props.memberId);
  const [confirming, setConfirming] = useState(false);
  const next = props.role === RoleAdmin ? RoleMember : RoleAdmin;
  const nextLabel = next === RoleAdmin ? "แอดมิน" : "สมาชิก";
  return (
    <Card>
      <CardContent className="flex flex-col gap-3">
        <h2 className="text-base font-semibold">บทบาท</h2>
        <div className="flex gap-2">
          {[RoleMember, RoleAdmin].map((r) => (
            <Button
              key={r}
              variant={props.role === r ? "default" : "outline"}
              size="lg"
              className="flex-1"
              aria-pressed={props.role === r}
              disabled={setRole.isPending}
              onClick={() => setConfirming(r !== props.role)}
            >
              {r === RoleAdmin ? "แอดมิน" : "สมาชิก"}
            </Button>
          ))}
        </div>
        {setRole.error && (
          <p role="alert" className="text-destructive text-sm">
            {errorText(setRole.error)}
          </p>
        )}
        {confirming && (
          <div className="bg-muted flex flex-col gap-3 rounded-md p-3">
            <p className="text-sm font-medium">
              เปลี่ยนบทบาทของ {props.nickname} เป็น{nextLabel}?
            </p>
            <div className="flex gap-2">
              <Button
                size="lg"
                className="flex-1"
                disabled={setRole.isPending}
                onClick={() =>
                  setRole.mutate(next, {
                    onSettled: () => setConfirming(false),
                  })
                }
              >
                ยืนยัน
              </Button>
              <Button
                variant="outline"
                size="lg"
                className="flex-1"
                disabled={setRole.isPending}
                onClick={() => setConfirming(false)}
              >
                ยกเลิก
              </Button>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
