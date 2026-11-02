import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Loader2, X } from "lucide-react";
import { useEffect, useId, useRef, useState } from "react";

import { uploadImage } from "@/api/campaigns";
import { ApiError } from "@/api/client";
import { imageUrl } from "@/lib/image";

import { Field } from "./field";

const tile =
  "bg-muted relative flex aspect-square items-center justify-center overflow-hidden rounded-md";

/**
 * `value` is the ordered stored image names. Picked files upload one at a time and each name is
 * appended as it lands; files still uploading are kept out of `value` and shown as placeholders.
 */
export function PhotosField(props: {
  label: string;
  value: string[];
  onChange: (next: string[]) => void;
  onUploadingChange: (uploading: boolean) => void;
}) {
  const id = useId();
  const [pending, setPending] = useState(0);
  const [error, setError] = useState<string | null>(null);
  // The upload loop outlives the render that started it, so it appends to the latest array.
  const latest = useRef(props.value);
  latest.current = props.value;

  const { onUploadingChange } = props;
  useEffect(() => {
    onUploadingChange(pending > 0);
  }, [pending, onUploadingChange]);

  function commit(next: string[]) {
    latest.current = next;
    props.onChange(next);
  }

  async function upload(files: File[]) {
    setError(null);
    setPending(files.length);
    for (const file of files) {
      try {
        commit([...latest.current, await uploadImage(file)]);
      } catch (e) {
        setError(e instanceof ApiError ? e.message : "ไม่สามารถอ่านรูปภาพนี้ได้");
      }
      setPending((n) => n - 1);
    }
  }

  return (
    <div className="flex flex-col gap-2">
      {(props.value.length > 0 || pending > 0) && (
        <ul className="grid grid-cols-3 gap-2">
          {props.value.map((name, i) => (
            <li key={name} className={tile}>
              <img
                src={imageUrl(name)}
                alt=""
                className="size-full object-cover"
              />
              <Button
                type="button"
                variant="secondary"
                size="icon-sm"
                aria-label={`ลบรูปที่ ${i + 1}`}
                className="absolute top-1 right-1"
                onClick={() => commit(latest.current.filter((n) => n !== name))}
              >
                <X aria-hidden />
              </Button>
            </li>
          ))}
          {Array.from({ length: pending }, (_, i) => (
            // biome-ignore lint/suspicious/noArrayIndexKey: placeholders have no identity
            <li key={i} className={tile}>
              <Loader2
                aria-hidden
                className="text-muted-foreground motion-safe:animate-spin"
              />
            </li>
          ))}
        </ul>
      )}
      <Field label={props.label}>
        <Input
          id={id}
          type="file"
          accept="image/*"
          multiple
          disabled={pending > 0}
          onChange={(e) => {
            const files = Array.from(e.target.files ?? []);
            e.target.value = "";
            if (files.length > 0) void upload(files);
          }}
        />
      </Field>
      {pending > 0 && (
        <p role="status" className="text-muted-foreground text-sm">
          กำลังอัปโหลด…
        </p>
      )}
      {error && (
        <p role="alert" className="text-destructive text-sm">
          {error}
        </p>
      )}
    </div>
  );
}
