import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useEffect, useRef, useState } from "react";

import type { Campaign } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

export function UploadedQr({ image }: { image: string }) {
  return (
    <img
      src={imageUrl(image)}
      alt="QR สำหรับโอนเงิน"
      className="bg-background mx-auto aspect-square w-full max-w-64 rounded-xl border object-contain p-2"
    />
  );
}

type CopyState = "idle" | "copied" | "failed";

const copyLabel: Record<CopyState, string> = {
  idle: "คัดลอกเลขบัญชี",
  copied: "คัดลอกแล้ว",
  failed: "คัดลอกไม่ได้",
};

function CopyAccountNo({ digits }: { digits: string }) {
  const [state, setState] = useState<CopyState>("idle");
  const timer = useRef<number>(undefined);
  useEffect(() => () => window.clearTimeout(timer.current), []);
  return (
    <Button
      type="button"
      variant="outline"
      className="w-full"
      aria-live="polite"
      onClick={async () => {
        const ok = await navigator.clipboard.writeText(digits).then(
          () => true,
          () => false,
        );
        setState(ok ? "copied" : "failed");
        window.clearTimeout(timer.current);
        timer.current = window.setTimeout(() => setState("idle"), 2000);
      }}
    >
      {copyLabel[state]}
    </Button>
  );
}

/** Bank name, account name and number with a copy button; null when the Campaign has no account. */
export function BankAccount({ campaign: c }: { campaign: Campaign }) {
  if (!c.bankAccountNo) return null;
  return (
    <div className="flex flex-col gap-3">
      <dl className="flex flex-col gap-0.5 text-sm">
        {c.bankName && (
          <div className="flex gap-2">
            <dt className="text-muted-foreground shrink-0">ธนาคาร</dt>
            <dd className="min-w-0 font-medium break-words">{c.bankName}</dd>
          </div>
        )}
        {c.bankAccountName && (
          <div className="flex gap-2">
            <dt className="text-muted-foreground shrink-0">ชื่อบัญชี</dt>
            <dd className="min-w-0 font-medium break-words">
              {c.bankAccountName}
            </dd>
          </div>
        )}
      </dl>
      <p className="text-xl font-bold tracking-wide break-all tabular-nums">
        {c.bankAccountNo}
      </p>
      <CopyAccountNo digits={c.bankAccountNo.replace(/\D/g, "")} />
    </div>
  );
}

/** Campaign-level transfer info for browsing; the per-Order payment card is PaymentInfo. */
export function BankDetails({ campaign: c }: { campaign: Campaign }) {
  if (!c.bankAccountNo && !c.qrImage) return null;
  return (
    <Card>
      <CardContent className="flex flex-col gap-4">
        <h2 className="text-base font-semibold">ช่องทางโอนเงิน</h2>
        {c.qrImage && <UploadedQr image={c.qrImage} />}
        <BankAccount campaign={c} />
      </CardContent>
    </Card>
  );
}
