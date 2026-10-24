import {
  type QueryClient,
  queryOptions,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";

import { shrinkImage } from "@/lib/image";

import { api } from "./client";
import type {
  Campaign,
  CampaignDetail,
  CampaignInput,
  ImageUpload,
  Option,
  OptionInput,
} from "./types.gen";

export const campaignsQuery = queryOptions({
  queryKey: ["campaigns"],
  queryFn: () => api<Campaign[]>("GET", "/campaigns"),
});

export const campaignQuery = (campaignId: string) =>
  queryOptions({
    queryKey: ["campaigns", campaignId],
    queryFn: () => api<CampaignDetail>("GET", `/campaigns/${campaignId}`),
  });

export const adminCampaignsQuery = queryOptions({
  queryKey: ["admin", "campaigns"],
  queryFn: () => api<Campaign[]>("GET", "/admin/campaigns"),
});

export const adminCampaignQuery = (campaignId: string) =>
  queryOptions({
    queryKey: ["admin", "campaigns", campaignId],
    queryFn: () => api<CampaignDetail>("GET", `/admin/campaigns/${campaignId}`),
  });

/** Admin edits change what Members see, so both views refetch. */
function invalidateCampaigns(qc: QueryClient) {
  return Promise.all([
    qc.invalidateQueries({ queryKey: ["campaigns"] }),
    qc.invalidateQueries({ queryKey: ["admin", "campaigns"] }),
  ]);
}

export function useCreateCampaign() {
  const qc = useQueryClient();
  const navigate = useNavigate();
  return useMutation({
    mutationFn: (input: CampaignInput) =>
      api<Campaign>("POST", "/admin/campaigns", input),
    onSuccess: async (campaign) => {
      await invalidateCampaigns(qc);
      await navigate({
        to: "/admin/campaigns/$campaignId",
        params: { campaignId: String(campaign.id) },
      });
    },
  });
}

export function useUpdateCampaign(campaignId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (input: CampaignInput) =>
      api<Campaign>("PUT", `/admin/campaigns/${campaignId}`, input),
    onSuccess: () => invalidateCampaigns(qc),
  });
}

/** `id` set = update that Option, unset = add one to the Campaign. */
export function useSaveOption(campaignId: string) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id?: number; input: OptionInput }) =>
      id === undefined
        ? api<Option>("POST", `/admin/campaigns/${campaignId}/options`, input)
        : api<Option>("PUT", `/admin/options/${id}`, input),
    onSuccess: () => invalidateCampaigns(qc),
  });
}

export function useDeleteOption() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (optionId: number) =>
      api<void>("DELETE", `/admin/options/${optionId}`),
    onSuccess: () => invalidateCampaigns(qc),
  });
}

/** Shrinks the photo client-side, uploads it, returns the stored image name. */
export async function uploadImage(file: File): Promise<string> {
  const form = new FormData();
  form.append("file", await shrinkImage(file), "image.jpg");
  const { name } = await api<ImageUpload>("POST", "/admin/images", form);
  return name;
}
