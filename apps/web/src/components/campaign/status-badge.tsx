import type { Campaign } from "@/api/types.gen";
import { CampaignDraft, CampaignOpen } from "@/api/types.gen";

const labels = { open: "เปิดรับ", closed: "ปิดแล้ว", draft: "ฉบับร่าง" } as const;

export function StatusBadge({ status }: { status: Campaign["status"] }) {
  const label = labels[status as keyof typeof labels] ?? status;
  const tone =
    status === CampaignOpen
      ? "bg-primary text-primary-foreground"
      : status === CampaignDraft
        ? "border"
        : "bg-muted text-muted-foreground";
  return <span className={`px-2 py-0.5 text-xs ${tone}`}>{label}</span>;
}

export function closeAtText(closeAt: string | null): string | null {
  return closeAt
    ? `ปิดรับ ${new Date(closeAt).toLocaleString("th-TH", { dateStyle: "medium", timeStyle: "short" })}`
    : null;
}
