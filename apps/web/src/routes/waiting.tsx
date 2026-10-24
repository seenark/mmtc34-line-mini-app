import { Button } from "@mmtc34/ui/components/button";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useEffect } from "react";

import { enforceStage, meQuery, useLogout } from "@/api/me";
import { MemberApproved } from "@/api/types.gen";

export const Route = createFileRoute("/waiting")({
  beforeLoad: ({ context }) => enforceStage(context.queryClient, "/waiting"),
  component: Waiting,
});

function Waiting() {
  const navigate = useNavigate();
  const logout = useLogout();
  const { data: me } = useQuery({ ...meQuery, refetchInterval: 15_000 });
  const approved = me?.member?.status === MemberApproved;

  useEffect(() => {
    if (approved) void navigate({ to: "/" });
  }, [approved, navigate]);

  return (
    <main className="mx-auto flex min-h-svh max-w-md flex-col items-center justify-center gap-4 p-6 text-center">
      <h1 className="text-2xl font-semibold">รอแอดมินอนุมัติ</h1>
      <p className="text-muted-foreground text-sm">
        ส่งใบสมัครเรียบร้อยแล้ว หน้านี้จะอัปเดตเองเมื่อได้รับอนุมัติ
      </p>
      <Button
        variant="outline"
        onClick={() => logout.mutate()}
        disabled={logout.isPending}
      >
        ออกจากระบบ
      </Button>
    </main>
  );
}
