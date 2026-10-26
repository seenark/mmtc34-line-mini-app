import { buttonVariants } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { membersQuery } from "@/api/members";

export const Route = createFileRoute("/admin/")({
  staticData: { nav: { slot: "admin", label: "ภาพรวม", order: 0 } },
  component: AdminHome,
});

function AdminHome() {
  const { data: members } = useQuery(membersQuery);
  return (
    <div className="flex flex-col gap-5">
      <h1 className="text-primary text-2xl font-extrabold">จัดการระบบ</h1>
      <Card>
        <CardContent className="flex items-center justify-between gap-3">
          <div className="flex min-w-0 flex-col">
            <h2 className="text-muted-foreground text-sm font-medium">สมาชิก</h2>
            <span className="text-primary text-3xl font-extrabold tabular-nums">
              {members?.length ?? "-"}
            </span>
          </div>
          <Link
            to="/admin/members"
            className={buttonVariants({ variant: "outline" })}
          >
            ดูรายการ
          </Link>
        </CardContent>
      </Card>
    </div>
  );
}
