import {
  type QueryClient,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { api } from "./client";
import { orderQuery } from "./orders";
import type { AdminOrderInput, Order, OrderInput } from "./types.gen";

/** An Admin write can change the Member's Order list, the Admin Order lists and the Campaign totals. */
async function settle(qc: QueryClient, order: Order) {
  qc.setQueryData(orderQuery(String(order.id)).queryKey, order);
  await Promise.all([
    qc.invalidateQueries({ queryKey: ["orders"] }),
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
        to: "/orders/$orderId",
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
      await navigate({ to: "/orders/$orderId", params: { orderId } });
    },
  });
}
