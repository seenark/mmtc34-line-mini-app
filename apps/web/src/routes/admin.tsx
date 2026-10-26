import { Button } from "@mmtc34/ui/components/button";
import { createFileRoute, Link, Outlet } from "@tanstack/react-router";

import { requireMember, useLogout } from "@/api/me";
import { Brand } from "@/components/brand";
import { ModeSwitch } from "@/components/mode-switch";
import { canSwitchMode, useRecordMode } from "@/lib/admin-mode";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/admin")({
  beforeLoad: async ({ context }) => ({
    member: await requireMember(context.queryClient, { admin: true }),
  }),
  component: AdminLayout,
});

function AdminLayout() {
  const { member } = Route.useRouteContext();
  const logout = useLogout();
  useRecordMode("admin");
  return (
    <div className="flex min-h-svh flex-col">
      <header className="bg-primary text-primary-foreground sticky top-0 z-20 pt-[env(safe-area-inset-top)]">
        <div className="mx-auto w-full max-w-3xl">
          <div className="flex h-14 items-center justify-between gap-3 px-4">
            <Link to="/admin" className="min-w-0 rounded-md">
              <Brand />
            </Link>
            {canSwitchMode(member) ? (
              <ModeSwitch current="admin" onNavy />
            ) : (
              <Button
                variant="ghost"
                className="text-primary-foreground hover:bg-primary-raised hover:text-primary-foreground"
                disabled={logout.isPending}
                onClick={() => logout.mutate()}
              >
                ออกจากระบบ
              </Button>
            )}
          </div>
          <div className="[scrollbar-width:none] overflow-x-auto px-2">
            <NavSlot slot="admin" />
          </div>
        </div>
      </header>
      <main className="mx-auto w-full max-w-3xl min-w-0 flex-1 p-4">
        <Outlet />
      </main>
    </div>
  );
}
