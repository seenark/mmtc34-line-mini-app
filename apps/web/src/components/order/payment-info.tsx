import { useQuery } from "@tanstack/react-query";

import { campaignQuery } from "@/api/campaigns";
import { BankDetails } from "@/components/campaign/bank-details";

/** Bank account and QR live on the Campaign, not the Order. A failed fetch just hides the card. */
export function PaymentInfo({ campaignId }: { campaignId: number }) {
  const { data } = useQuery(campaignQuery(String(campaignId)));
  return data ? <BankDetails campaign={data.campaign} /> : null;
}
