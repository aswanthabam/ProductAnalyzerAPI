import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";
import { Main } from "./main";

const inter = Inter({ subsets: ["latin"] });

export const metadata: Metadata = {
  title: "Products by AVC",
  description: "Products dashboard by AVC",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <Main>{children}</Main>
      </body>
    </html>
  );
}
