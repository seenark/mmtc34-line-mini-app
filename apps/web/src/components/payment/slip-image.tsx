import { useState } from "react";

import { slipUrl } from "@/api/payments";

/** Slip thumbnail; tapping opens the full image in an overlay. */
export function SlipImage({ name }: { name: string }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        aria-label="ดูสลิปขนาดเต็ม"
        className="block"
      >
        <img
          src={slipUrl(name)}
          alt="สลิป"
          className="bg-muted max-h-64 w-full object-contain"
        />
      </button>
      {open && (
        <button
          type="button"
          onClick={() => setOpen(false)}
          aria-label="ปิด"
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/80 p-2"
        >
          <img
            src={slipUrl(name)}
            alt="สลิป"
            className="max-h-full max-w-full object-contain"
          />
        </button>
      )}
    </>
  );
}
