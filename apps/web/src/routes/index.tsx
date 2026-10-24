import { queryOptions, useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";

import { api } from "@/api/client";
import type { Config } from "@/api/types.gen";

const config = queryOptions({
  queryKey: ["config"],
  queryFn: () => api<Config>("GET", "/config"),
});

export const Route = createFileRoute("/")({
  loader: ({ context }) => context.queryClient.ensureQueryData(config),
  component: Home,
});

function Home() {
  const { data } = useSuspenseQuery(config);
  return (
    <main className="container mx-auto max-w-3xl px-4 py-6">
      <h1 className="text-2xl font-semibold">แคมเปญห้องเรียน</h1>
      <dl className="mt-4 grid grid-cols-[auto_1fr] gap-x-4">
        <dt>liffId</dt>
        <dd>{data.liffId || "-"}</dd>
        <dt>devLogin</dt>
        <dd>{String(data.devLogin)}</dd>
      </dl>
    </main>
  );
}
