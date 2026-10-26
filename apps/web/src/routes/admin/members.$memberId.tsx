import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useQueries, useQuery, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { adminCampaignsQuery } from "@/api/campaigns";
import { campaignOrdersQuery } from "@/api/ledger";
import { errorText } from "@/api/me";
import { memberQuery } from "@/api/members";
import { RoleAdmin } from "@/api/types.gen";
import { LinkOrderDialog } from "@/components/order/link-order-dialog";
import { baht, OrderStatusBadge } from "@/components/order/order-status";
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
        className="text-secondary w-fit text-sm font-semibold whitespace-nowrap hover:underline"
      >
        ← สมาชิก
      </Link>
      <h1 className="text-primary flex flex-wrap items-center gap-2 text-2xl font-extrabold break-words">
        {m.nickname}
        {m.superAdmin && (
          <span className="bg-muted text-muted-foreground rounded-full px-2 py-0.5 text-xs font-normal whitespace-nowrap">
            ผู้ดูแลระบบหลัก
          </span>
        )}
        {!m.linked && (
          <span className="bg-muted text-muted-foreground rounded-full px-2 py-0.5 text-xs font-normal whitespace-nowrap">
            ไม่มี LINE
          </span>
        )}
      </h1>
      <Card>
        <CardContent>
          <dl className="grid grid-cols-[auto_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm">
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
      {!m.linked && !m.deleted && <UnlinkedOrders memberId={m.id} />}
    </div>
  );
}

/**
 * The API has no per-Member Order list, so the Member's Orders are the Campaign order lists
 * filtered by owner. The Campaign count is small (one class), so this is a handful of requests.
 */
function UnlinkedOrders({ memberId }: { memberId: number }) {
  const campaigns = useQuery(adminCampaignsQuery);
  const lists = useQueries({
    queries: (campaigns.data ?? []).map((c) =>
      campaignOrdersQuery(String(c.id)),
    ),
  });

  const loading = campaigns.isPending || lists.some((l) => l.isPending);
  const error = campaigns.error ?? lists.find((l) => l.error)?.error;
  const orders = lists
    .flatMap((l) => l.data ?? [])
    .filter((o) => o.memberId === memberId)
    .sort((a, b) => b.id - a.id);

  return (
    <section className="flex flex-col gap-2">
      <h2 className="text-base font-semibold">คำสั่งซื้อ</h2>
      {loading && (
        <p className="text-muted-foreground text-sm" role="status">
          กำลังโหลด…
        </p>
      )}
      {error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(error)}
        </p>
      )}
      {!loading && !error && orders.length === 0 && (
        <p className="text-muted-foreground text-sm">ยังไม่มีคำสั่งซื้อ</p>
      )}
      <ul className="flex flex-col gap-2">
        {orders.map((o) => (
          <li key={o.id}>
            <Card>
              <CardContent className="flex flex-col gap-3">
                <Link
                  to="/admin/orders/$orderId"
                  params={{ orderId: String(o.id) }}
                  className="flex min-w-0 flex-col gap-1 rounded-md"
                >
                  <div className="flex items-center justify-between gap-2">
                    <span className="min-w-0 font-semibold break-words">
                      #{o.id} {o.campaignTitle}
                    </span>
                    <OrderStatusBadge status={o.status} />
                  </div>
                  <p className="text-muted-foreground text-sm tabular-nums">
                    ยอดรวม {baht(o.total)}
                  </p>
                </Link>
                <LinkOrderDialog orderId={o.id} />
              </CardContent>
            </Card>
          </li>
        ))}
      </ul>
    </section>
  );
}
