import { buttonVariants } from "@mmtc34/ui/components/button";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

import { requireApproved } from "@/api/me";
import { RoleAdmin } from "@/api/types.gen";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/_app")({
  beforeLoad: async ({ context }) => ({
    member: await requireApproved(context.queryClient),
  }),
  component: MemberLayout,
});

function MemberLayout() {
  const { member } = Route.useRouteContext();
  return (
    <div className="mx-auto flex min-h-svh max-w-md flex-col">
      <main className="flex-1 p-4">
        <Outlet />
      </main>
      <footer className="bg-background sticky bottom-0 flex items-center justify-between gap-2 border-t p-3">
        <NavSlot slot="member" />
        {member.role === RoleAdmin && (
          <Link to="/admin" className={buttonVariants({ variant: "outline" })}>
            จัดการระบบ
          </Link>
        )}
      </footer>
    </div>
  );
}
