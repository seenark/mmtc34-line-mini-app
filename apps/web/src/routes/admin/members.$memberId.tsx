import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { memberQuery } from "@/api/members";
import { MemberApproved, RoleAdmin } from "@/api/types.gen";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/admin/members/$memberId")({
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(memberQuery(params.memberId)),
  component: MemberDetail,
});

function MemberDetail() {
  const { memberId } = Route.useParams();
  const { data: m } = useSuspenseQuery(memberQuery(memberId));
  const rows: [string, string][] = [
    ["สถานะ", m.status === MemberApproved ? "สมาชิก" : "รออนุมัติ"],
    ["บทบาท", m.role === RoleAdmin ? "แอดมิน" : "สมาชิก"],
    ["รุ่น", m.cohort],
    ["ชื่อเล่น", m.nickname],
    ["ชื่อจริง", m.realName],
    ["ที่อยู่", m.address || "-"],
    ["เบอร์โทร", m.phone || "-"],
  ];
  return (
    <div className="flex flex-col gap-4">
      <Link
        to="/admin/members"
        search={{
          status: m.status === MemberApproved ? "approved" : "applicant",
        }}
        className="text-sm underline"
      >
        ← สมาชิก
      </Link>
      <h1 className="text-2xl font-semibold">{m.nickname}</h1>
      <Card>
        <CardContent>
          <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-2 text-sm">
            {rows.map(([label, value]) => (
              <div key={label} className="contents">
                <dt className="text-muted-foreground">{label}</dt>
                <dd className="break-words whitespace-pre-wrap">{value}</dd>
              </div>
            ))}
          </dl>
        </CardContent>
      </Card>
      <NavSlot slot="adminMember" params={{ memberId }} />
    </div>
  );
}
