import { Card, CardContent } from "@mmtc34/ui/components/card";
import type { ReactNode } from "react";

import type { Option } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

export const priceText = (price: number | null) =>
  price === null ? "ระบุยอดเอง" : `${price.toLocaleString("th-TH")} บาท`;

function neededDetails(o: Option): string[] {
  return [
    o.sizes.length > 0 && "ขนาด",
    o.needsPrintName && "ชื่อสกรีน",
    o.needsNumber && "เบอร์",
  ].filter((d): d is string => d !== false);
}

/** `children` is the slot for per-Option actions (e.g. add to cart). */
export function OptionCard(props: { option: Option; children?: ReactNode }) {
  const o = props.option;
  const needed = neededDetails(o);
  return (
    <Card>
      {o.image && (
        <img
          src={imageUrl(o.image)}
          alt={o.name}
          className="aspect-square w-full object-cover"
        />
      )}
      <CardContent className="flex flex-col gap-2">
        <div className="flex items-start justify-between gap-3">
          <h3 className="min-w-0 text-base font-semibold break-words">
            {o.name}
          </h3>
          <span className="text-primary shrink-0 text-base font-bold tabular-nums">
            {priceText(o.price)}
          </span>
        </div>
        {o.description && (
          <p className="text-muted-foreground text-sm whitespace-pre-wrap">
            {o.description}
          </p>
        )}
        {(o.sizes.length > 0 || needed.length > 0) && (
          <dl className="text-muted-foreground grid grid-cols-[auto_minmax(0,1fr)] gap-x-2 gap-y-0.5 text-sm">
            {o.sizes.length > 0 && (
              <>
                <dt>ขนาด</dt>
                <dd className="text-foreground">{o.sizes.join(", ")}</dd>
              </>
            )}
            {needed.length > 0 && (
              <>
                <dt>ต้องระบุ</dt>
                <dd className="text-foreground">{needed.join(", ")}</dd>
              </>
            )}
          </dl>
        )}
        {props.children}
      </CardContent>
    </Card>
  );
}
