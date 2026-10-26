import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useQuery } from "@tanstack/react-query";

import { campaignQuery } from "@/api/campaigns";
import type { Order } from "@/api/types.gen";
import { BankAccount, UploadedQr } from "@/components/campaign/bank-details";

import { baht } from "./order-status";

/**
 * What to pay with: the amount owed, the Campaign's uploaded QR and the account to copy. Members
 * type the amount themselves. Bank details live on the Campaign, so a failed fetch just hides the card.
 */
export function PaymentInfo({ order }: { order: Order }) {
  const { data } = useQuery(campaignQuery(String(order.campaignId)));
  if (order.balance <= 0 || !data) return null;
  const c = data.campaign;
  if (!c.qrImage && !c.bankAccountNo) return null;
  return (
    <Card>
      <CardContent className="flex flex-col gap-4">
        <h2 className="text-base font-semibold">ช่องทางโอนเงิน</h2>
        <div className="flex flex-col items-center">
          <span className="text-muted-foreground text-sm">ยอดที่ต้องโอน</span>
          <span className="text-primary text-3xl font-extrabold tabular-nums">
            {baht(order.balance)}
          </span>
        </div>
        {c.qrImage && <UploadedQr image={c.qrImage} />}
        <BankAccount campaign={c} />
      </CardContent>
    </Card>
  );
}
