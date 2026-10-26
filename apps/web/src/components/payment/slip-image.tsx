import { Button } from "@mmtc34/ui/components/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogTitle,
} from "@mmtc34/ui/components/dialog";
import { useState } from "react";

import { slipUrl } from "@/api/payments";

/** Slip thumbnail; tapping opens the full image in a dialog. */
export function SlipImage({ name }: { name: string }) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <button
        type="button"
        onClick={() => setOpen(true)}
        aria-label="ดูสลิปขนาดเต็ม"
        className="bg-muted hover:bg-accent focus-visible:ring-ring block w-full rounded-md border transition-colors outline-none focus-visible:ring-2 active:opacity-80"
      >
        <img
          src={slipUrl(name)}
          alt="สลิป"
          className="max-h-64 w-full rounded-md object-contain"
        />
      </button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent
          showCloseButton={false}
          className="max-h-[calc(100svh-2rem)] grid-rows-[auto_minmax(0,1fr)_auto] sm:max-w-xl"
        >
          <DialogTitle className="text-base font-semibold">สลิป</DialogTitle>
          <img
            src={slipUrl(name)}
            alt="สลิป"
            className="bg-muted max-h-full min-h-0 w-full rounded-md object-contain"
          />
          <DialogClose render={<Button variant="outline" size="lg" />}>
            ปิด
          </DialogClose>
        </DialogContent>
      </Dialog>
    </>
  );
}
