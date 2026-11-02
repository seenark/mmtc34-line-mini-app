import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

import {
  type Campaign,
  CampaignClosed,
  CampaignDonation,
  CampaignDraft,
  CampaignEvent,
  type CampaignInput,
  CampaignMerch,
  CampaignOpen,
} from "@/api/types.gen";
import { configQuery } from "@/lib/liff";

import { Field, selectClass } from "./field";
import { ImageField } from "./image-field";
import { PhotosField } from "./photos-field";

const types = [
  { value: CampaignMerch, label: "สินค้า" },
  { value: CampaignDonation, label: "บริจาค" },
  { value: CampaignEvent, label: "กิจกรรม" },
];

const statuses = [
  { value: CampaignDraft, label: "ฉบับร่าง (สมาชิกไม่เห็น)" },
  { value: CampaignOpen, label: "เปิดรับ" },
  { value: CampaignClosed, label: "ปิดรับ" },
];

/** `datetime-local` wants local wall time without a zone. */
function toLocalInput(iso: string | null): string {
  if (!iso) return "";
  const d = new Date(iso);
  return new Date(d.getTime() - d.getTimezoneOffset() * 60_000)
    .toISOString()
    .slice(0, 16);
}

/** Shared by create and edit; `initial` unset = blank draft. */
export function CampaignForm(props: {
  initial?: Campaign;
  submitLabel: string;
  pending: boolean;
  error?: string;
  onSubmit: (input: CampaignInput) => void;
}) {
  const c = props.initial;
  const { data: config } = useQuery(configQuery);
  const [photos, setPhotos] = useState(c?.photos ?? []);
  const [uploading, setUploading] = useState(false);
  const [qrImage, setQrImage] = useState(c?.qrImage ?? "");
  const [bankCode, setBankCode] = useState(c?.bankCode ?? "");

  const banks = config?.banks ?? [];
  // Keeps the saved bank selected while the registry loads or if it no longer lists the code.
  const bankOptions =
    bankCode && !banks.some((b) => b.code === bankCode)
      ? [...banks, { code: bankCode, name: c?.bankName || bankCode }]
      : banks;

  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        const f = new FormData(e.currentTarget);
        const text = (name: string) => String(f.get(name) ?? "").trim();
        const closeAt = text("closeAt");
        props.onSubmit({
          type: text("type"),
          title: text("title"),
          description: text("description"),
          photos,
          bankCode,
          bankAccountNo: text("bankAccountNo"),
          bankAccountName: text("bankAccountName"),
          qrImage,
          status: text("status"),
          closeAt: closeAt ? new Date(closeAt).toISOString() : null,
        });
      }}
    >
      <Field label="ประเภท" required>
        <select
          id="type"
          name="type"
          className={selectClass}
          defaultValue={c?.type ?? CampaignMerch}
        >
          {types.map((t) => (
            <option key={t.value} value={t.value}>
              {t.label}
            </option>
          ))}
        </select>
      </Field>
      <Field label="ชื่ออีเวนต์" required>
        <Input id="title" name="title" required defaultValue={c?.title} />
      </Field>
      <Field label="รายละเอียด">
        <Textarea
          id="description"
          name="description"
          defaultValue={c?.description}
        />
      </Field>
      <PhotosField
        label="รูปภาพ"
        value={photos}
        onChange={setPhotos}
        onUploadingChange={setUploading}
      />
      <fieldset className="flex min-w-0 flex-col gap-4 rounded-xl border p-4">
        <legend className="px-1 text-base font-semibold">ช่องทางรับเงิน</legend>
        <Field label="ธนาคาร">
          <select
            id="bankCode"
            className={selectClass}
            value={bankCode}
            onChange={(e) => setBankCode(e.target.value)}
          >
            <option value="">ไม่ระบุธนาคาร</option>
            {bankOptions.map((b) => (
              <option key={b.code} value={b.code}>
                {b.name}
              </option>
            ))}
          </select>
        </Field>
        <Field label="เลขที่บัญชี">
          <Input
            id="bankAccountNo"
            name="bankAccountNo"
            inputMode="numeric"
            defaultValue={c?.bankAccountNo}
          />
        </Field>
        <Field label="ชื่อบัญชี">
          <Input
            id="bankAccountName"
            name="bankAccountName"
            defaultValue={c?.bankAccountName}
          />
        </Field>
        <ImageField
          label="รูป QR รับเงิน (เช่น จาก K PLUS)"
          value={qrImage}
          onChange={setQrImage}
        />
      </fieldset>
      <Field label="สถานะ" required>
        <select
          id="status"
          name="status"
          className={selectClass}
          defaultValue={c?.status ?? CampaignDraft}
        >
          {statuses.map((s) => (
            <option key={s.value} value={s.value}>
              {s.label}
            </option>
          ))}
        </select>
      </Field>
      <Field label="ปิดรับอัตโนมัติเมื่อ (เว้นว่างได้)">
        <Input
          id="closeAt"
          name="closeAt"
          type="datetime-local"
          defaultValue={toLocalInput(c?.closeAt ?? null)}
        />
      </Field>
      {props.error && (
        <p role="alert" className="text-destructive text-sm">
          {props.error}
        </p>
      )}
      <Button type="submit" size="lg" disabled={props.pending || uploading}>
        {props.submitLabel}
      </Button>
    </form>
  );
}
