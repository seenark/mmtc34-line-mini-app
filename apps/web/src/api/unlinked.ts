import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { api } from "./client";
import type { LinkInput, Member, ProfileInput } from "./types.gen";

/** Both writes change Member rows (a new one, or an Applicant folded into an Unlinked Member). */
function useMemberWrite<V>(
  send: (v: V) => Promise<Member>,
  memberId: (v: V, m: Member) => string,
) {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: send,
    onSuccess: async (m, v) => {
      await qc.invalidateQueries({ queryKey: ["admin", "members"] });
      await navigate({
        to: "/admin/members/$memberId",
        params: { memberId: memberId(v, m) },
      });
    },
  });
}

export const useCreateUnlinkedMember = () =>
  useMemberWrite<ProfileInput>(
    (input) => api<Member>("POST", "/admin/members", input),
    (_, m) => String(m.id),
  );

export const useLinkMember = (memberId: string) =>
  useMemberWrite<LinkInput>(
    (input) => api<Member>("POST", `/admin/members/${memberId}/link`, input),
    () => memberId,
  );
