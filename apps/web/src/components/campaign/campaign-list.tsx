import { Card } from "@mmtc34/ui/components/card";
import { Empty, EmptyDescription } from "@mmtc34/ui/components/empty";
import { Link } from "@tanstack/react-router";
import { Package } from "lucide-react";

import type { Campaign } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

import { closeAtText, StatusBadge } from "./status-badge";

const LINK = "group block rounded-xl";
const CARD =
  "w-full flex-row gap-0 py-0 transition-colors duration-(--dur-micro) ease-(--ease-out) group-hover:bg-accent group-active:bg-accent";
const THUMB = "bg-muted min-h-24 w-24 shrink-0";

/** `admin` links to the editor instead of the Member view and shows drafts' badge. */
export function CampaignList(props: {
  campaigns: Campaign[];
  admin?: boolean;
}) {
  if (props.campaigns.length === 0) {
    return (
      <Empty className="flex-none rounded-xl border">
        <EmptyDescription>ยังไม่มีอีเวนต์</EmptyDescription>
      </Empty>
    );
  }
  return (
    <ul className="flex flex-col gap-3">
      {props.campaigns.map((c) => {
        // The jersey cover is a square product photo, so the thumbnail is square and flush left.
        const card = (
          <Card className={CARD}>
            {c.photos[0] ? (
              <img
                src={imageUrl(c.photos[0])}
                alt=""
                loading="lazy"
                className={`${THUMB} object-cover`}
              />
            ) : (
              <div
                aria-hidden
                className={`${THUMB} text-muted-foreground flex items-center justify-center`}
              >
                <Package className="size-6" />
              </div>
            )}
            <div className="flex min-w-0 flex-1 flex-col justify-center gap-1.5 p-3">
              <span className="text-primary text-base font-semibold [overflow-wrap:anywhere]">
                {c.title}
              </span>
              <div className="flex flex-wrap items-center gap-x-2 gap-y-1">
                <StatusBadge status={c.status} />
                <span className="text-muted-foreground text-sm">
                  {closeAtText(c.closeAt) ?? "ไม่กำหนดวันปิดรับ"}
                </span>
              </div>
            </div>
          </Card>
        );
        const campaignId = String(c.id);
        return (
          <li key={c.id}>
            {props.admin ? (
              <Link
                to="/admin/campaigns/$campaignId"
                params={{ campaignId }}
                className={LINK}
              >
                {card}
              </Link>
            ) : (
              <Link
                to="/campaigns/$campaignId"
                params={{ campaignId }}
                className={LINK}
              >
                {card}
              </Link>
            )}
          </li>
        );
      })}
    </ul>
  );
}
