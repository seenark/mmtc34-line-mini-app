import {
  createFileRoute,
  Link,
  Outlet,
  redirect,
} from "@tanstack/react-router";

import { requireMember } from "@/api/me";
import { Brand } from "@/components/brand";
import { ModeSwitch } from "@/components/mode-switch";
import { canSwitchMode, launchMode, useRecordMode } from "@/lib/admin-mode";
import { NavSlot } from "@/lib/nav";

export const Route = createFileRoute("/_app")({
  beforeLoad: async ({ context, location }) => {
    const member = await requireMember(context.queryClient);
    // The Super Admin has no Orders of their own, so there is no member side to land on.
    if (member.superAdmin) throw redirect({ to: "/admin" });
    if (launchMode(member) === "admin" && location.pathname === "/")
      throw redirect({ to: "/admin" });
    return { member };
  },
  component: MemberLayout,
});

function MemberLayout() {
  const { member } = Route.useRouteContext();
  useRecordMode("member");
  return (
    <div className="bg-background mx-auto flex min-h-svh w-full max-w-md flex-col md:border-x">
      <header className="bg-background sticky top-0 z-20 border-b pt-[env(safe-area-inset-top)]">
        <div className="flex h-14 items-center justify-between gap-3 px-4">
          <Link to="/" className="text-primary min-w-0 rounded-md">
            <Brand />
          </Link>
          {canSwitchMode(member) && <ModeSwitch current="member" />}
        </div>
      </header>
      <main className="min-w-0 flex-1 p-4">
        <Outlet />
      </main>
      <footer className="bg-background sticky bottom-0 z-20 border-t pb-[env(safe-area-inset-bottom)]">
        <NavSlot slot="member" />
      </footer>
    </div>
  );
}
