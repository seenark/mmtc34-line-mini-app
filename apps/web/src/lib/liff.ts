import liff from "@line/liff";
import { queryOptions } from "@tanstack/react-query";

import { api } from "@/api/client";
import type { Config, DevLoginInput, LineLoginInput } from "@/api/types.gen";

export const configQuery = queryOptions({
  queryKey: ["config"],
  queryFn: () => api<Config>("GET", "/config"),
});

/**
 * LIFF login. Resolves true once the session cookie is set. Resolves false when LIFF is
 * redirecting the page to LINE, so the caller must not navigate.
 */
export async function signIn(config: Config): Promise<boolean> {
  if (!config.liffId) throw new Error("ยังไม่ได้ตั้งค่า LIFF");
  await liff.init({ liffId: config.liffId });
  const expiresAt = (liff.getDecodedIDToken()?.exp ?? 0) * 1000;
  if (!liff.isLoggedIn() || expiresAt <= Date.now()) {
    // LIFF keeps an expired ID token across visits in an external browser; only a fresh login replaces it.
    if (liff.isLoggedIn()) liff.logout();
    liff.login();
    return false;
  }
  const idToken = liff.getIDToken();
  if (!idToken) throw new Error("ไม่พบ ID token จาก LINE");
  await api<void>("POST", "/auth/line", { idToken } satisfies LineLoginInput);
  return true;
}

/** Dev login (server DEV_LOGIN on): any LINE userId, no LIFF. */
export async function devSignIn(lineUserId: string): Promise<boolean> {
  await api<void>("POST", "/auth/dev", {
    lineUserId,
  } satisfies DevLoginInput);
  return true;
}
