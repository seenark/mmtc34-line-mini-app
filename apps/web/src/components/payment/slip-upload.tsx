import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Loader2, Upload, X } from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";

import { errorText } from "@/api/me";
import { useUploadSlip } from "@/api/payments";
import { type Order, OrderPendingReview } from "@/api/types.gen";

const MAX_SLIPS = 6;

/**
 * The next step of paying, in place on the Order page: one tap opens the picker, then previews and
 * submit (up to MAX_SLIPS images in one request). The upload settles the Order (uploadSlip turns
 * false), so the confirmation is kept by this component's own mutation state rather than by the
 * parent's condition.
 */
export function SlipUpload({ order }: { order: Order }) {
  const upload = useUploadSlip(String(order.id));
  const input = useRef<HTMLInputElement>(null);
  const [files, setFiles] = useState<File[]>([]);
  const [over, setOver] = useState(false);
  const previews = useMemo(
    () => files.map((f) => URL.createObjectURL(f)),
    [files],
  );
  useEffect(
    () => () => {
      for (const url of previews) URL.revokeObjectURL(url);
    },
    [previews],
  );

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
            if (files.length === 0) return;
            upload.mutate(files, { onSuccess: () => setFiles([]) });
          }}
        >
          <input
            ref={input}
            id="slip-file"
            type="file"
            accept="image/*"
            multiple
            className="peer sr-only"
            tabIndex={files.length > 0 ? -1 : 0}
            disabled={upload.isPending}
            onChange={(e) => {
              const picked = Array.from(e.target.files ?? []);
              e.target.value = "";
              if (picked.length === 0) return;
              upload.reset();
              const next = [...files, ...picked];
              setOver(next.length > MAX_SLIPS);
              setFiles(next.slice(0, MAX_SLIPS));
            }}
          />
          {files.length > 0 ? (
            <>
              <ul className="grid grid-cols-3 gap-2">
                {previews.map((src, i) => (
                  <li
                    key={src}
                    className="bg-muted relative aspect-square overflow-hidden rounded-md"
                  >
                    <img
                      src={src}
                      alt={`ตัวอย่างสลิป ${i + 1}`}
                      className="size-full object-cover"
                    />
                    <Button
                      type="button"
                      variant="secondary"
                      size="icon-sm"
                      aria-label={`เอารูปที่ ${i + 1} ออก`}
                      className="absolute top-1 right-1"
                      disabled={upload.isPending}
                      onClick={() => {
                        setOver(false);
                        setFiles(files.filter((_, j) => j !== i));
                      }}
                    >
                      <X aria-hidden />
                    </Button>
                  </li>
                ))}
              </ul>
              <Button type="submit" size="lg" disabled={upload.isPending}>
                {upload.isPending && (
                  <Loader2 aria-hidden className="motion-safe:animate-spin" />
                )}
                {upload.isPending ? "กำลังอัปโหลด..." : "ส่งสลิป"}
              </Button>
              {files.length < MAX_SLIPS && (
                <Button
                  type="button"
                  variant="outline"
                  size="lg"
                  disabled={upload.isPending}
                  onClick={() => input.current?.click()}
                >
                  เพิ่มรูป
                </Button>
              )}
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
          {over && (
            <p role="status" className="text-muted-foreground text-sm">
              แนบได้สูงสุด {MAX_SLIPS} รูป
            </p>
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
