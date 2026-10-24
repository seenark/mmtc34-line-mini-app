import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { Input } from "@mmtc34/ui/components/input";
import { useSuspenseQuery } from "@tanstack/react-query";
import { useState } from "react";

import { errorText } from "@/api/me";
import { membersQuery } from "@/api/members";
import { useDeleteSponsor, useSaveSponsor } from "@/api/sponsors";
import type { SponsorPosition } from "@/api/types.gen";
import { Field, selectClass } from "@/components/campaign/field";
import { ImageField } from "@/components/campaign/image-field";

export function SponsorsEditor(props: {
  campaignId: string;
  positions: SponsorPosition[];
}) {
  const [adding, setAdding] = useState(false);
  return (
    <section className="flex flex-col gap-3">
      {props.positions.map((p) => (
        <SponsorForm key={p.id} campaignId={props.campaignId} position={p} />
      ))}
      {adding ? (
        <SponsorForm
          campaignId={props.campaignId}
          onSaved={() => setAdding(false)}
        />
      ) : (
        <Button variant="outline" size="lg" onClick={() => setAdding(true)}>
          เพิ่มตำแหน่งสปอนเซอร์
        </Button>
      )}
    </section>
  );
}

/** `position` unset = a new Sponsor Position; saving calls `onSaved`. */
function SponsorForm(props: {
  campaignId: string;
  position?: SponsorPosition;
  onSaved?: () => void;
}) {
  const p = props.position;
  const save = useSaveSponsor(props.campaignId);
  const remove = useDeleteSponsor(props.campaignId);
  const { data: members } = useSuspenseQuery(membersQuery("approved"));
  const [logo, setLogo] = useState(p?.logo ?? "");
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
                id: p?.id,
                input: {
                  name: text("name"),
                  sponsorName: text("sponsorName"),
                  logo,
                  memberId: text("memberId") ? Number(text("memberId")) : null,
                  sortOrder: Math.trunc(Number(text("sortOrder"))),
                },
              },
              { onSuccess: props.onSaved },
            );
          }}
        >
          <Field label="ชื่อตำแหน่ง" required>
            <Input
              id={`name-${p?.id}`}
              name="name"
              required
              defaultValue={p?.name}
            />
          </Field>
          <Field label="ชื่อสปอนเซอร์">
            <Input
              id={`sponsorName-${p?.id}`}
              name="sponsorName"
              defaultValue={p?.sponsorName}
            />
          </Field>
          <ImageField label="โลโก้" value={logo} onChange={setLogo} />
          <Field label="สมาชิกผู้ดูแล">
            <select
              id={`memberId-${p?.id}`}
              name="memberId"
              className={selectClass}
              defaultValue={p?.memberId ?? ""}
            >
              <option value="">ไม่ระบุ</option>
              {members.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.nickname}
                </option>
              ))}
            </select>
          </Field>
          <Field label="ลำดับการแสดง">
            <Input
              id={`sortOrder-${p?.id}`}
              name="sortOrder"
              type="number"
              step={1}
              defaultValue={p?.sortOrder ?? 0}
            />
          </Field>
          {error && (
            <p role="alert" className="text-destructive text-sm">
              {errorText(error)}
            </p>
          )}
          {save.isSuccess && p && (
            <p className="text-muted-foreground text-sm">บันทึกแล้ว</p>
          )}
          <div className="flex gap-2">
            <Button
              type="submit"
              className="flex-1"
              size="lg"
              disabled={save.isPending}
            >
              {p ? "บันทึกตำแหน่ง" : "เพิ่มตำแหน่ง"}
            </Button>
            {p && (
              <Button
                type="button"
                variant="destructive"
                size="lg"
                disabled={remove.isPending}
                onClick={() => {
                  if (confirm(`ลบตำแหน่ง ${p.name}?`)) remove.mutate(p.id);
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
