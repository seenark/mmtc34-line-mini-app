import { buttonVariants } from "@mmtc34/ui/components/button";
import { Link } from "@tanstack/react-router";

import type { Option } from "@/api/types.gen";
import { priceText } from "@/components/campaign/option-card";
import { type CartLine, estimateTotal } from "@/lib/cart";

/** Sits above the member layout's sticky bottom nav (about 3.6rem tall). */
export function CartBar(props: {
  campaignId: string;
  options: Option[];
  lines: CartLine[];
}) {
  if (props.lines.length === 0) return null;
  const count = props.lines.reduce((n, l) => n + l.quantity, 0);
  return (
    <div className="bg-background sticky bottom-16 flex items-center justify-between gap-2 border p-3">
      <div className="flex flex-col">
        <span className="text-sm font-medium">ตะกร้า {count} ชิ้น</span>
        <span className="text-muted-foreground text-xs">
          ประมาณ {priceText(estimateTotal(props.options, props.lines))}
        </span>
      </div>
      <Link
        to="/campaigns/$campaignId/checkout"
        params={{ campaignId: props.campaignId }}
        className={buttonVariants()}
      >
        ไปชำระเงิน
      </Link>
    </div>
  );
}
