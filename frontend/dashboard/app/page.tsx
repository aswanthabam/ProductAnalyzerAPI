"use client";
import { ProductInfo } from "@/utils/types";
import styles from "./page.module.css";
import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
const apiUrl: string = process.env.NEXT_PUBLIC_API_URL || "";
export default function Home() {
  const [products, setProducts] = useState<ProductInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();
  useEffect(() => {
    fetch(`${apiUrl}/api/product/list`)
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "success") {
          setProducts(data.data.products);
        } else {
          alert("Failed to fetch products");
        }
        setLoading(false);
      })
      .catch((err) => {
        alert("Failed to fetch products");
        setLoading(false);
      });
  }, []);

  return (
    <main className={styles.main}>
      {loading ? (
        <h2>Loading...</h2>
      ) : (
        <>
          <h2>Products</h2>
          <div className={styles.grid}>
            {products.map((product) => (
              <div
                onClick={() => {
                  router.push(`/product/?product=${product.code}`);
                }}
                key={product.code}
                className={styles.card}
              >
                <h3>{product.name}</h3>
                <p>Code: {product.code}</p>
                <p>Total Visits: {product.total_visits}</p>
                <p>Monthly Visits: {product.monthly_visits}</p>
              </div>
            ))}
          </div>
        </>
      )}
    </main>
  );
}
