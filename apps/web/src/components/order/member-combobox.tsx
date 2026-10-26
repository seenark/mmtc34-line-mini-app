import { Button } from "@mmtc34/ui/components/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@mmtc34/ui/components/command";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@mmtc34/ui/components/popover";
import { ChevronsUpDown, UserPlus } from "lucide-react";
import { useState } from "react";

import { errorText } from "@/api/me";
import type { Member } from "@/api/types.gen";
import { useCreateUnlinkedMember } from "@/api/unlinked";

const details = (m: Member) =>
  [m.cohort && `รุ่น ${m.cohort}`, m.realName].filter(Boolean).join(" · ");

const matches = (m: Member, query: string) =>
  [m.nickname, m.realName, m.cohort].some((field) =>
    field.toLowerCase().includes(query),
  );

/**
 * Searchable Member picker for placing an Order on someone's behalf. The Super Admin and deleted
 * Members are not classmates and never appear. Typing a nickname nobody has offers to add an
 * Unlinked Member on the spot and selects them.
 */
export function MemberCombobox(props: {
  id: string;
  members: Member[];
  value: Member | null;
  onChange: (member: Member) => void;
}) {
  const [open, setOpen] = useState(false);
  const [query, setQuery] = useState("");
  const create = useCreateUnlinkedMember();

  const typed = query.trim();
  const found = props.members
    .filter((m) => !m.superAdmin && !m.deleted)
    .filter((m) => matches(m, typed.toLowerCase()));

  const choose = (member: Member) => {
    props.onChange(member);
    setOpen(false);
    setQuery("");
  };

  return (
    <Popover
      open={open}
      onOpenChange={(next) => {
        setOpen(next);
        if (!next) {
          setQuery("");
          create.reset();
        }
      }}
    >
      <PopoverTrigger
        render={
          <Button
            id={props.id}
            type="button"
            variant="outline"
            className="w-full justify-between font-normal"
          />
        }
      >
        <span className="min-w-0 truncate">
          {props.value ? props.value.nickname : "เลือกสมาชิก"}
          {props.value && !props.value.linked && (
            <span className="text-muted-foreground"> · ไม่มี LINE</span>
          )}
        </span>
        <ChevronsUpDown aria-hidden className="text-muted-foreground" />
      </PopoverTrigger>
      <PopoverContent
        align="start"
        className="w-(--anchor-width) max-w-[calc(100vw-1rem)] gap-0 p-0"
      >
        <Command shouldFilter={false}>
          <CommandInput
            value={query}
            onValueChange={setQuery}
            placeholder="ค้นหาชื่อเล่น ชื่อจริง หรือรุ่น"
          />
          <CommandList>
            <CommandEmpty className="text-muted-foreground">
              ยังไม่มีสมาชิก
            </CommandEmpty>
            {found.length > 0 && (
              <CommandGroup>
                {found.map((m) => (
                  <CommandItem
                    key={m.id}
                    value={String(m.id)}
                    data-checked={props.value?.id === m.id}
                    onSelect={() => choose(m)}
                  >
                    <span className="flex min-w-0 flex-col">
                      <span className="flex items-center gap-2">
                        <span className="min-w-0 truncate font-medium">
                          {m.nickname}
                        </span>
                        {!m.linked && (
                          <span className="text-muted-foreground shrink-0 rounded-full border px-2 text-xs">
                            ไม่มี LINE
                          </span>
                        )}
                      </span>
                      {details(m) && (
                        <span className="text-muted-foreground truncate text-xs">
                          {details(m)}
                        </span>
                      )}
                    </span>
                  </CommandItem>
                ))}
              </CommandGroup>
            )}
            {typed !== "" && found.length === 0 && (
              <CommandGroup>
                <CommandItem
                  value="create"
                  disabled={create.isPending}
                  onSelect={() =>
                    create.mutate(
                      {
                        cohort: "",
                        nickname: typed,
                        realName: "",
                        address: "",
                        phone: "",
                      },
                      { onSuccess: choose },
                    )
                  }
                >
                  <UserPlus aria-hidden />
                  <span className="min-w-0 [overflow-wrap:anywhere]">
                    เพิ่ม "{typed}" (ไม่มี LINE)
                  </span>
                </CommandItem>
              </CommandGroup>
            )}
          </CommandList>
          {create.error && (
            <p role="alert" className="text-destructive border-t p-2 text-sm">
              {errorText(create.error)}
            </p>
          )}
        </Command>
      </PopoverContent>
    </Popover>
  );
}
