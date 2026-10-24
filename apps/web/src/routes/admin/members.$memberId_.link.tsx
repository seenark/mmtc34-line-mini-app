import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { errorText } from "@/api/me";
import { memberQuery, membersQuery } from "@/api/members";
import { useLinkMember } from "@/api/unlinked";

// The `$memberId_` escape keeps this out of members.$memberId, whose page renders no <Outlet />.
export const Route = createFileRoute("/admin/members/$memberId_/link")({
  staticData: {
    nav: { slot: "adminMember", label: "ผูกบัญชี LINE", order: 20 },
  },
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(memberQuery(params.memberId)),
      context.queryClient.ensureQueryData(membersQuery("applicant")),
    ]),
  component: LinkMember,
});

function LinkMember() {
  const { memberId } = Route.useParams();
  const { data: m } = useSuspenseQuery(memberQuery(memberId));
  const { data: applicants } = useSuspenseQuery(membersQuery("applicant"));
  const link = useLinkMember(memberId);

  return (
    <div className="flex flex-col gap-4">
      <Link
        to="/admin/members/$memberId"
        params={{ memberId }}
        className="text-sm underline"
      >
        ← {m.nickname}
      </Link>
      <h1 className="text-2xl font-semibold">ผูกบัญชี LINE</h1>
      {m.linked ? (
        <p role="status" className="bg-muted p-3 text-sm font-medium">
          บัญชีนี้ผูก LINE แล้ว
        </p>
      ) : (
        <>
          <p className="text-muted-foreground text-sm">
            เลือกผู้สมัครที่ใช้ LINE ของ {m.nickname}
          </p>
          {link.error && (
            <p role="alert" className="text-destructive text-sm">
              {errorText(link.error)}
            </p>
          )}
          {applicants.length === 0 && (
            <p className="text-muted-foreground text-sm">ไม่มีผู้สมัคร</p>
          )}
          <ul className="flex flex-col gap-2">
            {applicants.map((a) => (
              <li key={a.id}>
                <Card>
                  <CardContent className="flex items-center justify-between gap-2">
                    <div className="min-w-0">
                      <div className="font-medium">{a.nickname}</div>
                      <div className="text-muted-foreground truncate text-sm">
                        {a.realName} · รุ่น {a.cohort}
                      </div>
                    </div>
                    <Button
                      size="sm"
                      disabled={link.isPending}
                      onClick={() => {
                        if (
                          confirm(`ผูก LINE ของ ${a.nickname} กับ ${m.nickname}?`)
                        )
                          link.mutate({ applicantId: a.id });
                      }}
                    >
                      ผูกบัญชีนี้
                    </Button>
                  </CardContent>
                </Card>
              </li>
            ))}
          </ul>
        </>
      )}
    </div>
  );
}
