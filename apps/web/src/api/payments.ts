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

/** Shrinks every picked image and sends them as repeated `file` fields. */
async function slipForm(files: File[]) {
  const form = new FormData();
  for (const file of await Promise.all(files.map(shrinkImage))) {
    form.append("file", file, "slip.jpg");
  }
  return form;
}

export async function uploadSlip(orderId: string, files: File[]) {
  return api<Order>("POST", `/orders/${orderId}/slip`, await slipForm(files));
}

export function useUploadSlip(orderId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (files: File[]) => uploadSlip(orderId, files),
    onSuccess: (order) => settle(qc, order),
  });
}

/** Admin: add images to an existing Slip Payment. */
export function useAddSlipImages(orderId: string, paymentId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (files: File[]) =>
      api<Order>(
        "POST",
        `/admin/orders/${orderId}/payments/${paymentId}/slips`,
        await slipForm(files),
      ),
    onSuccess: (order) => settle(qc, order),
  });
}

/** Admin: remove one image; the server refuses the Payment's last one. */
export function useRemoveSlipImage(orderId: string, paymentId: number) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (name: string) =>
      api<Order>(
        "DELETE",
        `/admin/orders/${orderId}/payments/${paymentId}/slips/${name}`,
      ),
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
