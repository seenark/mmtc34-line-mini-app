import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { useMutation } from "@tanstack/react-query";
import { useId } from "react";

import { uploadImage } from "@/api/campaigns";
import { ApiError } from "@/api/client";
import { imageUrl } from "@/lib/image";

import { Field } from "./field";

/** Uploads on pick; `value` is the stored image name ("" = none). */
export function ImageField(props: {
  label: string;
  value: string;
  onChange: (name: string) => void;
}) {
  const id = useId();
  const upload = useMutation({
    mutationFn: uploadImage,
    onSuccess: props.onChange,
  });
  return (
    <div className="flex flex-col gap-2">
      {props.value && (
        <div className="flex items-end gap-3">
          <img
            src={imageUrl(props.value)}
            alt=""
            className="bg-muted size-24 rounded-md border object-contain"
          />
          <Button
            type="button"
            variant="outline"
            onClick={() => props.onChange("")}
          >
            ลบรูป
            <span className="sr-only">{props.label}</span>
          </Button>
        </div>
      )}
      <Field label={props.label}>
        <Input
          id={id}
          type="file"
          accept="image/*"
          disabled={upload.isPending}
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) upload.mutate(file);
            e.target.value = "";
          }}
        />
      </Field>
      {upload.isPending && (
        <p role="status" className="text-muted-foreground text-sm">
          กำลังอัปโหลด…
        </p>
      )}
      {upload.error && (
        <p role="alert" className="text-destructive text-sm">
          {upload.error instanceof ApiError
            ? upload.error.message
            : "ไม่สามารถอ่านรูปภาพนี้ได้"}
        </p>
      )}
    </div>
  );
}
