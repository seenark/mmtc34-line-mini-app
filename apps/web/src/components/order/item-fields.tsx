import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { Minus, Plus } from "lucide-react";

import type { Option } from "@/api/types.gen";
import { Field, selectClass } from "@/components/campaign/field";
import { type ItemDetails, MAX_QUANTITY } from "@/lib/cart";

export const blankDetails: ItemDetails = {
  size: "",
  printName: "",
  number: "",
  quantity: 1,
  amount: null,
};

/**
 * The inputs an Option asks for (size, print name, number, or an amount for no-price Options)
 * plus quantity. Must sit inside a <form>: required/min/max are enforced by native validation.
 */
export function ItemFields(props: {
  id: string;
  option: Option;
  value: ItemDetails;
  onChange: (value: ItemDetails) => void;
}) {
  const { id, option: o, value } = props;
  const set = (patch: Partial<ItemDetails>) =>
    props.onChange({ ...value, ...patch });
  return (
    <div className="flex flex-col gap-3">
      {o.sizes.length > 0 && (
        <Field label="ขนาด" required>
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
      {o.price === null ? (
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
          <Label htmlFor={`${id}-quantity`}>จำนวน</Label>
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
              className="w-14 text-center"
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
