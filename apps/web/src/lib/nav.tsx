import { buttonVariants } from "@mmtc34/ui/components/button";
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
  | "adminOrder"; // buttons on the admin Order page

export interface NavEntry {
  slot: NavSlotName;
  label: string; // Thai
  order: number;
  badge?: UseQueryOptions<number>; // e.g. pending Slip count
  /** adminOrder only: shown when order.actions[when] is true (the server's state machine decides). */
  when?: keyof OrderActions;
}

declare module "@tanstack/react-router" {
  interface StaticDataRouteOption {
    nav?: NavEntry;
  }
}

const TAB =
  "flex h-11 items-center border-b-2 border-transparent px-3 text-sm font-semibold whitespace-nowrap -outline-offset-2 transition-colors";
const CAMPAIGN_TABS = {
  nav: "flex gap-1 overflow-x-auto border-b",
  link: `${TAB} -mb-px text-muted-foreground hover:text-foreground data-[status=active]:border-secondary data-[status=active]:text-secondary`,
};
const ACTIONS = {
  nav: "flex flex-wrap gap-2",
  link: buttonVariants({ variant: "outline" }),
};

// Labels stay on one line; a row that does not fit scrolls in its container instead of wrapping.
const SLOT_STYLES: Record<NavSlotName, { nav: string; link: string }> = {
  member: {
    nav: "grid auto-cols-fr grid-flow-col",
    link: "text-muted-foreground hover:text-foreground data-[status=active]:border-secondary data-[status=active]:text-secondary flex min-h-14 items-center justify-center border-t-2 border-transparent px-2 text-sm font-semibold whitespace-nowrap -outline-offset-2 transition-colors",
  },
  admin: {
    nav: "flex w-max gap-1",
    link: `${TAB} text-primary-muted-foreground hover:text-primary-foreground data-[status=active]:border-primary-foreground data-[status=active]:text-primary-foreground`,
  },
  memberCampaign: CAMPAIGN_TABS,
  adminCampaign: CAMPAIGN_TABS,
  adminMember: ACTIONS,
  adminOrder: ACTIONS,
};

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

  const style = SLOT_STYLES[props.slot];
  return (
    <nav className={style.nav}>
      {entries.map(({ to, nav }) => (
        // entries come from the runtime route tree, so `to` and `params` cannot be statically typed
        <Link
          key={to}
          to={to as never}
          params={props.params as never}
          // a slot's root entry (e.g. /admin/) must not stay active on its siblings
          activeOptions={{
            exact: entries.some((o) => o.to !== to && o.to.startsWith(to)),
          }}
          className={style.link}
        >
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
    <span className="bg-destructive text-primary-foreground ml-1.5 rounded-full px-1.5 text-xs font-semibold">
      {data}
    </span>
  ) : null;
}
