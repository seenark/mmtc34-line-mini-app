import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "./client";
import type {
  Member,
  OrderLinkInfo,
  OrderLinkInput,
  OrderLinkResult,
  ProfileInput,
} from "./types.gen";

/** Adds an Unlinked Member (nickname only); the caller decides where to go next. */
export function useCreateUnlinkedMember() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: ProfileInput) =>
      api<Member>("POST", "/admin/members", input),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin", "members"] }),
  });
}

/** Who owns the Order and how many of their other Orders a link could move. 409 once the owner has LINE. */
export const orderLinkInfoQuery = (orderId: string) =>
  queryOptions({
    queryKey: ["admin", "order-link", orderId],
    queryFn: () => api<OrderLinkInfo>("GET", `/admin/orders/${orderId}/link`),
  });

/**
 * Moves the Order (or all of the Unlinked Member's Orders) to a Member who has a LINE account.
 * `onLinked` runs before the invalidations: their refetch unmounts the dialog, and a mutate-level
 * callback does not fire once its component is gone.
 */
export function useLinkOrder(
  orderId: string,
  onLinked: (result: OrderLinkResult) => void,
) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: OrderLinkInput) =>
      api<OrderLinkResult>("POST", `/admin/orders/${orderId}/link`, input),
    onSuccess: async (result) => {
      onLinked(result);
      await Promise.all([
        qc.invalidateQueries({ queryKey: ["orders"] }),
        qc.invalidateQueries({ queryKey: ["admin"] }),
      ]);
    },
  });
}
