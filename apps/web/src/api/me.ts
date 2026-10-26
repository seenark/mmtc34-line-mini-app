import {
  type QueryClient,
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { redirect, useNavigate } from "@tanstack/react-router";

import { clearAllCarts } from "@/lib/cart";
import { LineSignInError } from "@/lib/liff";

import { ApiError, api } from "./client";
import {
  type AdminLoginInput,
  type Me,
  type Member,
  MemberApproved,
  type ProfileInput,
  RoleAdmin,
} from "./types.gen";

/** Server errors and LINE sign-in errors are already Thai; anything else gets a generic Thai message. */
export function errorText(e: unknown): string {
  if (e instanceof ApiError || e instanceof LineSignInError) return e.message;
  return "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง";
}

/** null = no valid session. */
export const meQuery = queryOptions({
  queryKey: ["me"],
  staleTime: 10_000,
  queryFn: async (): Promise<Me | null> => {
    try {
      return await api<Me>("GET", "/me");
    } catch (e) {
      if (e instanceof ApiError && e.status === 401) return null;
      throw e;
    }
  },
});

type Stage = "/login" | "/onboarding" | "/waiting" | "/";

/** The one screen a viewer may be on: no session, not onboarded, Applicant, or approved Member. */
function stageOf(me: Me | null): Stage {
  if (!me) return "/login";
  if (!me.member) return "/onboarding";
  return me.member.status === MemberApproved ? "/" : "/waiting";
}

/** Redirects to the viewer's own stage unless it is `expected`. */
export async function enforceStage(qc: QueryClient, expected: Stage) {
  const stage = stageOf(await qc.fetchQuery(meQuery));
  if (stage !== expected) throw redirect({ to: stage });
}

/** Gate for the approved-Member and Admin layouts. Returns the fresh Member. */
export async function requireApproved(
  qc: QueryClient,
  opts?: { admin?: boolean },
): Promise<Member> {
  const me = await qc.fetchQuery(meQuery);
  const stage = stageOf(me);
  if (stage !== "/" || !me?.member) throw redirect({ to: stage });
  if (opts?.admin && me.member.role !== RoleAdmin) throw redirect({ to: "/" });
  return me.member;
}

/** Super Admin login (server SUPER_ADMIN_PASSWORD set): username + password, no LINE. */
export async function adminSignIn(username: string, password: string) {
  await api<void>("POST", "/auth/admin", {
    username,
    password,
  } satisfies AdminLoginInput);
}

export function useOnboard() {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: ProfileInput) => api<Member>("POST", "/me", input),
    onSuccess: async () => {
      await qc.resetQueries({ queryKey: meQuery.queryKey });
      await navigate({ to: "/" });
    },
  });
}

export function useLogout() {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: () => api<void>("POST", "/auth/logout"),
    onSuccess: async () => {
      clearAllCarts();
      qc.clear();
      await navigate({ to: "/login" });
    },
  });
}
