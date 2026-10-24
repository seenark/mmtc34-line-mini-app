import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Link } from "@tanstack/react-router";

import type { Campaign } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

import { closeAtText, StatusBadge } from "./status-badge";

/** `admin` links to the editor instead of the Member view and shows drafts' badge. */
export function CampaignList(props: {
  campaigns: Campaign[];
  admin?: boolean;
}) {
  if (props.campaigns.length === 0) {
    return <p className="text-muted-foreground text-sm">ยังไม่มีแคมเปญ</p>;
  }
  return (
    <ul className="flex flex-col gap-3">
      {props.campaigns.map((c) => {
        const card = (
          <Card>
            {c.coverImage && (
              <img
                src={imageUrl(c.coverImage)}
                alt=""
                className="aspect-video w-full object-cover"
              />
            )}
            <CardContent className="flex flex-col gap-1">
              <div className="flex items-start justify-between gap-2">
                <span className="text-sm font-medium">{c.title}</span>
                <StatusBadge status={c.status} />
              </div>
              <span className="text-muted-foreground">
                {closeAtText(c.closeAt) ?? "ไม่กำหนดวันปิดรับ"}
              </span>
            </CardContent>
          </Card>
        );
        const campaignId = String(c.id);
        return (
          <li key={c.id}>
            {props.admin ? (
              <Link to="/admin/campaigns/$campaignId" params={{ campaignId }}>
                {card}
              </Link>
            ) : (
              <Link to="/campaigns/$campaignId" params={{ campaignId }}>
                {card}
              </Link>
            )}
          </li>
        );
      })}
    </ul>
  );
}
