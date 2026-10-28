import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { api } from "./client";
import { myAddressesQuery } from "./profile";
import type { Order, OrderDetailsInput, OrderInput } from "./types.gen";

export const ordersQuery = queryOptions({
  queryKey: ["orders"],
  queryFn: () => api<Order[]>("GET", "/orders"),
});

export const orderQuery = (orderId: string) =>
  queryOptions({
    queryKey: ["orders", orderId],
    queryFn: () => api<Order>("GET", `/orders/${orderId}`),
  });

/** `onPlaced` runs before navigating: the page unmounts, so it cannot be a mutate-level callback. */
export function useCheckout(campaignId: string, onPlaced: () => void) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: OrderInput) =>
      api<Order>("POST", `/campaigns/${campaignId}/orders`, input),
    onSuccess: async (order) => {
      onPlaced();
      qc.setQueryData(orderQuery(String(order.id)).queryKey, order);
      await Promise.all([
        qc.invalidateQueries({ queryKey: ordersQuery.queryKey }),
        qc.invalidateQueries({ queryKey: myAddressesQuery.queryKey }),
      ]);
      await navigate({
        to: "/orders/$orderId",
        params: { orderId: String(order.id) },
      });
    },
  });
}

export function useCancelOrder(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: () => api<Order>("POST", `/orders/${orderId}/cancel`),
    onSuccess: (order) => {
      qc.setQueryData(orderQuery(orderId).queryKey, order);
      return qc.invalidateQueries({ queryKey: ordersQuery.queryKey });
    },
  });
}

export function useEditOrderDetails(orderId: string) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: OrderDetailsInput) =>
      api<Order>("PUT", `/orders/${orderId}/details`, input),
    onSuccess: async (order) => {
      qc.setQueryData(orderQuery(orderId).queryKey, order);
      await Promise.all([
        qc.invalidateQueries({ queryKey: ordersQuery.queryKey }),
        qc.invalidateQueries({ queryKey: myAddressesQuery.queryKey }),
      ]);
      await navigate({ to: "/orders/$orderId", params: { orderId } });
    },
  });
}
