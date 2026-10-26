import { Button } from "@mmtc34/ui/components/button";
import { Card, CardContent } from "@mmtc34/ui/components/card";
import { useSuspenseQuery } from "@tanstack/react-query";
import { createFileRoute } from "@tanstack/react-router";
import { Check, Loader2, LogOut } from "lucide-react";
import { useState } from "react";

import { errorText, meQuery, useLogout } from "@/api/me";
import { useDeleteAccount, useUpdateProfile } from "@/api/profile";
import { ProfileForm } from "@/components/profile-form";

export const Route = createFileRoute("/_app/profile")({
  staticData: { nav: { slot: "member", label: "โปรไฟล์", order: 30 } },
  component: Profile,
});

function Profile() {
  const { data: me } = useSuspenseQuery(meQuery);
  const update = useUpdateProfile();
  const logout = useLogout();
  if (!me?.member) return null;
  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-primary text-2xl font-extrabold">โปรไฟล์</h1>
      <section className="flex flex-col gap-3">
        <ProfileForm
          submitLabel="บันทึก"
          pending={update.isPending}
          error={update.error ? errorText(update.error) : undefined}
          initial={me.member}
          onSubmit={(input) => update.mutate(input)}
        />
        {update.isSuccess && (
          <p
            role="status"
            className="text-secondary flex items-center gap-1.5 text-sm font-medium"
          >
            <Check aria-hidden className="size-4 shrink-0" />
            บันทึกข้อมูลแล้ว
          </p>
        )}
        <p className="text-muted-foreground text-sm">
          คำสั่งซื้อที่สั่งไปแล้วจะใช้ที่อยู่เดิมตอนสั่งซื้อ การแก้ไขที่อยู่ไม่เปลี่ยนคำสั่งซื้อเก่า
        </p>
      </section>
      <Button
        variant="outline"
        size="lg"
        disabled={logout.isPending}
        onClick={() => logout.mutate()}
      >
        {logout.isPending ? (
          <Loader2 aria-hidden className="motion-safe:animate-spin" />
        ) : (
          <LogOut aria-hidden />
        )}
        ออกจากระบบ
      </Button>
      {!me.member.superAdmin && <DeleteAccount />}
    </div>
  );
}

function DeleteAccount() {
  const del = useDeleteAccount();
  const [confirming, setConfirming] = useState(false);
  return (
    <Card className="ring-destructive/30">
      <CardContent className="flex flex-col gap-3">
        <h2 className="text-destructive text-base font-semibold">ลบบัญชี</h2>
        <p className="text-sm">
          ข้อมูลติดต่อของคุณ (ชื่อ ที่อยู่ เบอร์โทร) จะถูกลบถาวร
          แต่ประวัติการชำระเงินและคำสั่งซื้อจะยังอยู่ในระบบ
        </p>
        {del.error && (
          <p
            role="alert"
            className="bg-destructive/10 text-destructive rounded-md px-3 py-2 text-sm"
          >
            {errorText(del.error)}
          </p>
        )}
        {confirming ? (
          <div className="flex flex-col gap-2">
            <p className="text-sm font-medium">ยืนยันลบบัญชี? ย้อนกลับไม่ได้</p>
            <div className="flex gap-2">
              <Button
                variant="destructive"
                size="lg"
                className="flex-1"
                disabled={del.isPending}
                onClick={() => del.mutate()}
              >
                {del.isPending && (
                  <Loader2 aria-hidden className="motion-safe:animate-spin" />
                )}
                ยืนยันลบบัญชี
              </Button>
              <Button
                variant="outline"
                size="lg"
                className="flex-1"
                disabled={del.isPending}
                onClick={() => setConfirming(false)}
              >
                ยกเลิก
              </Button>
            </div>
          </div>
        ) : (
          <Button
            variant="destructive"
            size="lg"
            onClick={() => setConfirming(true)}
          >
            ลบบัญชี
          </Button>
        )}
      </CardContent>
    </Card>
  );
}
