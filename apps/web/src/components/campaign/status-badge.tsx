import { cn } from "@mmtc34/ui/lib/utils";
import type { ReactNode } from "react";

import type { Campaign } from "@/api/types.gen";
import { CampaignDraft, CampaignOpen } from "@/api/types.gen";

const labels = { open: "เปิดรับ", closed: "ปิดแล้ว", draft: "ฉบับร่าง" } as const;

/** Pill shared by every status chip (Campaign, Order, Payment). Tone classes come from the caller. */
export function ToneBadge(props: {
  tone: string;
  /** Colour class for the leading dot; omit for a text-only pill. */
  dot?: string;
  children: ReactNode;
}) {
  return (
    <span
      className={cn(
        "inline-flex shrink-0 items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-semibold whitespace-nowrap",
        props.tone,
      )}
    >
      {props.dot && (
        <span aria-hidden className={cn("size-1.5 rounded-full", props.dot)} />
      )}
      {props.children}
    </span>
  );
}

export function StatusBadge({ status }: { status: Campaign["status"] }) {
  const label = labels[status as keyof typeof labels] ?? status;
  if (status === CampaignOpen) {
    return (
      <ToneBadge tone="bg-accent text-accent-foreground" dot="bg-secondary">
        {label}
      </ToneBadge>
    );
  }
  if (status === CampaignDraft) {
    return (
      <ToneBadge tone="text-muted-foreground ring-input ring-1 ring-inset">
        {label}
      </ToneBadge>
    );
  }
  return <ToneBadge tone="bg-muted text-muted-foreground">{label}</ToneBadge>;
}

export function closeAtText(closeAt: string | null): string | null {
  return closeAt
    ? `ปิดรับ ${new Date(closeAt).toLocaleString("th-TH", { dateStyle: "medium", timeStyle: "short" })}`
    : null;
}
