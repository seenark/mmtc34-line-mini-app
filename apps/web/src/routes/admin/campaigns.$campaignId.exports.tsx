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
      <Label>
        <Checkbox
          checked={includePending}
          onCheckedChange={setIncludePending}
        />
        รวมคำสั่งซื้อที่รอตรวจสลิป
      </Label>
      {EXPORTS.map((e) => (
        <Card key={e.file}>
          <CardContent className="flex flex-col gap-3">
            <div>
              <p className="font-medium">{e.label}</p>
              <p className="text-muted-foreground text-sm">{e.description}</p>
            </div>
            <a
              download
              href={`${API}/admin/campaigns/${campaignId}/exports/${e.file}.csv${query}`}
              className={buttonVariants()}
            >
              ดาวน์โหลด
            </a>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
