import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Suspense, useId, useState } from "react";

import {
  memberAddressesQuery,
  useCreateOrderOnBehalf,
} from "@/api/admin-orders";
import { adminCampaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { membersQuery } from "@/api/members";
import { DeliveryShip, type Member, type Option } from "@/api/types.gen";
import { Field } from "@/components/campaign/field";
import Loader from "@/components/loader";
import { MemberCombobox } from "@/components/order/member-combobox";
import { OrderForm } from "@/components/order/order-form";

export const Route = createFileRoute("/admin/campaigns/$campaignId/new-order")({
  staticData: { nav: { slot: "adminCampaign", label: "สั่งแทน", order: 30 } },
  loader: ({ context, params }) =>
    Promise.all([
      context.queryClient.ensureQueryData(
        adminCampaignQuery(params.campaignId),
      ),
      context.queryClient.ensureQueryData(membersQuery),
    ]),
  component: NewOrder,
});

function NewOrder() {
  const { campaignId } = Route.useParams();
  const id = useId();
  const { data: detail } = useSuspenseQuery(adminCampaignQuery(campaignId));
  const { data: members } = useSuspenseQuery(membersQuery);
  const [member, setMember] = useState<Member | null>(null);

  return (
    <div className="flex flex-col gap-4">
      <h2 className="text-base font-semibold">สั่งแทนสมาชิก</h2>
      <Field label="สมาชิก" required>
        <MemberCombobox
          id={`${id}-member`}
          members={members}
          value={member}
          onChange={setMember}
        />
      </Field>
      {member && (
        <Suspense fallback={<Loader />}>
          <MemberOrderForm
            key={member.id}
            campaignId={campaignId}
            member={member}
            options={detail.options}
          />
        </Suspense>
      )}
    </div>
  );
}

/** Split out so the Member's saved addresses load behind their own Suspense, not the whole page. */
function MemberOrderForm(props: {
  campaignId: string;
  member: Member;
  options: Option[];
}) {
  const { member } = props;
  const { data: addresses } = useSuspenseQuery(memberAddressesQuery(member.id));
  const place = useCreateOrderOnBehalf(props.campaignId);

  return (
    <OrderForm
      options={props.options}
      initial={{
        delivery: DeliveryShip,
        shipName: "",
        shipAddress: "",
        shipPhone: "",
        items: [],
      }}
      addresses={addresses}
      submitLabel="สั่งแทน"
      pending={place.isPending}
      error={place.error ? errorText(place.error) : undefined}
      onSubmit={(order) => place.mutate({ memberId: member.id, order })}
    />
  );
}
