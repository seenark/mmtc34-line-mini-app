import { Button } from "@mmtc34/ui/components/button";
import { Input } from "@mmtc34/ui/components/input";
import { Label } from "@mmtc34/ui/components/label";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import {
  createFileRoute,
  Link,
  redirect,
  useNavigate,
} from "@tanstack/react-router";
import { useState } from "react";

import { adminSignIn, enforceStage, errorText, meQuery } from "@/api/me";
import { configQuery } from "@/lib/liff";

export const Route = createFileRoute("/admin-login")({
  beforeLoad: async ({ context }) => {
    await enforceStage(context.queryClient, "/login");
    const config = await context.queryClient.ensureQueryData(configQuery);
    if (!config.adminLogin) throw redirect({ to: "/login" });
  },
  component: AdminLogin,
});

function AdminLogin() {
  const qc = useQueryClient();
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const login = useMutation({
    mutationFn: () => adminSignIn(username.trim(), password),
    onSuccess: async () => {
      await qc.resetQueries({ queryKey: meQuery.queryKey });
      await navigate({ to: "/admin" });
    },
  });

  return (
    <main className="mx-auto flex min-h-svh max-w-md flex-col justify-center gap-6 p-6">
      <h1 className="text-2xl font-semibold">เข้าสู่ระบบผู้ดูแลระบบ</h1>
      <form
        className="flex flex-col gap-4"
        onSubmit={(e) => {
          e.preventDefault();
          login.mutate();
        }}
      >
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="username">ชื่อผู้ใช้</Label>
          <Input
            id="username"
            required
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <Label htmlFor="password">รหัสผ่าน</Label>
          <Input
            id="password"
            type="password"
            required
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        {login.error && (
          <p role="alert" className="text-destructive text-sm">
            {errorText(login.error)}
          </p>
        )}
        <Button type="submit" size="lg" disabled={login.isPending}>
          เข้าสู่ระบบ
        </Button>
      </form>
      <Link to="/login" className="text-center text-sm underline">
        กลับไปเข้าสู่ระบบด้วย LINE
      </Link>
    </main>
  );
}
