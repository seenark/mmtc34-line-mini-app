import { Button } from "@mmtc34/ui/components/button";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect, useMemo, useState } from "react";

import { errorText } from "@/api/me";
import { useUploadSlip } from "@/api/payments";

export const Route = createFileRoute("/_app/orders/$orderId/slip")({
  staticData: {
    nav: {
      slot: "orderAction",
      label: "อัปโหลดสลิป",
      order: 10,
      when: "uploadSlip",
    },
  },
  component: UploadSlip,
});

function UploadSlip() {
  const { orderId } = Route.useParams();
  const navigate = useNavigate();
  const upload = useUploadSlip(orderId);
  const [file, setFile] = useState<File>();
  const preview = useMemo(
    () => (file ? URL.createObjectURL(file) : undefined),
    [file],
  );
  useEffect(() => {
    if (preview) return () => URL.revokeObjectURL(preview);
  }, [preview]);

  return (
    <form
      className="flex flex-col gap-3"
      onSubmit={(e) => {
        e.preventDefault();
        if (!file) return;
        upload.mutate(file, {
          onSuccess: () =>
            navigate({ to: "/orders/$orderId", params: { orderId } }),
        });
      }}
    >
      <h2 className="text-sm font-medium">อัปโหลดสลิป</h2>
      <input
        type="file"
        accept="image/*"
        aria-label="เลือกรูปสลิป"
        disabled={upload.isPending}
        onChange={(e) => {
          upload.reset();
          setFile(e.target.files?.[0]);
        }}
        className="text-sm"
      />
      {preview && (
        <img
          src={preview}
          alt="ตัวอย่างสลิป"
          className="bg-muted max-h-96 w-full object-contain"
        />
      )}
      <Button type="submit" size="lg" disabled={!file || upload.isPending}>
        {upload.isPending ? "กำลังอัปโหลด..." : "ส่งสลิป"}
      </Button>
      {upload.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(upload.error)}
        </p>
      )}
    </form>
  );
}
