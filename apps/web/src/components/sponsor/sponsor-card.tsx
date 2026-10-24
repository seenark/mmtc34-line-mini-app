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
            className="size-16 shrink-0 object-contain"
          />
        )}
        <div className="flex min-w-0 flex-col gap-0.5">
          <span className="text-sm font-medium">{p.name}</span>
          <span className="text-sm">{p.sponsorName || "ว่าง"}</span>
          {p.memberNickname && (
            <span className="text-muted-foreground">{p.memberNickname}</span>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
