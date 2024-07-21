"use client";

import Topbar from "@/components/Topbar/Topbar";

export function Main({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <Topbar />
      <div className="container">{children}</div>
    </>
  );
}
