import liff from "@line/liff";
import { queryOptions } from "@tanstack/react-query";

import { api } from "@/api/client";
import type { Config, DevLoginInput, LineLoginInput } from "@/api/types.gen";

export const configQuery = queryOptions({
  queryKey: ["config"],
  queryFn: () => api<Config>("GET", "/config"),
});

let liffReady: Promise<void> | undefined;

/** The URL the page was opened with, captured before liff.init rewrites it; LINE login returns here. */
const launchUrl = window.location.href;

/**
 * Runs liff.init once per page load. It must run before the router touches the URL: LINE hands
 * the login result to the endpoint URL as a fragment or `code` query, and a router redirect
 * (e.g. `/` → `/login`) would drop it, so LIFF would ask LINE to log in again forever.
 */
export function initLiff(liffId: string): Promise<void> {
  liffReady ??= liff.init({ liffId }).catch((err: unknown) => {
    liffReady = undefined;
    throw err;
  });
  return liffReady;
}

const ID_TOKEN_REFRESH_WINDOW_MS = 30_000;

/** A LINE sign-in problem with a Thai message the user can act on (or report). */
export class LineSignInError extends Error {}

/**
 * LIFF login. Resolves true once the session cookie is set. Resolves false when LIFF is
 * redirecting the page to LINE, so the caller must not navigate.
 *
 * Inside the LINE app this never calls liff.login(): LINE logs the user in during liff.init, and a
 * login() there only reloads the page into the same state, which is an endless "Logging in" loop.
 * Problems inside LINE surface as an error the user can read instead.
 */
export async function signIn(config: Config): Promise<boolean> {
  if (!config.liffId) throw new LineSignInError("ยังไม่ได้ตั้งค่า LIFF");
  await initLiff(config.liffId).catch((err: unknown) => {
    const code =
      err instanceof Error && "code" in err ? String(err.code) : String(err);
    throw new LineSignInError(`เริ่มต้น LINE ไม่สำเร็จ (${code})`);
  });
  const inClient = liff.isInClient();
  if (!liff.isLoggedIn()) {
    if (inClient)
      throw new LineSignInError("LINE ยังไม่ได้เข้าสู่ระบบ กรุณาปิดแล้วเปิดแอปใหม่");
    liff.login({ redirectUri: launchUrl });
    return false;
  }
  const exp = liff.getDecodedIDToken()?.exp;
  if (
    typeof exp === "number" &&
    exp * 1000 <= Date.now() + ID_TOKEN_REFRESH_WINDOW_MS
  ) {
    if (inClient)
      throw new LineSignInError("การเข้าสู่ระบบ LINE หมดอายุ กรุณาปิดแล้วเปิดแอปใหม่");
    // An external browser keeps an expired ID token across visits; only a fresh login replaces it.
    liff.logout();
    liff.login({ redirectUri: launchUrl });
    return false;
  }
  const idToken = liff.getIDToken();
  if (!idToken) {
    throw new LineSignInError(
      "LINE ไม่ส่ง ID token มา (ตรวจว่าเปิด scope openid ใน LINE Developers)",
    );
  }
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
