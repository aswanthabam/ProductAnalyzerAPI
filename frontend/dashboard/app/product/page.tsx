"use client";
import { useEffect, useState } from "react";
import styles from "./page.module.css";
import { useSearchParams } from "next/navigation";
import {
  ProductInfo,
  ProductLocationsResponse,
  RequestData,
} from "@/utils/types";
const apiUrl: string = process.env.NEXT_PUBLIC_API_URL || "";
export default function Product(p: any) {
  const params = useSearchParams();
  const product = params!.get("product");
  const [locationData, setLocationData] =
    useState<ProductLocationsResponse | null>(null);
  const [productInfo, setProductInfo] = useState<ProductInfo | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalRequests, setTotalRequests] = useState(0);
  const [requests, setRequests] = useState<RequestData[] | null>(null);
  useEffect(() => {
    fetch(`${apiUrl}/api/product/${product}/info`)
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
  useEffect(() => {
    if (productInfo == null) {
      return;
    }
    fetch(`${apiUrl}/api/product/${product}/locations`)
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "success") {
          console.log(data.data);
          setLocationData(data.data);
        } else {
          alert(data.message);
        }
      })
      .catch((err) => {
        alert("Failed to fetch top locations");
      });
  }, [productInfo]);
  useEffect(() => {
    if (productInfo == null) {
      return;
    }
    fetch(`${apiUrl}/api/product/${product}/requests?page=${page}`)
      .then((res) => res.json())
      .then((data) => {
        if (data.status == "success") {
          console.log(data.data);
          setTotalPages(data.data.total_pages);
          setTotalRequests(data.data.total);
          setRequests(data.data.requests);
        } else {
          alert(data.message);
        }
      })
      .catch((err) => {
        alert("Failed to fetch requests");
      });
  }, [productInfo, page]);
  return (
    <main className={styles.main}>
      <h2>{productInfo?.name} Product Information</h2>
      {!locationData ? (
        <h1>Loading...</h1>
      ) : (
        <div className={styles.grid}>
          <div>
            <h3>Top Countries</h3>
            <div className={styles.card}>
              <ul className={styles.ul}>
                {locationData?.top_countries.map((country) => (
                  <li key={country.location}>
                    <strong>{country.location}</strong> ({country.count})
                  </li>
                ))}
              </ul>
            </div>
          </div>
          <div>
            <h3>Top Regions</h3>
            <div className={styles.card}>
              <ul className={styles.ul}>
                {locationData.top_regions.map((region) => (
                  <li key={region.location}>
                    <strong>{region.location}</strong>: {region.count}
                  </li>
                ))}
              </ul>
            </div>
          </div>
          <div>
            <h3>Top Cities</h3>
            <div className={styles.card}>
              <ul className={styles.ul}>
                {locationData.top_cities.map((city) => (
                  <li key={city.location}>
                    <strong>{city.location}</strong>: {city.count}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      )}
      <h3>Latest Requests</h3>

      {!requests ? (
        <h1>Loading ...</h1>
      ) : (
        <>
          <div className={styles.pagination}>
            <button
              className={styles.btn}
              onClick={() => {
                setRequests(null);
                setPage(page - 1);
              }}
              disabled={page == 1}
            >
              Previous
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button
              className={styles.btn}
              onClick={() => {
                setRequests(null);
                setPage(page + 1);
              }}
              disabled={page == totalPages}
            >
              Next
            </button>
          </div>
          <div className={styles.grid + " " + styles.row}>
            {requests.map((visit) => (
              <div
                key={visit.ip}
                className={styles.card}
                onClick={(e) => {
                  e.currentTarget.classList.toggle(styles.active);
                }}
              >
                <div className={styles.cardHeader}>
                  <p>
                    <strong>
                      {visit.path} ({visit.method})
                    </strong>
                  </p>
                  <p>{visit.ip}</p>
                  <p>{new Date(visit.time).toLocaleString()}</p>
                  <p>{visit.location}</p>
                </div>
                <div className={styles.cardBody}>
                  <p>
                    <b>ISP:</b> {visit.isp}
                  </p>
                  <p>
                    <b>Postal:</b> {visit.postal || "N/A"}
                  </p>
                  <p>
                    <b>Timezone:</b> {visit.timezone}
                  </p>
                  <p>
                    <b>User Agent:</b> {visit.user_agent}
                  </p>
                </div>
              </div>
            ))}
          </div>
          <div className={styles.pagination}>
            <button
              className={styles.btn}
              onClick={() => {
                setRequests(null);
                setPage(page - 1);
              }}
              disabled={page == 1}
            >
              Previous
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button
              className={styles.btn}
              onClick={() => {
                setRequests(null);
                setPage(page + 1);
              }}
              disabled={page == totalPages}
            >
              Next
            </button>
          </div>
        </>
      )}
    </main>
  );
}
