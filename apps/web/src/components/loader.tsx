import { Loader2 } from "lucide-react";

export default function Loader() {
  return (
    <div
      role="status"
      className="text-secondary flex h-full items-center justify-center pt-8"
    >
      <Loader2 aria-hidden className="size-6 motion-safe:animate-spin" />
      <span className="sr-only">กำลังโหลด</span>
    </div>
  );
}
