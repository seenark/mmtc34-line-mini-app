import type { QueryClient } from "@tanstack/react-query";
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";

import { ApiError } from "@/api/client";

import "../index.css";

export interface RouterAppContext {
  queryClient: QueryClient;
}

export const Route = createRootRouteWithContext<RouterAppContext>()({
  component: Outlet,
  errorComponent: ({ error }) => (
    <p role="alert" className="p-4 text-red-700">
      {error instanceof ApiError
        ? error.message
        : "เกิดข้อผิดพลาด กรุณาลองใหม่อีกครั้ง"}
    </p>
  ),
});
