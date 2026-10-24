import { buttonVariants } from "@mmtc34/ui/components/button";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

import { requireApproved } from "@/api/me";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/admin")({
  beforeLoad: async ({ context }) => {
    await requireApproved(context.queryClient, { admin: true });
  },
  component: AdminLayout,
});

function AdminLayout() {
  return (
    <div className="mx-auto flex min-h-svh max-w-3xl flex-col">
      <header className="flex items-center justify-between gap-2 border-b p-3">
        {/* nav.tsx renders a bare flex row; scoped variants keep labels on one line and the row scrollable */}
        <div className="min-w-0 overflow-x-auto [&_a]:whitespace-nowrap [&_nav]:w-max [&_nav]:gap-3">
          <NavSlot slot="admin" />
        </div>
        <Link
          to="/"
          className={`${buttonVariants({ variant: "outline" })} shrink-0 whitespace-nowrap`}
        >
          กลับหน้าสมาชิก
        </Link>
      </header>
      <main className="flex-1 p-4">
        <Outlet />
      </main>
    </div>
  );
}
