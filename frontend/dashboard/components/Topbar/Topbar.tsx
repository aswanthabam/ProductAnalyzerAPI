import Link from "next/link";
import styles from "./Topbar.module.css";
import { useEffect, useState } from "react";
import { IsAuthenticated } from "@/utils/api";
import { usePathname, useRouter } from "next/navigation";

export default function Topbar() {
  const [authenticated, setAuthenticated] = useState(false);
  const router = useRouter();
  const path = usePathname();
  useEffect(() => {
    setAuthenticated(IsAuthenticated());
  }, [path]);
  return (
    <div className={styles.topbar}>
      <div className={styles.topbarWrapper}>
        <div className={styles.topLeft}>
          <Link href="/" className={styles.title}>
            AVC Products Dashboard
          </Link>
        </div>
        <div className={styles.topRight}>
          {authenticated ? (
            <Link
              onClick={(e) => {
                e.preventDefault();
                localStorage.removeItem("token");
                router.push("/");
              }}
              href="/"
              className={styles.login}
            >
              Logout
            </Link>
          ) : (
            <Link href="/login" className={styles.login}>
              Login
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
