import { Button } from "@mmtc34/ui/components/button";
import { Checkbox } from "@mmtc34/ui/components/checkbox";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@mmtc34/ui/components/command";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@mmtc34/ui/components/dialog";
import { Label } from "@mmtc34/ui/components/label";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@mmtc34/ui/components/popover";
import { useQuery } from "@tanstack/react-query";
import { ChevronsUpDownIcon, Link2Icon, Loader2Icon } from "lucide-react";
import { useId, useState } from "react";
import { toast } from "sonner";

import { ApiError } from "@/api/client";
import { errorText } from "@/api/me";
import { membersQuery } from "@/api/members";
import type { Member, OrderLinkResult } from "@/api/types.gen";
import { orderLinkInfoQuery, useLinkOrder } from "@/api/unlinked";

/** cmdk's fuzzy scorer does poorly on Thai, so the picker matches plain substrings. */
const matchesSearch = (_value: string, search: string, keywords?: string[]) =>
  (keywords ?? []).join(" ").toLowerCase().includes(search.trim().toLowerCase())
    ? 1
    : 0;

function movedToast(unlinked: string, target: string, r: OrderLinkResult) {
  toast.success(`ย้ายคำสั่งซื้อแล้ว ${r.moved} รายการไปที่ ${target}`, {
    description: r.unlinkedRemoved
      ? `ลบ ${unlinked} ออกจากรายชื่อแล้ว เพราะไม่เหลือคำสั่งซื้อ`
      : `${unlinked} ยังอยู่ในรายชื่อ`,
  });
}

/** The body mounts only while the dialog is open, so the link info is fetched lazily and the form resets on every open. */
export function LinkOrderDialog({ orderId }: { orderId: number }) {
  const [open, setOpen] = useState(false);
  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger
        render={<Button variant="outline" size="sm" className="self-start" />}
      >
        <Link2Icon aria-hidden />
        ลิงก์กับบัญชี LINE
      </DialogTrigger>
      <DialogContent>
        <LinkOrderBody orderId={orderId} onDone={() => setOpen(false)} />
      </DialogContent>
    </Dialog>
  );
}

function LinkOrderBody(props: { orderId: number; onDone: () => void }) {
  const checkId = useId();
  const orderId = String(props.orderId);
  const info = useQuery({ ...orderLinkInfoQuery(orderId), retry: false });
  const members = useQuery(membersQuery);
  const [target, setTarget] = useState<Member | null>(null);
  const unlinked = info.data?.unlinkedMember.nickname;
  const link = useLinkOrder(orderId, (result) => {
    if (target && unlinked) movedToast(unlinked, target.nickname, result);
    props.onDone();
  });
  const [pickerOpen, setPickerOpen] = useState(false);
  const [allOrders, setAllOrders] = useState(true);

  const others = info.data?.otherOrders ?? 0;
  // Targets have a LINE account, are live, and are not the Super Admin.
  const candidates = (members.data ?? []).filter(
    (m) => m.linked && !m.deleted && !m.superAdmin,
  );

  function confirm() {
    if (!target || !unlinked) return;
    link.mutate({ memberId: target.id, allOrders: others > 0 && allOrders });
  }

  return (
    <>
      <DialogHeader>
        <DialogTitle>ลิงก์คำสั่งซื้อกับบัญชี LINE</DialogTitle>
        <DialogDescription>
          เลือกสมาชิกที่มีบัญชี LINE เพื่อย้ายคำสั่งซื้อของสมาชิกที่ยังไม่มี LINE ไปให้
        </DialogDescription>
      </DialogHeader>

      {/* a refetch after the move can 409; keep the form while the mutation finishes */}
      {!info.data && info.isPending && (
        <p className="text-muted-foreground text-sm" role="status">
          กำลังโหลด…
        </p>
      )}
      {!info.data && info.error && (
        <p role="alert" className="text-destructive text-sm">
          {info.error instanceof ApiError && info.error.status === 409
            ? "คำสั่งซื้อนี้ไม่ได้เป็นของสมาชิกที่ไม่มี LINE แล้ว"
            : errorText(info.error)}
        </p>
      )}

      {info.data && (
        <div className="flex flex-col gap-4">
          <Popover open={pickerOpen} onOpenChange={setPickerOpen}>
            <PopoverTrigger
              render={
                <Button
                  variant="outline"
                  className="w-full justify-between font-normal"
                  disabled={members.isPending || link.isPending}
                  aria-invalid={members.isError || undefined}
                />
              }
            >
              <span
                className={`min-w-0 truncate ${target ? "" : "text-muted-foreground"}`}
              >
                {members.isPending
                  ? "กำลังโหลดรายชื่อ…"
                  : target
                    ? target.nickname
                    : "เลือกสมาชิกที่มี LINE"}
              </span>
              <ChevronsUpDownIcon aria-hidden className="opacity-50" />
            </PopoverTrigger>
            <PopoverContent
              align="start"
              className="w-(--anchor-width) gap-0 p-0"
            >
              <Command filter={matchesSearch}>
                <CommandInput placeholder="ค้นหาชื่อเล่น ชื่อจริง หรือรุ่น" />
                <CommandList>
                  <CommandEmpty>ไม่พบสมาชิก</CommandEmpty>
                  <CommandGroup>
                    {candidates.map((m) => (
                      <CommandItem
                        key={m.id}
                        value={String(m.id)}
                        keywords={[m.nickname, m.realName, m.cohort]}
                        data-checked={target?.id === m.id}
                        onSelect={() => {
                          setTarget(m);
                          setPickerOpen(false);
                        }}
                      >
                        <span className="min-w-0 truncate">
                          {m.nickname}
                          <span className="text-muted-foreground">
                            {m.realName && ` · ${m.realName}`}
                            {m.cohort && ` · รุ่น ${m.cohort}`}
                          </span>
                        </span>
                      </CommandItem>
                    ))}
                  </CommandGroup>
                </CommandList>
              </Command>
            </PopoverContent>
          </Popover>
          {members.isError && (
            <p role="alert" className="text-destructive text-sm">
              {errorText(members.error)}
            </p>
          )}

          {target && (
            <p
              aria-live="polite"
              className="bg-muted rounded-md px-3 py-2 text-sm break-words"
            >
              ย้ายคำสั่งซื้อของ <strong className="font-semibold">{unlinked}</strong>{" "}
              ไปที่ <strong className="font-semibold">{target.nickname}</strong>
            </p>
          )}

          {others > 0 && (
            <div className="flex items-start gap-2">
              <Checkbox
                id={checkId}
                checked={allOrders}
                onCheckedChange={setAllOrders}
                disabled={link.isPending}
                className="mt-0.5"
              />
              <Label htmlFor={checkId} className="items-start">
                ย้ายคำสั่งซื้ออื่นของ {unlinked} อีก {others} รายการด้วย
              </Label>
            </div>
          )}
        </div>
      )}

      {link.error && (
        <p role="alert" className="text-destructive text-sm">
          {errorText(link.error)}
        </p>
      )}

      <DialogFooter>
        <DialogClose
          render={<Button variant="outline" disabled={link.isPending} />}
        >
          ยกเลิก
        </DialogClose>
        <Button
          onClick={confirm}
          disabled={!info.data || !target || link.isPending}
          aria-busy={link.isPending || undefined}
        >
          {link.isPending && (
            <Loader2Icon
              aria-hidden
              className="animate-spin motion-reduce:animate-none"
            />
          )}
          {link.isPending ? "กำลังย้าย…" : "ยืนยันการย้าย"}
        </Button>
      </DialogFooter>
    </>
  );
}
