import { queryOptions } from "@tanstack/react-query";

import { api } from "./client";
import type { Member } from "./types.gen";

export const membersQuery = queryOptions({
  queryKey: ["admin", "members"],
  queryFn: () => api<Member[]>("GET", "/admin/members"),
});

export const memberQuery = (memberId: string) =>
  queryOptions({
    queryKey: ["admin", "members", memberId],
    queryFn: () => api<Member>("GET", `/admin/members/${memberId}`),
  });
