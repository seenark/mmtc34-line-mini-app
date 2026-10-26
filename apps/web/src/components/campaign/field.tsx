import { Label } from "@mmtc34/ui/components/label";
import type { ReactElement, ReactNode } from "react";

/** `aside` sits beside the control on its row (e.g. a button next to a select). */
export function Field(props: {
  label: string;
  required?: boolean;
  aside?: ReactNode;
  children: ReactElement<{ id: string }>;
}) {
  return (
    <div className="flex flex-col gap-1.5">
      <Label htmlFor={props.children.props.id} className="gap-1 font-medium">
        {props.label}
        {props.required && (
          <span aria-hidden className="text-destructive">
            *
          </span>
        )}
      </Label>
      {props.aside ? (
        <div className="flex items-center gap-2">
          {props.children}
          {props.aside}
        </div>
      ) : (
        props.children
      )}
    </div>
  );
}

/** The ui package has no Select; this matches Input's look. */
export const selectClass =
  "border-input h-10 w-full min-w-0 rounded-md border bg-transparent px-3 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-1 aria-invalid:ring-destructive/20";
