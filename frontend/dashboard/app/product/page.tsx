"use client";
import { useEffect, useState } from "react";
import styles from "./page.module.css";
import { ProductDetailResponse } from "@/utils/types";
import { useSearchParams } from "next/navigation";

export default function Product(p: any) {
  const params = useSearchParams();
  const product = params!.get("product");
  const [productInfo, setProductInfo] = useState<ProductDetailResponse | null>(
    null
  );
  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/product/${product}/info`)
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "success") {
          console.log(data.data);
          setProductInfo(data.data);
        } else {
          alert(data.message);
        }
      })
      .catch((err) => {
        alert("Failed to fetch product info");
      });
  }, [product]);
  return (
    <main className={styles.main}>
      {productInfo == null ? (
        <h2>Loading...</h2>
      ) : (
        <>
          <h2>{productInfo?.name} Product Information</h2>
          <div className={styles.grid}>
            <div>
              <h3>Top Countries</h3>
              <div className={styles.card}>
                {productInfo?.countries.map((country) => (
                  <p key={country.name}>
                    {country.name}: {country.count}
                  </p>
                ))}
              </div>
            </div>
            <div>
              <h3>Top Regions</h3>
              <div className={styles.card}>
                {productInfo?.regions.map((region) => (
                  <p key={region.name}>
                    {region.name}: {region.count}
                  </p>
                ))}
              </div>
            </div>
            <div>
              <h3>Top Cities</h3>
              <div className={styles.card}>
                {productInfo?.cities.map((city) => (
                  <p key={city.name}>
                    {city.name}: {city.count}
                  </p>
                ))}
              </div>
            </div>
          </div>
          <h3>Latest Visits</h3>
          <div className={styles.grid}>
            {productInfo?.latest_visits.map((visit) => (
              <div key={visit.ip} className={styles.card}>
                <p>IP: {visit.ip}</p>
                <p>Time: {visit.time}</p>
                <p>ISP: {visit.isp}</p>
                <p>Postal: {visit.postal}</p>
                <p>Timezone: {visit.timezone}</p>
                <p>City: {visit.city}</p>
                <p>Region: {visit.region}</p>
                <p>Country: {visit.country}</p>
                <p>Continent: {visit.continent}</p>
              </div>
            ))}
          </div>
        </>
      )}
    </main>
  );
}
