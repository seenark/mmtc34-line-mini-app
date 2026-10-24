import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useState } from "react";

import type { Campaign } from "@/api/types.gen";
import { imageUrl } from "@/lib/image";

export function BankDetails({ campaign: c }: { campaign: Campaign }) {
  const [copied, setCopied] = useState(false);
  if (!c.bankAccountNo && !c.qrImage) return null;
  return (
    <Card>
      <CardContent className="flex flex-col gap-2">
        <h2 className="text-sm font-medium">ช่องทางโอนเงิน</h2>
        {c.bankAccountNo && (
          <>
            <p>
              {c.bankName} {c.bankAccountName}
            </p>
            <div className="flex items-center justify-between gap-2">
              <span className="text-base font-semibold tabular-nums">
                {c.bankAccountNo}
              </span>
              <Button
                variant="outline"
                onClick={async () => {
                  await navigator.clipboard.writeText(c.bankAccountNo);
                  setCopied(true);
                  setTimeout(() => setCopied(false), 2000);
                }}
              >
                {copied ? "คัดลอกแล้ว" : "คัดลอกเลขบัญชี"}
              </Button>
            </div>
          </>
        )}
        {c.qrImage && (
          <img
            src={imageUrl(c.qrImage)}
            alt="QR สำหรับโอนเงิน"
            className="mx-auto w-full max-w-64"
          />
        )}
      </CardContent>
    </Card>
  );
}
