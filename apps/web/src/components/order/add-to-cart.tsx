import { Button } from "@mmtc34/ui/components/button";
import { Check, ShoppingBag } from "lucide-react";
import { useId, useState } from "react";

import type { Option } from "@/api/types.gen";
import type { CartLine } from "@/lib/cart";

import { blankDetails, ItemFields } from "./item-fields";

export function AddToCart(props: {
  option: Option;
  onAdd: (line: CartLine) => void;
}) {
  const id = useId();
  const [details, setDetails] = useState(blankDetails);
  const [added, setAdded] = useState(false);
  return (
    <form
      className="mt-2 flex flex-col gap-3 border-t pt-3"
      onSubmit={(e) => {
        e.preventDefault();
        props.onAdd({ optionId: props.option.id, ...details });
        setDetails(blankDetails);
        setAdded(true);
      }}
    >
      <ItemFields
        id={id}
        option={props.option}
        value={details}
        onChange={(d) => {
          setDetails(d);
          setAdded(false);
        }}
      />
      <Button type="submit" size="lg">
        <ShoppingBag aria-hidden />
        เพิ่มลงตะกร้า
      </Button>
      {/* Always mounted so screen readers announce the text when it appears. */}
      <div role="status" className="empty:hidden">
        {added && (
          <p className="text-secondary flex items-center gap-1.5 text-sm font-medium">
            <Check aria-hidden className="size-4 shrink-0" />
            เพิ่มลงตะกร้าแล้ว
          </p>
        )}
      </div>
    </form>
  );
}
