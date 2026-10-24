import type { UseQueryOptions } from "@tanstack/react-query";
import type { OrderActions } from "../api/types.gen";

// Navigation is derived from the routes themselves: a route file declares
//   export const Route = createFileRoute("/admin/payments")({
//     staticData: { nav: { slot: "admin", label: "ตรวจสลิป", order: 20, badge: pendingSlipCount } },
//   })
// and the layout for that slot renders <NavSlot slot="admin" />. No hand-kept menu list exists,
// so parallel tickets add screens and tabs without editing a shared file.

export type NavSlotName =
  | "member" // bottom bar of the approved-Member layout (_app.tsx)
  | "admin" // admin layout sidebar (admin.tsx)
  | "memberCampaign" // tabs on _app/campaigns.$campaignId.tsx: รายละเอียด, สปอนเซอร์, รายชื่อ
  | "adminCampaign" // tabs on admin/campaigns.$campaignId.tsx: แก้ไข, ออเดอร์, สั่งแทน, สปอนเซอร์, ส่งออก
  | "adminMember" // actions on admin/members.$memberId.tsx: แก้ไข, เชื่อมบัญชี
  | "orderAction"; // buttons on the Order page: อัปโหลดสลิป, แก้ไข, รับเงินสด, ปรับยอด

export interface NavEntry {
  slot: NavSlotName;
  label: string; // Thai
  order: number;
  badge?: UseQueryOptions<number>; // e.g. pending Slip count
  /** orderAction only: shown when order.actions[when] is true (server state machine decides). */
  when?: keyof OrderActions;
}

declare module "@tanstack/react-router" {
  interface StaticDataRouteOption {
    nav?: NavEntry;
  }
}

/** Links to every route whose staticData.nav.slot === slot, sorted by order; params pass through. */
export function NavSlot(props: {
  slot: NavSlotName;
  params?: Record<string, string>;
  actions?: OrderActions;
}): JSX.Element {
  throw new Error("not implemented");
}
