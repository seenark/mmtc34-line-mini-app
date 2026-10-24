import liff from "@line/liff";
import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";

import { enforceStage, errorText, meQuery } from "@/api/me";
import { configQuery, devSignIn, signIn } from "@/lib/liff";

export const Route = createFileRoute("/login")({
  beforeLoad: ({ context }) => enforceStage(context.queryClient, "/login"),
  loader: ({ context }) => context.queryClient.ensureQueryData(configQuery),
  component: Login,
});

function Login() {
  const config = Route.useLoaderData();
  const qc = useQueryClient();
  const navigate = useNavigate();
  const [lineUserId, setLineUserId] = useState("");
  const login = useMutation({
    mutationFn: () =>
      config.devLogin ? devSignIn(lineUserId.trim()) : signIn(config),
    onSuccess: async (signedIn) => {
      if (!signedIn) return;
      await qc.resetQueries({ queryKey: meQuery.queryKey });
      await navigate({ to: "/" });
    },
  });

  // isInClient is a user-agent check that works before liff.init (which signIn runs).
  const { mutate } = login;
  useEffect(() => {
    if (!config.devLogin && liff.isInClient()) mutate();
  }, [config.devLogin, mutate]);

  return (
    <main className="mx-auto flex min-h-svh max-w-md flex-col justify-center gap-6 p-6">
      <h1 className="text-2xl font-semibold">แคมเปญห้องเรียน</h1>
      <form
        className="flex flex-col gap-4"
        onSubmit={(e) => {
          e.preventDefault();
          login.mutate();
        }}
      >
        {config.devLogin && (
          <div className="flex flex-col gap-1.5">
            <Label htmlFor="lineUserId">LINE userId (โหมดทดสอบ)</Label>
            <Input
              id="lineUserId"
              required
              autoComplete="off"
              value={lineUserId}
              onChange={(e) => setLineUserId(e.target.value)}
            />
          </div>
        )}
        {login.error && (
          <p role="alert" className="text-destructive text-sm">
            {errorText(login.error)}
          </p>
        )}
        <Button type="submit" size="lg" disabled={login.isPending}>
          {config.devLogin ? "เข้าสู่ระบบ" : "เข้าสู่ระบบด้วย LINE"}
        </Button>
      </form>
      {config.adminLogin && (
        <Link to="/admin-login" className="text-center text-sm underline">
          เข้าสู่ระบบผู้ดูแลระบบ
        </Link>
      )}
    </main>
  );
}
