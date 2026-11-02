import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { ChevronLeft, ChevronRight, Loader2, X } from "lucide-react";
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

  // Resolved by name against the latest array so a move during an in-flight upload keeps the new photos.
  function move(name: string, by: -1 | 1) {
    const next = [...latest.current];
    const from = next.indexOf(name);
    const to = from + by;
    if (from < 0 || to < 0 || to >= next.length) return;
    next.splice(to, 0, ...next.splice(from, 1));
    commit(next);
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
              {i === 0 && (
                <span className="bg-primary text-primary-foreground absolute top-1 left-1 rounded-full px-2 py-0.5 text-xs font-medium">
                  ปก
                </span>
              )}
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
              <Button
                type="button"
                variant="secondary"
                size="icon-sm"
                aria-label={`เลื่อนรูปที่ ${i + 1} ไปทางซ้าย`}
                className="absolute bottom-1 left-1"
                disabled={i === 0}
                onClick={() => move(name, -1)}
              >
                <ChevronLeft aria-hidden />
              </Button>
              <Button
                type="button"
                variant="secondary"
                size="icon-sm"
                aria-label={`เลื่อนรูปที่ ${i + 1} ไปทางขวา`}
                className="absolute right-1 bottom-1"
                disabled={i === props.value.length - 1}
                onClick={() => move(name, 1)}
              >
                <ChevronRight aria-hidden />
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
