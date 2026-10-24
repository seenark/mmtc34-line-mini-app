import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "./client";
import type { Member } from "./types.gen";

export type MemberFilter = "applicant" | "approved";

export const membersQuery = (status: MemberFilter) =>
  queryOptions({
    queryKey: ["admin", "members", { status }],
    queryFn: () => api<Member[]>("GET", `/admin/members?status=${status}`),
  });

export const memberQuery = (memberId: string) =>
  queryOptions({
    queryKey: ["admin", "members", memberId],
    queryFn: () => api<Member>("GET", `/admin/members/${memberId}`),
  });

function useReviewMutation<T>(action: "approve" | "reject") {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (memberId: number) =>
      api<T>("POST", `/admin/members/${memberId}/${action}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["admin", "members"] }),
  });
}

export const useApproveMember = () => useReviewMutation<Member>("approve");
export const useRejectMember = () => useReviewMutation<void>("reject");
