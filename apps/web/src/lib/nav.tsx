import { type UseQueryOptions, useQuery } from "@tanstack/react-query";
import { type AnyRoute, Link, useRouter } from "@tanstack/react-router";

import type { OrderActions } from "@/api/types.gen";

// Navigation is derived from the routes themselves: a route file declares
//   staticData: { nav: { slot: "admin", label: "ตรวจสลิป", order: 20, badge: pendingSlipCount } }
// and the layout for that slot renders <NavSlot slot="admin" />. No hand-kept menu list exists.

export type NavSlotName =
  | "member" // bottom bar of the approved-Member layout (_app.tsx)
  | "admin" // admin layout sidebar (admin.tsx)
  | "memberCampaign" // tabs on _app/campaigns.$campaignId.tsx
  | "adminCampaign" // tabs on admin/campaigns.$campaignId.tsx
  | "adminMember" // actions on admin/members.$memberId.tsx
  | "orderAction"; // buttons on the Order page

export interface NavEntry {
  slot: NavSlotName;
  label: string; // Thai
  order: number;
  badge?: UseQueryOptions<number>; // e.g. pending Slip count
  /** orderAction only: shown when order.actions[when] is true (the server's state machine decides). */
  when?: keyof OrderActions;
}

declare module "@tanstack/react-router" {
  interface StaticDataRouteOption {
    nav?: NavEntry;
  }
}

/** Links to every route whose staticData.nav.slot === slot, sorted by order. */
export function NavSlot(props: {
  slot: NavSlotName;
  params?: Record<string, string>;
  actions?: OrderActions;
}) {
  const { routesById } = useRouter();
  const routes: AnyRoute[] = Object.values(routesById);
  const entries = routes
    .flatMap((route) => {
      const nav = route.options.staticData?.nav;
      return nav?.slot === props.slot &&
        (!nav.when || props.actions?.[nav.when])
        ? [{ to: route.fullPath, nav }]
        : [];
    })
    .sort((a, b) => a.nav.order - b.nav.order);

  return (
    <nav className="flex gap-2">
      {entries.map(({ to, nav }) => (
        // entries come from the runtime route tree, so `to` and `params` cannot be statically typed
        <Link key={to} to={to as never} params={props.params as never}>
          {nav.label}
          {nav.badge && <Badge options={nav.badge} />}
        </Link>
      ))}
    </nav>
  );
}

function Badge({ options }: { options: UseQueryOptions<number> }) {
  const { data } = useQuery(options);
  return data ? (
    <span className="ml-1 rounded-full bg-red-600 px-1.5 text-xs text-white">
      {data}
    </span>
  ) : null;
}
