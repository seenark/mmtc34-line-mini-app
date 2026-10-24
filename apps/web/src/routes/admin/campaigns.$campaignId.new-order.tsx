import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { useId, useState } from "react";

import { useCreateOrderOnBehalf } from "@/api/admin-orders";
import { adminCampaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { membersQuery } from "@/api/members";
import { DeliveryShip } from "@/api/types.gen";
import { Field, selectClass } from "@/components/campaign/field";
import { OrderForm } from "@/components/order/order-form";

export const Route = createFileRoute("/admin/campaigns/$campaignId/new-order")({
  staticData: { nav: { slot: "adminCampaign", label: "สั่งแทน", order: 30 } },
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(
        adminCampaignQuery(params.campaignId),
      ),
      context.queryClient.ensureQueryData(membersQuery("approved")),
    ]),
  component: NewOrder,
});

function NewOrder() {
  const { campaignId } = Route.useParams();
  const id = useId();
  const { data: detail } = useSuspenseQuery(adminCampaignQuery(campaignId));
  const { data: members } = useSuspenseQuery(membersQuery("approved"));
  const place = useCreateOrderOnBehalf(campaignId);
  const [memberId, setMemberId] = useState("");
  const member = members.find((m) => String(m.id) === memberId);

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-xl font-semibold">สั่งแทนสมาชิก</h2>
      <Field label="สมาชิก" required>
        <select
          id={`${id}-member`}
          className={selectClass}
          value={memberId}
          onChange={(e) => setMemberId(e.target.value)}
        >
          <option value="">เลือกสมาชิก</option>
          {members.map((m) => (
            <option key={m.id} value={m.id}>
              {m.nickname}
              {m.realName && ` (${m.realName})`}
              {!m.linked && " · ไม่มี LINE"}
            </option>
          ))}
        </select>
      </Field>
      {member && (
        <OrderForm
          key={member.id}
          options={detail.options}
          initial={{
            delivery: DeliveryShip,
            shipName: member.realName,
            shipAddress: member.address,
            shipPhone: member.phone,
            items: [],
          }}
          submitLabel="สั่งแทน"
          pending={place.isPending}
          error={place.error ? errorText(place.error) : undefined}
          onSubmit={(order) => place.mutate({ memberId: member.id, order })}
        />
      )}
    </div>
  );
}
