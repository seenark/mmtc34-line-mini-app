import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/_app/")({
  component: Home,
});

function Home() {
  const { member } = Route.useRouteContext();
  return <h1 className="text-2xl font-semibold">สวัสดี {member.nickname}</h1>;
}
