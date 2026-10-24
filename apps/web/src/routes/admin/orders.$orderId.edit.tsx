import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute, Link } from "@tanstack/react-router";

import { useEditOrder } from "@/api/admin-orders";
import { adminCampaignQuery } from "@/api/campaigns";
import { errorText } from "@/api/me";
import { orderQuery } from "@/api/orders";
import { OrderForm } from "@/components/order/order-form";

export const Route = createFileRoute("/admin/orders/$orderId/edit")({
  staticData: {
    nav: { slot: "orderAction", label: "แก้ไขคำสั่งซื้อ", order: 30, when: "edit" },
  },
  loader: async ({ context, params }) => {
    const order = await context.queryClient.ensureQueryData(
      orderQuery(params.orderId),
    );
    await context.queryClient.ensureQueryData(
      adminCampaignQuery(String(order.campaignId)),
    );
  },
  component: EditOrder,
});

function EditOrder() {
  const { orderId } = Route.useParams();
  const { data: order } = useSuspenseQuery(orderQuery(orderId));
  const { data: detail } = useSuspenseQuery(
    adminCampaignQuery(String(order.campaignId)),
  );
  const save = useEditOrder(orderId);
  const price = new Map(detail.options.map((o) => [o.id, o.price]));

  return (
    <div className="flex flex-col gap-4">
      <Link
        to="/orders/$orderId"
        params={{ orderId }}
        className="text-sm underline"
      >
        ← คำสั่งซื้อ #{order.id}
      </Link>
      <h1 className="text-2xl font-semibold">แก้ไขคำสั่งซื้อ #{order.id}</h1>
      <OrderForm
        options={detail.options}
        initial={{
          delivery: order.delivery,
          shipName: order.shipName,
          shipAddress: order.shipAddress,
          shipPhone: order.shipPhone,
          items: order.items.map((i) => ({
            id: i.id,
            optionId: i.optionId,
            size: i.size,
            printName: i.printName,
            number: i.number,
            quantity: i.quantity,
            // a no-price Option's amount is stored as the unit price
            amount: price.get(i.optionId) === null ? i.unitPrice : null,
          })),
        }}
        submitLabel="บันทึก"
        pending={save.isPending}
        error={save.error ? errorText(save.error) : undefined}
        onSubmit={(input) => save.mutate(input)}
      />
    </div>
  );
}
