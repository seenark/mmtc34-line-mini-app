import { Label } from "@mmtc34/ui/components/label";
import type { ReactElement } from "react";

export function Field(props: {
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

/** The ui package has no Select; this matches Input's look. */
export const selectClass =
  "border-input h-8 w-full rounded-none border bg-transparent px-2 text-xs outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-1";
