import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { Textarea } from "@mmtc34/ui/components/textarea";
import type { ReactElement } from "react";

import type { ProfileInput } from "@/api/types.gen";

/** Cohort, nickname and real name are required; address and phone are optional. */
export function ProfileForm(props: {
  submitLabel: string;
  pending: boolean;
  error?: string;
  onSubmit: (input: ProfileInput) => void;
}) {
  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={(e) => {
        e.preventDefault();
        const f = new FormData(e.currentTarget);
        const text = (name: keyof ProfileInput) =>
          String(f.get(name) ?? "").trim();
        props.onSubmit({
          cohort: text("cohort"),
          nickname: text("nickname"),
          realName: text("realName"),
          address: text("address"),
          phone: text("phone"),
        });
      }}
    >
      <Field label="รุ่น" required>
        <Input id="cohort" name="cohort" required autoComplete="off" />
      </Field>
      <Field label="ชื่อเล่น" required>
        <Input id="nickname" name="nickname" required autoComplete="nickname" />
      </Field>
      <Field label="ชื่อจริง" required>
        <Input id="realName" name="realName" required autoComplete="name" />
      </Field>
      <Field label="ที่อยู่">
        <Textarea id="address" name="address" autoComplete="street-address" />
      </Field>
      <Field label="เบอร์โทร">
        <Input id="phone" name="phone" type="tel" autoComplete="tel" />
      </Field>
      {props.error && (
        <p role="alert" className="text-destructive text-sm">
          {props.error}
        </p>
      )}
      <Button type="submit" size="lg" disabled={props.pending}>
        {props.submitLabel}
      </Button>
    </form>
  );
}

function Field(props: {
  label: string;
  required?: boolean;
  children: ReactElement<{ id: string }>;
}) {
  return (
    <div className="flex flex-col gap-1.5">
      <Label htmlFor={props.children.props.id}>
        {props.label}
        {props.required && <span className="text-destructive">*</span>}
      </Label>
      {props.children}
    </div>
  );
}
