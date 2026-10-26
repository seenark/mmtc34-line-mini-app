import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { clearAllCarts } from "@/lib/cart";

import { api } from "./client";
import { meQuery } from "./me";
import type { AddressChoice, Member, ProfileInput, Role } from "./types.gen";

/** Shipping triples to offer at checkout: the profile, then earlier Orders. Refetched after an Order is placed. */
export const myAddressesQuery = queryOptions({
  queryKey: ["me", "addresses"],
  queryFn: () => api<AddressChoice[]>("GET", "/me/addresses"),
});

export function useUpdateProfile() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: ProfileInput) => api<Member>("PUT", "/me", input),
    onSuccess: () => qc.invalidateQueries({ queryKey: meQuery.queryKey }),
  });
}

/** The server erases contact data and clears the session cookie; payment records stay. */
export function useDeleteAccount() {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: () => api<void>("DELETE", "/me"),
    onSuccess: async () => {
      clearAllCarts();
      qc.clear();
      await navigate({ to: "/login" });
    },
  });
}

/** An Admin may edit their own Member here too, so `me` is refreshed as well. */
function useAdminMemberMutation<T>(request: (input: T) => Promise<Member>) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: request,
    onSuccess: () =>
      Promise.all([
        qc.invalidateQueries({ queryKey: ["admin", "members"] }),
        qc.invalidateQueries({ queryKey: meQuery.queryKey }),
      ]),
  });
}

export function useAdminEditMember(memberId: string) {
  return useAdminMemberMutation((input: ProfileInput) =>
    api<Member>("PUT", `/admin/members/${memberId}`, input),
  );
}

export function useSetRole(memberId: string) {
  return useAdminMemberMutation((role: Role) =>
    api<Member>("PUT", `/admin/members/${memberId}/role`, { role }),
  );
}
