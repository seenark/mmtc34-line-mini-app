import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Checkbox } from "@mmtc34/ui/components/checkbox";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { useId, useState } from "react";

import { useDeleteOption, useSaveOption } from "@/api/campaigns";
import { errorText } from "@/api/me";
import type { Option } from "@/api/types.gen";

import { Field } from "./field";
import { ImageField } from "./image-field";

const JERSEY_SIZES = "K2XS,K3XS,KXS,2XS,S,M,L,XL,2XL,3XL,4XL,6XL";

export function OptionsEditor(props: {
  campaignId: string;
  options: Option[];
}) {
  const [adding, setAdding] = useState(false);
  return (
    <section className="flex flex-col gap-3">
      <h2 className="text-base font-semibold">ตัวเลือก</h2>
      {props.options.map((o) => (
        <OptionForm key={o.id} campaignId={props.campaignId} option={o} />
      ))}
      {adding ? (
        <OptionForm
          campaignId={props.campaignId}
          onSaved={() => setAdding(false)}
        />
      ) : (
        <Button variant="outline" size="lg" onClick={() => setAdding(true)}>
          เพิ่มตัวเลือก
        </Button>
      )}
    </section>
  );
}

/** `option` unset = a new Option; saving calls `onSaved`. */
function OptionForm(props: {
  campaignId: string;
  option?: Option;
  onSaved?: () => void;
}) {
  const o = props.option;
  const uid = useId();
  const save = useSaveOption(props.campaignId);
  const remove = useDeleteOption();
  const [image, setImage] = useState(o?.image ?? "");
  const [sizeChart, setSizeChart] = useState(o?.sizeChart ?? "");
  const [memberChosen, setMemberChosen] = useState(
    o ? o.price === null : false,
  );
  const [sizes, setSizes] = useState(o?.sizes.join(",") ?? "");
  const [needsPrintName, setNeedsPrintName] = useState(
    o?.needsPrintName ?? false,
  );
  const [needsNumber, setNeedsNumber] = useState(o?.needsNumber ?? false);
  const error = save.error ?? remove.error;
  const sizeList = sizes
    .split(",")
    .map((s) => s.trim())
    .filter(Boolean);
  const hasSizes = sizeList.length > 0;

  return (
    <Card>
      <CardContent>
        <form
          className="flex flex-col gap-4"
          onSubmit={(e) => {
            e.preventDefault();
            const f = new FormData(e.currentTarget);
            const text = (name: string) => String(f.get(name) ?? "").trim();
            save.mutate(
              {
                id: o?.id,
                input: {
                  name: text("name"),
                  description: text("description"),
                  image,
                  sizeChart: hasSizes ? sizeChart : "",
                  price: memberChosen
                    ? null
                    : Math.trunc(Number(text("price"))),
                  sizes: sizeList,
                  needsPrintName,
                  needsNumber,
                  sortOrder: Math.trunc(Number(text("sortOrder"))),
                },
              },
              { onSuccess: props.onSaved },
            );
          }}
        >
          <Field label="ชื่อตัวเลือก" required>
            <Input
              id={`${uid}-name`}
              name="name"
              required
              defaultValue={o?.name}
            />
          </Field>
          <Field label="รายละเอียด">
            <Textarea
              id={`${uid}-description`}
              name="description"
              defaultValue={o?.description}
            />
          </Field>
          <ImageField label="รูป" value={image} onChange={setImage} />
          <Label className="min-h-9 cursor-pointer">
            <Checkbox
              checked={memberChosen}
              onCheckedChange={setMemberChosen}
            />
            ให้สมาชิกระบุยอดเอง
          </Label>
          {!memberChosen && (
            <Field label="ราคา (บาท)" required>
              <Input
                id={`${uid}-price`}
                name="price"
                type="number"
                inputMode="numeric"
                min={0}
                step={1}
                required
                defaultValue={o?.price ?? ""}
              />
            </Field>
          )}
          <Field label="ขนาด (คั่นด้วยจุลภาค)">
            <Input
              id={`${uid}-sizes`}
              value={sizes}
              onChange={(e) => setSizes(e.target.value)}
            />
          </Field>
          <Button
            type="button"
            variant="outline"
            className="self-start"
            onClick={() => setSizes(JERSEY_SIZES)}
          >
            ใช้ไซซ์เสื้อมาตรฐาน
          </Button>
          {hasSizes && (
            <ImageField
              label="รูปตารางไซซ์"
              value={sizeChart}
              onChange={setSizeChart}
            />
          )}
          <Label className="min-h-9 cursor-pointer">
            <Checkbox
              checked={needsPrintName}
              onCheckedChange={setNeedsPrintName}
            />
            ต้องระบุชื่อสกรีน
          </Label>
          <Label className="min-h-9 cursor-pointer">
            <Checkbox checked={needsNumber} onCheckedChange={setNeedsNumber} />
            ต้องระบุเบอร์
          </Label>
          <Field label="ลำดับการแสดง">
            <Input
              id={`${uid}-sortOrder`}
              name="sortOrder"
              type="number"
              inputMode="numeric"
              step={1}
              defaultValue={o?.sortOrder ?? 0}
            />
          </Field>
          {error && (
            <p role="alert" className="text-destructive text-sm">
              {errorText(error)}
            </p>
          )}
          {save.isSuccess && o && (
            <p role="status" className="text-muted-foreground text-sm">
              บันทึกแล้ว
            </p>
          )}
          <div className="flex gap-2">
            <Button
              type="submit"
              className="min-w-0 flex-1"
              size="lg"
              disabled={save.isPending}
            >
              {o ? "บันทึกตัวเลือก" : "เพิ่มตัวเลือก"}
            </Button>
            {o && (
              <Button
                type="button"
                variant="destructive"
                size="lg"
                disabled={remove.isPending}
                onClick={() => {
                  if (confirm(`ลบตัวเลือก ${o.name}?`)) remove.mutate(o.id);
                }}
              >
                ลบ
              </Button>
            )}
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
