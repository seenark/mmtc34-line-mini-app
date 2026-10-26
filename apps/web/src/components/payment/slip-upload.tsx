import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Loader2, Upload } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";

import { errorText } from "@/api/me";
import { useUploadSlip } from "@/api/payments";
import { type Order, OrderPendingReview } from "@/api/types.gen";

/**
 * The next step of paying, in place on the Order page: one tap opens the picker, then preview and
 * submit. The upload settles the Order (uploadSlip turns false), so the confirmation is kept by
 * this component's own mutation state rather than by the parent's condition.
 */
export function SlipUpload({ order }: { order: Order }) {
  const upload = useUploadSlip(String(order.id));
  const input = useRef<HTMLInputElement>(null);
  const [file, setFile] = useState<File>();
  const preview = useMemo(
    () => (file ? URL.createObjectURL(file) : undefined),
    [file],
  );
  useEffect(() => {
    if (preview) return () => URL.revokeObjectURL(preview);
  }, [preview]);

  if (upload.isSuccess) {
    // On the admin page the Slip is reviewed in place, so the note goes once it is reviewed.
    if (order.status !== OrderPendingReview) return null;
    return (
      <p role="status" className="bg-muted rounded-md p-3 text-sm">
        ส่งสลิปแล้ว รอแอดมินตรวจ
      </p>
    );
  }
  if (!order.actions.uploadSlip && !upload.isPending) return null;

  return (
    <Card>
      <CardContent>
        <form
          className="flex flex-col gap-3"
          onSubmit={(e) => {
            e.preventDefault();
            if (!file) return;
            upload.mutate(file, { onSuccess: () => setFile(undefined) });
          }}
        >
          <input
            ref={input}
            id="slip-file"
            type="file"
            accept="image/*"
            className="peer sr-only"
            tabIndex={file ? -1 : 0}
            disabled={upload.isPending}
            onChange={(e) => {
              const picked = e.target.files?.[0];
              e.target.value = "";
              if (!picked) return;
              upload.reset();
              setFile(picked);
            }}
          />
          {preview ? (
            <>
              <img
                src={preview}
                alt="ตัวอย่างสลิป"
                className="bg-muted max-h-96 w-full rounded-xl object-contain"
              />
              <Button type="submit" size="lg" disabled={upload.isPending}>
                {upload.isPending && (
                  <Loader2 aria-hidden className="motion-safe:animate-spin" />
                )}
                {upload.isPending ? "กำลังอัปโหลด..." : "ส่งสลิป"}
              </Button>
              <Button
                type="button"
                variant="outline"
                size="lg"
                disabled={upload.isPending}
                onClick={() => input.current?.click()}
              >
                เลือกรูปใหม่
              </Button>
            </>
          ) : (
            <label
              htmlFor="slip-file"
              className="border-secondary bg-accent text-primary peer-focus-visible:ring-ring hover:bg-accent/70 active:bg-accent/50 flex min-h-32 cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed p-5 text-center transition-colors duration-(--dur-micro) ease-(--ease-out) peer-focus-visible:ring-2"
            >
              <Upload aria-hidden className="size-8" />
              <span className="text-base font-semibold">แนบสลิปการโอน</span>
              <span className="text-muted-foreground text-sm">
                โอนแล้วกดตรงนี้เพื่อเลือกรูปสลิป
              </span>
            </label>
          )}
          {upload.error && (
            <p
              role="alert"
              className="bg-destructive/10 text-destructive rounded-md px-3 py-2 text-sm"
            >
              {errorText(upload.error)}
            </p>
          )}
        </form>
      </CardContent>
    </Card>
  );
}
