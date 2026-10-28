import { Button } from "@mmtc34/ui/components/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@mmtc34/ui/components/dialog";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { Minus, Plus, ZoomIn, ZoomOut } from "lucide-react";
import { useState } from "react";

import type { Option } from "@/api/types.gen";
import { Field, selectClass } from "@/components/campaign/field";
import { type ItemDetails, MAX_QUANTITY } from "@/lib/cart";
import { imageUrl } from "@/lib/image";

export const blankDetails: ItemDetails = {
  size: "",
  printName: "",
  number: "",
  quantity: 1,
  amount: null,
};

/** Mounted only while the dialog is open, so zoom resets on each open. */
function SizeChartViewer({ image }: { image: string }) {
  const [zoomed, setZoomed] = useState(false);
  return (
    <>
      <DialogHeader className="flex-row items-center justify-between gap-2 pr-10">
        <DialogTitle>ตารางไซซ์</DialogTitle>
        <Button
          type="button"
          variant="outline"
          size="sm"
          aria-pressed={zoomed}
          onClick={() => setZoomed((z) => !z)}
        >
          {zoomed ? <ZoomOut /> : <ZoomIn />}
          {zoomed ? "ย่อ" : "ขยาย"}
        </Button>
      </DialogHeader>
      <div className="bg-muted max-h-[calc(100dvh-9rem)] overflow-auto overscroll-contain rounded-md border">
        <img
          src={imageUrl(image)}
          alt="ตารางไซซ์"
          className={zoomed ? "h-auto w-[220%] max-w-none" : "h-auto w-full"}
        />
      </div>
    </>
  );
}

function SizeChartButton({ image }: { image: string }) {
  return (
    <Dialog>
      <DialogTrigger
        render={<Button type="button" variant="outline" className="shrink-0" />}
      >
        ดูตารางไซซ์
      </DialogTrigger>
      <DialogContent className="sm:max-w-2xl">
        <SizeChartViewer image={image} />
      </DialogContent>
    </Dialog>
  );
}

/**
 * The inputs an Option asks for (size, print name, number, or an amount for no-price Options)
 * plus quantity. Must sit inside a <form>: required/min/max are enforced by native validation.
 */
export function ItemFields(props: {
  id: string;
  option: Option;
  value: ItemDetails;
  onChange: (value: ItemDetails) => void;
  /** Hide quantity and amount: a placed Order's money fields are fixed. */
  detailsOnly?: boolean;
}) {
  const { id, option: o, value } = props;
  const set = (patch: Partial<ItemDetails>) =>
    props.onChange({ ...value, ...patch });
  return (
    <div className="flex flex-col gap-3">
      {o.sizes.length > 0 && (
        <Field
          label="ขนาด"
          required
          aside={o.sizeChart && <SizeChartButton image={o.sizeChart} />}
        >
          <select
            id={`${id}-size`}
            required
            className={selectClass}
            value={value.size}
            onChange={(e) => set({ size: e.target.value })}
          >
            <option value="">เลือกขนาด</option>
            {o.sizes.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
        </Field>
      )}
      {o.needsPrintName && (
        <Field label="ชื่อสกรีน" required>
          <Input
            id={`${id}-printName`}
            required
            autoComplete="off"
            value={value.printName}
            onChange={(e) => set({ printName: e.target.value })}
          />
        </Field>
      )}
      {o.needsNumber && (
        <Field label="เบอร์" required>
          <Input
            id={`${id}-number`}
            required
            autoComplete="off"
            value={value.number}
            onChange={(e) => set({ number: e.target.value })}
          />
        </Field>
      )}
      {props.detailsOnly ? null : o.price === null ? (
        <Field label="จำนวนเงิน (บาท)" required>
          <Input
            id={`${id}-amount`}
            type="number"
            inputMode="numeric"
            required
            min={1}
            step={1}
            value={value.amount ?? ""}
            onChange={(e) =>
              set({
                amount: e.target.value === "" ? null : Number(e.target.value),
              })
            }
          />
        </Field>
      ) : (
        <div className="flex items-center justify-between gap-2">
          <Label htmlFor={`${id}-quantity`} className="font-medium">
            จำนวน
          </Label>
          <div className="flex items-center gap-1">
            <Button
              type="button"
              variant="outline"
              size="icon"
              aria-label="ลดจำนวน"
              disabled={value.quantity <= 1}
              onClick={() => set({ quantity: value.quantity - 1 })}
            >
              <Minus />
            </Button>
            <Input
              id={`${id}-quantity`}
              type="number"
              inputMode="numeric"
              required
              min={1}
              max={MAX_QUANTITY}
              step={1}
              className="w-14 text-center tabular-nums"
              value={value.quantity}
              onChange={(e) => set({ quantity: Number(e.target.value) })}
            />
            <Button
              type="button"
              variant="outline"
              size="icon"
              aria-label="เพิ่มจำนวน"
              disabled={value.quantity >= MAX_QUANTITY}
              onClick={() => set({ quantity: value.quantity + 1 })}
            >
              <Plus />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
