import {
  type QueryClient,
  type QueryKey,
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

import { shrinkImage } from "@/lib/image";

import { api } from "./client";
import { orderQuery, ordersQuery } from "./orders";
import type { AmountInput, Order, RejectInput } from "./types.gen";

/** Absolute-path URL of an uploaded Slip; the server allows only its owner and Admins. */
export const slipUrl = (name: string) =>
  `${import.meta.env.BASE_URL}api/slips/${name}`;

export const pendingPaymentsQuery = queryOptions({
  queryKey: ["payments", "pending"],
  queryFn: () => api<Order[]>("GET", "/admin/payments/pending"),
});

/** Nav badge: how many Slips wait for review. */
export const pendingSlipCount = queryOptions({
  queryKey: ["payments", "pending", "count"] as QueryKey,
  queryFn: async () =>
    (await api<Order[]>("GET", "/admin/payments/pending")).length,
  refetchInterval: 30_000,
});

/** Every Slip mutation returns the updated Order; the Slip queue and "my orders" must refetch. */
function settle(qc: QueryClient, order: Order) {
  qc.setQueryData(orderQuery(String(order.id)).queryKey, order);
  return Promise.all([
    qc.invalidateQueries({ queryKey: pendingPaymentsQuery.queryKey }),
    qc.invalidateQueries({ queryKey: ordersQuery.queryKey }),
  ]);
}

export async function uploadSlip(orderId: string, file: File) {
  const form = new FormData();
  form.append("file", await shrinkImage(file), "slip.jpg");
  return api<Order>("POST", `/orders/${orderId}/slip`, form);
}

export function useUploadSlip(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => uploadSlip(orderId, file),
    onSuccess: (order) => settle(qc, order),
  });
}

export function useConfirmSlip(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (amount: number) =>
      api<Order>("POST", `/admin/orders/${orderId}/confirm`, {
        amount,
      } satisfies AmountInput),
    onSuccess: (order) => settle(qc, order),
  });
}

export function useRejectSlip(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (reason: string) =>
      api<Order>("POST", `/admin/orders/${orderId}/reject`, {
        reason,
      } satisfies RejectInput),
    onSuccess: (order) => settle(qc, order),
  });
}
