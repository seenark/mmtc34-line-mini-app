import {
  type QueryClient,
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { api } from "./client";
import { orderQuery, ordersQuery } from "./orders";
import type { AdjustmentInput, AmountInput, Order } from "./types.gen";

/** `status` is an effective Order status; undefined lists every Order. */
export const campaignOrdersQuery = (campaignId: string, status?: string) =>
  queryOptions({
    queryKey: ["admin", "campaigns", campaignId, "orders", status ?? "all"],
    queryFn: () =>
      api<Order[]>(
        "GET",
        `/admin/campaigns/${campaignId}/orders${status ? `?status=${status}` : ""}`,
      ),
  });

/** Money moved: the Order page, the Member's list and every Admin Order list are stale. */
async function settle(qc: QueryClient, order: Order) {
  qc.setQueryData(orderQuery(String(order.id)).queryKey, order);
  await Promise.all([
    qc.invalidateQueries({ queryKey: ordersQuery.queryKey }),
    qc.invalidateQueries({ queryKey: ["admin"] }),
  ]);
}

function useLedgerMutation<T>(
  orderId: string,
  path: string,
  body: (input: T) => unknown,
) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: T) =>
      api<Order>("POST", `/admin/orders/${orderId}/${path}`, body(input)),
    onSuccess: async (order) => {
      await settle(qc, order);
      await navigate({ to: "/orders/$orderId", params: { orderId } });
    },
  });
}

export const useRecordCash = (orderId: string) =>
  useLedgerMutation<number>(
    orderId,
    "cash",
    (amount) => ({ amount }) satisfies AmountInput,
  );

export const useAddAdjustment = (orderId: string) =>
  useLedgerMutation<AdjustmentInput>(orderId, "adjustments", (input) => input);
