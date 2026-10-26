import { cn } from "@mmtc34/ui/lib/utils";
import { Link } from "@tanstack/react-router";

import { type Mode, recordMode } from "@/lib/admin-mode";

const MODES = [
  { mode: "member", to: "/", label: "โหมดสมาชิก" },
  { mode: "admin", to: "/admin", label: "โหมดแอดมิน" },
] as const;

const TONES = {
  light: {
    track: "bg-muted",
    on: "bg-secondary text-secondary-foreground",
    off: "text-muted-foreground hover:text-foreground",
  },
  navy: {
    track: "bg-primary-raised",
    on: "bg-background text-primary",
    off: "text-primary-muted-foreground hover:text-primary-foreground",
  },
};

/** Two-state switch between the member shell and the admin shell. `onNavy` is for the admin header. */
export function ModeSwitch({
  current,
  onNavy = false,
}: {
  current: Mode;
  onNavy?: boolean;
}) {
  const tone = onNavy ? TONES.navy : TONES.light;
  return (
    <nav
      aria-label="โหมด"
      className={cn("flex shrink-0 rounded-full p-0.5", tone.track)}
    >
      {MODES.map(({ mode, to, label }) => {
        const selected = mode === current;
        return (
          <Link
            key={mode}
            to={to}
            activeOptions={{ exact: true }}
            onClick={() => recordMode(mode)}
            aria-current={selected ? "true" : undefined}
            className={cn(
              "flex h-9 items-center rounded-full px-3 text-sm font-semibold whitespace-nowrap transition-colors",
              selected ? tone.on : tone.off,
            )}
          >
            {label}
          </Link>
        );
      })}
    </nav>
  );
}
