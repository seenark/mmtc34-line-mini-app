import { buttonVariants } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Checkbox } from "@mmtc34/ui/components/checkbox";
import { Label } from "@mmtc34/ui/components/label";
import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";

import { API } from "@/api/client";

export const Route = createFileRoute("/admin/campaigns/$campaignId/exports")({
  staticData: {
    nav: { slot: "adminCampaign", label: "ส่งออก", order: 40 },
  },
  component: AdminExports,
});

const EXPORTS = [
  {
    file: "factory",
    label: "รายการส่งโรงงาน",
    description:
      "รายการสินค้าแต่ละชิ้น ไซส์ ชื่อและเบอร์บนเสื้อ พร้อมยอดรวมแยกตามแบบและไซส์",
  },
  {
    file: "shipping",
    label: "รายการจัดส่ง",
    description: "คำสั่งซื้อที่เลือกจัดส่ง พร้อมชื่อผู้รับ เบอร์โทร และที่อยู่",
  },
  {
    file: "pickup",
    label: "รายการรับเอง",
    description: "คำสั่งซื้อที่เลือกรับเอง พร้อมชื่อสมาชิกและรายการสินค้า",
  },
] as const;

function AdminExports() {
  const { campaignId } = Route.useParams();
  const [includePending, setIncludePending] = useState(false);
  const query = includePending ? "?includePending=1" : "";

  return (
    <div className="flex flex-col gap-4">
      <p className="text-muted-foreground text-sm">
        ไฟล์ CSV นับเฉพาะคำสั่งซื้อที่ชำระเงินแล้ว
      </p>
      <Label className="hover:bg-accent min-h-10 w-fit max-w-full cursor-pointer gap-3 rounded-md px-1 transition-colors">
        <Checkbox
          checked={includePending}
          onCheckedChange={setIncludePending}
        />
        รวมคำสั่งซื้อที่รอตรวจสลิป
      </Label>
      <ul className="flex flex-col gap-3">
        {EXPORTS.map((e) => (
          <li key={e.file}>
            <Card>
              <CardContent className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between sm:gap-4">
                <div className="min-w-0">
                  <h2 className="text-base font-semibold">{e.label}</h2>
                  <p className="text-muted-foreground text-sm">
                    {e.description}
                  </p>
                </div>
                <a
                  download
                  href={`${API}/admin/campaigns/${campaignId}/exports/${e.file}.csv${query}`}
                  className={buttonVariants({ className: "sm:shrink-0" })}
                >
                  ดาวน์โหลด
                </a>
              </CardContent>
            </Card>
          </li>
        ))}
      </ul>
    </div>
  );
}
