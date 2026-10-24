import { queryOptions } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { api } from "../../api/client";
import type { Order } from "../../api/types.gen";

// Ticket .5. Adding this file is the whole "add an admin screen + badge" change:
// routeTree.gen.ts picks the route up and <NavSlot slot="admin"> picks the nav entry up.

export const pendingSlips = queryOptions({
  queryKey: ["admin", "payments", "pending"],
  queryFn: () => api<Order[]>("GET", "/admin/payments/pending"),
});

const pendingSlipCount = { ...pendingSlips, select: (orders: Order[]) => orders.length };

export const Route = createFileRoute("/admin/payments")({
  staticData: { nav: { slot: "admin", label: "ตรวจสลิป", order: 20, badge: pendingSlipCount } },
  loader: ({ context }) => context.queryClient.ensureQueryData(pendingSlips),
  component: SlipQueue,
});

/** Slip image, Campaign, nickname, Total; confirm (amount prefilled to total) or reject (reason). */
function SlipQueue(): JSX.Element {
  throw new Error("not implemented");
}
