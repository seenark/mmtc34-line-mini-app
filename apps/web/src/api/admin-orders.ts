import {
  type QueryClient,
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { api } from "./client";
import { orderQuery } from "./orders";
import { myAddressesQuery } from "./profile";
import type {
  AddressChoice,
  AdminOrderInput,
  Order,
  OrderInput,
} from "./types.gen";

/** The Member's shipping triples for an Order placed or edited on their behalf. */
export const memberAddressesQuery = (memberId: number) =>
  queryOptions({
    queryKey: ["admin", "members", String(memberId), "addresses"],
    queryFn: () =>
      api<AddressChoice[]>("GET", `/admin/members/${memberId}/addresses`),
  });

/** An Admin write can change the Member's Order list, the Admin Order lists and the Campaign totals. */
async function settle(qc: QueryClient, order: Order) {
  qc.setQueryData(orderQuery(String(order.id)).queryKey, order);
  await Promise.all([
    qc.invalidateQueries({ queryKey: ["orders"] }),
    qc.invalidateQueries({ queryKey: myAddressesQuery.queryKey }),
    qc.invalidateQueries({ queryKey: ["admin"] }),
  ]);
}

export function useCreateOrderOnBehalf(campaignId: string) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: AdminOrderInput) =>
      api<Order>("POST", `/admin/campaigns/${campaignId}/orders`, input),
    onSuccess: async (order) => {
      await settle(qc, order);
      await navigate({
        to: "/admin/orders/$orderId",
        params: { orderId: String(order.id) },
      });
    },
  });
}

export function useEditOrder(orderId: string) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: OrderInput) =>
      api<Order>("PUT", `/admin/orders/${orderId}`, input),
    onSuccess: async (order) => {
      await settle(qc, order);
      await navigate({ to: "/admin/orders/$orderId", params: { orderId } });
    },
  });
}

export function useCancelOrderAsAdmin(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api<Order>("POST", `/orders/${orderId}/cancel`),
    onSuccess: (order) => settle(qc, order),
  });
}
