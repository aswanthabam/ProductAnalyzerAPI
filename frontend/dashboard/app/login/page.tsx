"use client";
import { ProductInfo } from "@/utils/types";
import styles from "./page.module.css";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Get, Post } from "@/utils/api";
const apiUrl: string = process.env.NEXT_PUBLIC_API_URL || "";
export default function Login() {
  const [username, setUsername] = useState<string | null>(null);
  const [password, setPassword] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const login = async () => {
    if (
      username == null ||
      password == null ||
      username == "" ||
      password == ""
    ) {
      alert("Please enter a username and password");
      return;
    }
    setLoading(true);
    Post(`${apiUrl}/api/user/authenticate`, {
      headers: { "content-type": "application/x-www-form-urlencoded" },
      body: new URLSearchParams({ email: username, password: password }),
    })
      .then(async (res) => {
        setLoading(false);
        var data = await res.json();
        if (data.status == "success") {
          localStorage.setItem("token", data.data.access_token);
          alert("Logged in successfully");
          router.push("/");
        } else {
          alert(data.message);
        }
      })
      .catch((err) => {
        setLoading(false);
        console.log(err);
        alert("Failed to login");
      });
  };
  return (
    <div className={styles.login}>
      <h2>Login</h2>
      <div className={styles.form}>
        <label htmlFor="username">Username</label>
        <input
          onChange={(e) => {
            setUsername(e.target.value);
          }}
          className={styles.input}
          type="text"
          id="username"
          name="username"
        />
        <label htmlFor="password">Password</label>
        <input
          onChange={(e) => {
            setPassword(e.target.value);
          }}
          type="password"
          className={styles.input}
          id="password"
          name="password"
        />
        <button
          className={styles.button}
          onClick={login}
          disabled={loading}
          type="submit"
        >
          {loading ? "Signing In..." : "Login"}
        </button>
      </div>
    </div>
  );
}
