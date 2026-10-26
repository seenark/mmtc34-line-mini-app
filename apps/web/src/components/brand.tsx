import { cn } from "@mmtc34/ui/lib/utils";

/** The MMTC Events wordmark. Inherits text colour so it works on paper and on navy. */
export function Brand({ className }: { className?: string }) {
  return (
    <span
      className={cn(
        "flex min-w-0 items-center gap-2 text-base font-extrabold tracking-tight",
        className,
      )}
    >
      <span
        aria-hidden
        className="bg-secondary size-2.5 shrink-0 rounded-[3px]"
      />
      <span className="truncate">MMTC Events</span>
    </span>
  );
}
