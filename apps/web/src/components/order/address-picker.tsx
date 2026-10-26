import { Input } from "@mmtc34/ui/components/input";
import { Textarea } from "@mmtc34/ui/components/textarea";
import { type ReactNode, useState } from "react";

import type { AddressChoice } from "@/api/types.gen";
import { Field } from "@/components/campaign/field";

type Selection = { kind: "card"; card: AddressChoice } | { kind: "new" };

const blank: AddressChoice = { name: "", address: "", phone: "" };

const same = (a: AddressChoice, b: AddressChoice) =>
  a.name.trim() === b.name.trim() &&
  a.address.trim() === b.address.trim() &&
  a.phone.trim() === b.phone.trim();

export interface ShipAddress {
  /** Saved triples to pick from; `current` is among them when it is known. */
  cards: AddressChoice[];
  current?: AddressChoice;
  selection: Selection;
  select: (selection: Selection) => void;
  /** The new-address fields. */
  draft: AddressChoice;
  setDraft: (draft: AddressChoice) => void;
  /** What the Order's ship fields will be. */
  ship: AddressChoice;
}

/**
 * The ship fields of one Order. `choices` are the owner's saved triples (profile first, then earlier
 * Orders); `current` is the Order's own snapshot when editing. The first card, or `current`, starts
 * selected; with no cards the Order gets a new address. Nothing here writes the profile.
 */
export function useShipAddress(
  choices: AddressChoice[],
  current?: AddressChoice,
): ShipAddress {
  const known = current && current.address.trim() !== "" ? current : undefined;
  const cards =
    known && !choices.some((c) => same(c, known))
      ? [known, ...choices]
      : choices;
  const [selection, setSelection] = useState<Selection>(() => {
    const first = known ?? cards[0];
    return first ? { kind: "card", card: first } : { kind: "new" };
  });
  const [draft, setDraft] = useState(blank);
  return {
    cards,
    current: known,
    selection,
    select: setSelection,
    draft,
    setDraft,
    ship: selection.kind === "card" ? selection.card : draft,
  };
}

function Choice(props: {
  name: string;
  checked: boolean;
  onSelect: () => void;
  children: ReactNode;
}) {
  return (
    <label className="group/addr has-focus-visible:ring-ring has-checked:border-secondary hover:bg-accent flex min-h-11 cursor-pointer items-start gap-3 rounded-xl border p-3 transition-colors has-focus-visible:ring-2">
      <input
        type="radio"
        name={props.name}
        checked={props.checked}
        onChange={props.onSelect}
        className="sr-only"
      />
      <span
        aria-hidden
        className="border-input bg-background group-has-checked/addr:border-secondary mt-0.5 size-4 shrink-0 rounded-full border transition-colors group-has-checked/addr:border-[5px]"
      />
      <span className="flex min-w-0 flex-col gap-0.5 text-sm">
        {props.children}
      </span>
    </label>
  );
}

export function AddressPicker(props: {
  id: string;
  address: ShipAddress;
  /** Own address: let the browser autofill. On someone's behalf it would fill the Admin's own. */
  autofill?: boolean;
}) {
  const { id, address, autofill = false } = props;
  const { cards, current, selection, select, draft, setDraft } = address;
  const showFields = cards.length === 0 || selection.kind === "new";
  const group = `${id}-address`;

  return (
    <div className="flex flex-col gap-3">
      {cards.length > 0 && (
        <fieldset className="flex flex-col gap-2">
          <legend className="mb-1 text-sm font-medium">เลือกที่อยู่จัดส่ง</legend>
          {cards.map((c) => (
            <Choice
              key={`${c.name}|${c.address}|${c.phone}`}
              name={group}
              checked={selection.kind === "card" && same(selection.card, c)}
              onSelect={() => select({ kind: "card", card: c })}
            >
              {current && same(current, c) && (
                <span className="text-secondary text-xs font-semibold">
                  ที่อยู่ปัจจุบัน
                </span>
              )}
              <span className="font-medium [overflow-wrap:anywhere]">
                {c.name}
                {c.phone && (
                  <span className="text-muted-foreground font-normal">
                    {" "}
                    · {c.phone}
                  </span>
                )}
              </span>
              <span className="text-muted-foreground [overflow-wrap:anywhere] whitespace-pre-wrap">
                {c.address}
              </span>
            </Choice>
          ))}
          <Choice
            name={group}
            checked={selection.kind === "new"}
            onSelect={() => select({ kind: "new" })}
          >
            <span className="font-medium">ที่อยู่ใหม่</span>
          </Choice>
        </fieldset>
      )}
      {showFields && (
        <div className="flex flex-col gap-3">
          <Field label="ชื่อผู้รับ" required>
            <Input
              id={`${id}-shipName`}
              required
              autoComplete={autofill ? "name" : "off"}
              value={draft.name}
              onChange={(e) => setDraft({ ...draft, name: e.target.value })}
            />
          </Field>
          <Field label="ที่อยู่จัดส่ง" required>
            <Textarea
              id={`${id}-shipAddress`}
              required
              autoComplete={autofill ? "street-address" : "off"}
              value={draft.address}
              onChange={(e) => setDraft({ ...draft, address: e.target.value })}
            />
          </Field>
          <Field label="เบอร์โทร" required>
            <Input
              id={`${id}-shipPhone`}
              type="tel"
              required
              autoComplete={autofill ? "tel" : "off"}
              value={draft.phone}
              onChange={(e) => setDraft({ ...draft, phone: e.target.value })}
            />
          </Field>
        </div>
      )}
    </div>
  );
}
