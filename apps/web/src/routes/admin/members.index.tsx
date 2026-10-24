import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";
import { errorText } from "@/api/me";
import {
  type MemberFilter,
  membersQuery,
  useApproveMember,
  useRejectMember,
} from "@/api/members";

const tabs: { status: MemberFilter; label: string }[] = [
  { status: "applicant", label: "รออนุมัติ" },
  { status: "approved", label: "สมาชิก" },
];

export const Route = createFileRoute("/admin/members/")({
  staticData: { nav: { slot: "admin", label: "สมาชิก", order: 10 } },
  validateSearch: (search: { status?: string }) => ({
    status: (search.status === "approved"
      ? "approved"
      : "applicant") as MemberFilter,
  }),
  loaderDeps: ({ search }) => ({ status: search.status }),
  loader: ({ context, deps }) =>
    context.queryClient.ensureQueryData(membersQuery(deps.status)),
  component: Members,
});

function Members() {
  const { status } = Route.useSearch();
  const { data: members } = useSuspenseQuery(membersQuery(status));
  const approve = useApproveMember();
  const reject = useRejectMember();
  const error = approve.error ?? reject.error;

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">สมาชิก</h1>
      <div className="flex gap-2">
        {tabs.map((tab) => (
          <Link
            key={tab.status}
            to="/admin/members"
            search={{ status: tab.status }}
            className={`rounded-none border px-3 py-1.5 text-sm ${
              tab.status === status ? "bg-primary text-primary-foreground" : ""
            }`}
          >
            {tab.label}
          </Link>
        ))}
      </div>
      {error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(error)}
        </p>
      )}
      {members.length === 0 && (
        <p className="text-muted-foreground text-sm">ไม่มีรายการ</p>
      )}
      <ul className="flex flex-col gap-2">
        {members.map((m) => (
          <li key={m.id}>
            <Card>
              <CardContent className="flex items-center justify-between gap-2">
                <Link
                  to="/admin/members/$memberId"
                  params={{ memberId: String(m.id) }}
                  className="min-w-0 flex-1"
                >
                  <div className="font-medium">{m.nickname}</div>
                  <div className="text-muted-foreground truncate text-sm">
                    {m.realName} · รุ่น {m.cohort}
                  </div>
                </Link>
                {status === "applicant" && (
                  <div className="flex gap-2">
                    <Button
                      size="sm"
                      disabled={approve.isPending || reject.isPending}
                      onClick={() => approve.mutate(m.id)}
                    >
                      อนุมัติ
                    </Button>
                    <Button
                      size="sm"
                      variant="destructive"
                      disabled={approve.isPending || reject.isPending}
                      onClick={() => {
                        if (confirm(`ปฏิเสธใบสมัครของ ${m.nickname}?`))
                          reject.mutate(m.id);
                      }}
                    >
                      ปฏิเสธ
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>
          </li>
        ))}
      </ul>
    </div>
  );
}
