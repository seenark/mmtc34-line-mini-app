import type { ReactNode } from "react";

/** Navy masthead over a white form column, shared by the LINE and Super Admin sign-in screens. */
export function AuthShell({
  heading,
  lede,
  children,
}: {
  heading: string;
  lede?: string;
  children: ReactNode;
}) {
  return (
    <div className="flex min-h-svh flex-col">
      <header className="bg-primary text-primary-foreground pt-[env(safe-area-inset-top)]">
        <div className="mx-auto w-full max-w-md px-6 pt-16 pb-10">
          <h1 className="text-3xl font-extrabold tracking-tight [overflow-wrap:anywhere]">
            {heading}
          </h1>
          {lede && (
            <p className="text-primary-muted-foreground mt-2 text-sm">{lede}</p>
          )}
        </div>
      </header>
      <main className="mx-auto flex w-full max-w-md flex-1 flex-col gap-6 px-6 py-8 pb-[calc(2rem+env(safe-area-inset-bottom))]">
        {children}
      </main>
    </div>
  );
}
