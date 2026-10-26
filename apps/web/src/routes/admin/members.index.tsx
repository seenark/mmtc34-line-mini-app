import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { membersQuery } from "@/api/members";
import { type Member, RoleAdmin } from "@/api/types.gen";

export const Route = createFileRoute("/admin/members/")({
  staticData: { nav: { slot: "admin", label: "สมาชิก", order: 10 } },
  loader: ({ context }) => context.queryClient.ensureQueryData(membersQuery),
  component: Members,
});

const BADGE =
  "shrink-0 rounded-full px-2 py-0.5 text-xs font-semibold whitespace-nowrap";

function RoleBadge({ m }: { m: Member }) {
  if (m.superAdmin) {
    return (
      <span className={`${BADGE} bg-primary text-primary-foreground`}>
        ผู้ดูแลระบบหลัก
      </span>
    );
  }
  if (m.role === RoleAdmin) {
    return (
      <span className={`${BADGE} bg-accent text-accent-foreground`}>แอดมิน</span>
    );
  }
  return null;
}

function Members() {
  const { data: members } = useSuspenseQuery(membersQuery);

  return (
    <div className="flex flex-col gap-5">
      <h1 className="text-primary text-2xl font-extrabold">สมาชิก</h1>
      {members.length === 0 && (
        <p className="text-muted-foreground rounded-xl border border-dashed p-6 text-center text-sm">
          ไม่มีรายการ
        </p>
      )}
      <ul className="flex flex-col gap-2">
        {members.map((m) => {
          const meta = [m.realName, m.cohort && `รุ่น ${m.cohort}`].filter(
            Boolean,
          );
          return (
            <li key={m.id}>
              <Link
                to="/admin/members/$memberId"
                params={{ memberId: String(m.id) }}
                className="group block min-w-0 rounded-xl"
              >
                <Card
                  size="sm"
                  className="group-hover:bg-accent group-active:bg-accent transition-colors"
                >
                  <CardContent className="flex min-w-0 flex-col gap-0.5">
                    <div className="flex min-w-0 items-center gap-2">
                      <span
                        className={`min-w-0 truncate text-sm font-semibold ${m.deleted ? "text-muted-foreground" : ""}`}
                      >
                        {m.nickname}
                      </span>
                      <RoleBadge m={m} />
                      {!m.linked && (
                        <span
                          className={`${BADGE} bg-muted text-muted-foreground`}
                        >
                          ไม่มี LINE
                        </span>
                      )}
                      {m.deleted && (
                        <span
                          className={`${BADGE} bg-muted text-muted-foreground`}
                        >
                          ลบแล้ว
                        </span>
                      )}
                    </div>
                    {meta.length > 0 && (
                      <span className="text-muted-foreground truncate text-sm">
                        {meta.join(" · ")}
                      </span>
                    )}
                  </CardContent>
                </Card>
              </Link>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
