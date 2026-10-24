import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@mmtc34/ui/components/card";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { membersQuery } from "@/api/members";

export const Route = createFileRoute("/admin/")({
  staticData: { nav: { slot: "admin", label: "ภาพรวม", order: 0 } },
  component: AdminHome,
});

function AdminHome() {
  const { data: applicants } = useQuery(membersQuery("applicant"));
  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">จัดการระบบ</h1>
      <Card>
        <CardHeader>
          <CardTitle>ใบสมัครรออนุมัติ</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-between">
          <span className="text-2xl font-semibold">
            {applicants?.length ?? "-"}
          </span>
          <Link
            to="/admin/members"
            search={{ status: "applicant" }}
            className="text-sm underline"
          >
            ดูรายการ
          </Link>
        </CardContent>
      </Card>
    </div>
  );
}
