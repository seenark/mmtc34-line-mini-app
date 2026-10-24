import { queryOptions } from "@tanstack/react-query";

import { api } from "./client";
import type { RosterEntry } from "./types.gen";

export const rosterQuery = (campaignId: string) =>
  queryOptions({
    queryKey: ["campaigns", campaignId, "roster"],
    queryFn: () => api<RosterEntry[]>("GET", `/campaigns/${campaignId}/roster`),
  });
