import { Card, CardContent } from "@mmtc34/ui/components/card";

import type { SponsorPosition } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

export function SponsorCard(props: { position: SponsorPosition }) {
  const p = props.position;
  return (
    <Card>
      <CardContent className="flex items-center gap-3">
        {p.logo && (
          <img
            src={imageUrl(p.logo)}
            alt={p.sponsorName}
            loading="lazy"
            className="bg-muted size-16 shrink-0 rounded-md object-contain p-1"
          />
        )}
        <div className="flex min-w-0 flex-col gap-0.5">
          <span className="text-muted-foreground text-sm [overflow-wrap:anywhere]">
            {p.name}
          </span>
          {p.sponsorName ? (
            <span className="text-primary text-base font-semibold [overflow-wrap:anywhere]">
              {p.sponsorName}
            </span>
          ) : (
            <span className="text-muted-foreground text-base">ว่าง</span>
          )}
          {p.memberNickname && (
            <span className="text-muted-foreground text-sm [overflow-wrap:anywhere]">
              {p.memberNickname}
            </span>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
