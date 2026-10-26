import { useEffect } from "react";

import { type Member, RoleAdmin } from "@/api/types.gen";

export type Mode = "member" | "admin";

const KEY = "mode";

/** The last mode this device was in. Storage can be blocked in a webview, so failure reads as member. */
export function storedMode(): Mode {
  try {
    return localStorage.getItem(KEY) === "admin" ? "admin" : "member";
  } catch {
    return "member";
  }
}

export function recordMode(mode: Mode) {
  try {
    localStorage.setItem(KEY, mode);
  } catch {
    // not remembered; the next launch opens in member mode
  }
}

/** Only an Admin signed in with LINE has Orders of their own. The Super Admin lives in the admin shell. */
export function canSwitchMode(member: Member): boolean {
  return member.role === RoleAdmin && !member.superAdmin;
}

let launching = true;

/**
 * The mode to open in. Only the first call per page load honors the stored mode, so the back
 * button from /admin to / does not bounce the Admin straight back.
 */
export function launchMode(member: Member): Mode {
  const atLaunch = launching;
  launching = false;
  return atLaunch && canSwitchMode(member) ? storedMode() : "member";
}

/** Records the mode when a shell is actually shown, not when the router merely preloads it. */
export function useRecordMode(mode: Mode) {
  useEffect(() => recordMode(mode), [mode]);
}
