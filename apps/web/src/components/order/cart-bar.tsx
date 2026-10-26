import { buttonVariants } from "@mmtc34/ui/components/button";
import { Link } from "@tanstack/react-router";
import { ShoppingBag } from "lucide-react";

import type { Option } from "@/api/types.gen";
import { priceText } from "@/components/campaign/option-card";
import { type CartLine, estimateTotal } from "@/lib/cart";

/**
 * Sits above the member layout's sticky bottom nav (about 3.6rem tall plus the safe-area inset).
 * The trailing spacer is the bar's own height, so the page can scroll the last action clear of it.
 */
export function CartBar(props: {
  campaignId: string;
  options: Option[];
  lines: CartLine[];
}) {
  if (props.lines.length === 0) return null;
  const count = props.lines.reduce((n, l) => n + l.quantity, 0);
  return (
    <>
      <div className="bg-card ring-foreground/10 sticky bottom-[calc(3.6rem+env(safe-area-inset-bottom))] z-10 flex items-center justify-between gap-3 rounded-xl p-3 ring-1">
        <div className="flex min-w-0 items-center gap-3">
          <span
            aria-hidden
            className="bg-accent text-accent-foreground flex size-10 shrink-0 items-center justify-center rounded-md"
          >
            <ShoppingBag className="size-5" />
          </span>
          <div className="flex min-w-0 flex-col">
            <span className="text-primary text-sm font-semibold">
              ตะกร้า {count} ชิ้น
            </span>
            <span className="text-muted-foreground text-sm tabular-nums">
              ประมาณ {priceText(estimateTotal(props.options, props.lines))}
            </span>
          </div>
        </div>
        <Link
          to="/campaigns/$campaignId/checkout"
          params={{ campaignId: props.campaignId }}
          className={buttonVariants()}
        >
          ไปชำระเงิน
        </Link>
      </div>
      <div aria-hidden className="h-16" />
    </>
  );
}
