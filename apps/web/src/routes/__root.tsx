import { Toaster } from "@mmtc34/ui/components/sonner";
import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

import { errorText } from "@/api/me";

import "../index.css";

export interface RouterAppContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterAppContext>()({
  component: () => (
    <>
      <Outlet />
      <Toaster />
    </>
  ),
  errorComponent: ({ error }) => (
    <p role="alert" className="text-destructive p-4">
      {errorText(error)}
    </p>
  ),
});
