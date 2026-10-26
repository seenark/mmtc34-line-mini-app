import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { createRouter, RouterProvider } from "@tanstack/react-router";
import ReactDOM from "react-dom/client";
import { configQuery, initLiff } from "@/lib/liff";

import Loader from "./components/loader";
import { routeTree } from "./routeTree.gen";

const queryClient = new QueryClient();

const router = createRouter({
  routeTree,
  basepath: import.meta.env.BASE_URL,
  defaultPreload: "intent",
  scrollRestoration: true,
  defaultPendingComponent: () => <Loader />,
  context: { queryClient },
});

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("app");

if (!rootElement) {
  throw new Error("Root element not found");
}

async function start(root: HTMLElement) {
  const config = await queryClient.fetchQuery(configQuery).catch(() => null);
  if (config?.liffId && !config.devLogin) {
    await initLiff(config.liffId).catch((err) =>
      console.error("liff.init", err),
    );
  }
  ReactDOM.createRoot(root).render(
    <QueryClientProvider client={queryClient}>
      <RouterProvider router={router} />
    </QueryClientProvider>,
  );
}

if (!rootElement.innerHTML) {
  void start(rootElement);
}
