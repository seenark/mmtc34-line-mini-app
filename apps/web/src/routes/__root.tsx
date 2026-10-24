import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

import { errorText } from "@/api/me";

import "../index.css";

export interface RouterAppContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterAppContext>()({
  component: Outlet,
  errorComponent: ({ error }) => (
    <p role="alert" className="p-4 text-red-700">
      {errorText(error)}
    </p>
  ),
});
