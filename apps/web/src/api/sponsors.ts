import {
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";

import { api } from "./client";
import type { SponsorInput, SponsorPosition } from "./types.gen";

export const sponsorsQuery = (campaignId: string) =>
  queryOptions({
    queryKey: ["campaigns", campaignId, "sponsors"],
    queryFn: () =>
      api<SponsorPosition[]>("GET", `/campaigns/${campaignId}/sponsors`),
  });

/** `id` set = update that Sponsor Position, unset = add one to the Campaign. */
export function useSaveSponsor(campaignId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id?: number; input: SponsorInput }) =>
      id === undefined
        ? api<SponsorPosition>(
            "POST",
            `/admin/campaigns/${campaignId}/sponsors`,
            input,
          )
        : api<SponsorPosition>("PUT", `/admin/sponsors/${id}`, input),
    onSuccess: () => qc.invalidateQueries(sponsorsQuery(campaignId)),
  });
}

export function useDeleteSponsor(campaignId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (sponsorId: number) =>
      api<void>("DELETE", `/admin/sponsors/${sponsorId}`),
    onSuccess: () => qc.invalidateQueries(sponsorsQuery(campaignId)),
  });
}
