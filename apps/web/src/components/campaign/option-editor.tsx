import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Checkbox } from "@mmtc34/ui/components/checkbox";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { useState } from "react";

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
      <h2 className="text-lg font-semibold">ตัวเลือก</h2>
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
  const save = useSaveOption(props.campaignId);
  const remove = useDeleteOption();
  const [image, setImage] = useState(o?.image ?? "");
  const [memberChosen, setMemberChosen] = useState(
    o ? o.price === null : false,
  );
  const [sizes, setSizes] = useState(o?.sizes.join(",") ?? "");
  const [needsPrintName, setNeedsPrintName] = useState(
    o?.needsPrintName ?? false,
  );
  const [needsNumber, setNeedsNumber] = useState(o?.needsNumber ?? false);
  const error = save.error ?? remove.error;

  return (
    <Card>
      <CardContent>
        <form
          className="flex flex-col gap-3"
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
                  price: memberChosen
                    ? null
                    : Math.trunc(Number(text("price"))),
                  sizes: sizes
                    .split(",")
                    .map((s) => s.trim())
                    .filter(Boolean),
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
              id={`name-${o?.id}`}
              name="name"
              required
              defaultValue={o?.name}
            />
          </Field>
          <Field label="รายละเอียด">
            <Textarea
              id={`description-${o?.id}`}
              name="description"
              defaultValue={o?.description}
            />
          </Field>
          <ImageField label="รูป" value={image} onChange={setImage} />
          <Label>
            <Checkbox
              checked={memberChosen}
              onCheckedChange={setMemberChosen}
            />
            ให้สมาชิกระบุยอดเอง
          </Label>
          {!memberChosen && (
            <Field label="ราคา (บาท)" required>
              <Input
                id={`price-${o?.id}`}
                name="price"
                type="number"
                min={0}
                step={1}
                required
                defaultValue={o?.price ?? ""}
              />
            </Field>
          )}
          <Field label="ขนาด (คั่นด้วยจุลภาค)">
            <Input
              id={`sizes-${o?.id}`}
              value={sizes}
              onChange={(e) => setSizes(e.target.value)}
            />
          </Field>
          <Button
            type="button"
            variant="outline"
            onClick={() => setSizes(JERSEY_SIZES)}
          >
            ใช้ไซซ์เสื้อมาตรฐาน
          </Button>
          <Label>
            <Checkbox
              checked={needsPrintName}
              onCheckedChange={setNeedsPrintName}
            />
            ต้องระบุชื่อสกรีน
          </Label>
          <Label>
            <Checkbox checked={needsNumber} onCheckedChange={setNeedsNumber} />
            ต้องระบุเบอร์
          </Label>
          <Field label="ลำดับการแสดง">
            <Input
              id={`sortOrder-${o?.id}`}
              name="sortOrder"
              type="number"
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
            <p className="text-muted-foreground text-sm">บันทึกแล้ว</p>
          )}
          <div className="flex gap-2">
            <Button
              type="submit"
              className="flex-1"
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
