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
      <img
        src={`${import.meta.env.BASE_URL}icon-192.png`}
        alt=""
        width={28}
        height={28}
        className="size-7 shrink-0"
      />
      <span className="truncate">MMTC Events</span>
    </span>
  );
}
