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
      <CardContent className="flex flex-col gap-1.5">
        <div className="flex items-start justify-between gap-2">
          <span className="text-sm font-medium">{o.name}</span>
          <span className="shrink-0 text-sm font-semibold">
            {priceText(o.price)}
          </span>
        </div>
        {o.description && (
          <p className="text-muted-foreground whitespace-pre-wrap">
            {o.description}
          </p>
        )}
        {o.sizes.length > 0 && <p>ขนาด: {o.sizes.join(", ")}</p>}
        {needed.length > 0 && <p>ต้องระบุ: {needed.join(", ")}</p>}
        {props.children}
      </CardContent>
    </Card>
  );
}
