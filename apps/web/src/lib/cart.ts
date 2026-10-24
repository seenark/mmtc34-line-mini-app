import { useRouteContext } from "@tanstack/react-router";
import { useSyncExternalStore } from "react";

import type { Option } from "@/api/types.gen";

export const MAX_QUANTITY = 99;

/** One cart row. `amount` is set only for Options without a fixed price. */
export interface CartLine {
  optionId: number;
  size: string;
  printName: string;
  number: string;
  quantity: number;
  amount: number | null;
}

export type ItemDetails = Omit<CartLine, "optionId">;

const PREFIX = "cart:";
const EMPTY: CartLine[] = [];
const listeners = new Set<() => void>();
// useSyncExternalStore needs a referentially stable snapshot, so re-parse only when the raw string changes.
const parsed = new Map<string, { raw: string | null; lines: CartLine[] }>();

function read(key: string): CartLine[] {
  const raw = localStorage.getItem(key);
  const hit = parsed.get(key);
  if (hit && hit.raw === raw) return hit.lines;
  let lines = EMPTY;
  if (raw) {
    try {
      const value: unknown = JSON.parse(raw);
      if (Array.isArray(value)) lines = value as CartLine[];
    } catch {
      // a corrupt entry reads as an empty cart
    }
  }
  parsed.set(key, { raw, lines });
  return lines;
}

function write(key: string, lines: CartLine[]) {
  if (lines.length > 0) localStorage.setItem(key, JSON.stringify(lines));
  else localStorage.removeItem(key);
  for (const notify of listeners) notify();
}

function subscribe(notify: () => void) {
  listeners.add(notify);
  const onStorage = (e: StorageEvent) => {
    if (e.key === null || e.key.startsWith(PREFIX)) notify();
  };
  window.addEventListener("storage", onStorage);
  return () => {
    listeners.delete(notify);
    window.removeEventListener("storage", onStorage);
  };
}

/** The signed-in Member's cart for one Campaign, persisted in localStorage as `cart:<memberId>:<campaignId>`. */
export function useCart(campaignId: string) {
  const memberId = useRouteContext({
    from: "/_app",
    select: (c) => c.member.id,
  });
  const key = `${PREFIX}${memberId}:${campaignId}`;
  const lines = useSyncExternalStore(subscribe, () => read(key));
  return {
    lines,
    /** An identical line (same Option and details) just gains quantity. */
    add(line: CartLine) {
      const current = read(key);
      const i = current.findIndex(
        (l) =>
          l.optionId === line.optionId &&
          l.size === line.size &&
          l.printName === line.printName &&
          l.number === line.number &&
          l.amount === line.amount,
      );
      write(
        key,
        i < 0
          ? [...current, line]
          : current.map((l, j) =>
              j === i
                ? {
                    ...l,
                    quantity: Math.min(
                      MAX_QUANTITY,
                      l.quantity + line.quantity,
                    ),
                  }
                : l,
            ),
      );
    },
    update(index: number, details: ItemDetails) {
      write(
        key,
        read(key).map((l, j) => (j === index ? { ...l, ...details } : l)),
      );
    },
    remove(index: number) {
      write(
        key,
        read(key).filter((_, j) => j !== index),
      );
    },
    clear() {
      write(key, EMPTY);
    },
  };
}

/** Logout: the next person on this device must not inherit a cart. */
export function clearAllCarts() {
  for (const key of Object.keys(localStorage)) {
    if (key.startsWith(PREFIX)) localStorage.removeItem(key);
  }
  parsed.clear();
  for (const notify of listeners) notify();
}

/** A display estimate only; the server computes the real Total at checkout. */
export function estimateTotal(options: Option[], lines: CartLine[]): number {
  const price = new Map(options.map((o) => [o.id, o.price]));
  return lines.reduce(
    (sum, l) => sum + (l.amount ?? price.get(l.optionId) ?? 0) * l.quantity,
    0,
  );
}
