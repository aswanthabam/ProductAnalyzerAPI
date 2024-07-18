export type ProductInfo = {
  name: string;
  code: string;
  total_visits: number;
  monthly_visits: number;
};

export type RequestData = {
  ip: string | null;
  time: string | null;
  isp: string | null;
  postal: string | null;
  timezone: string | null;
  city: string | null;
  region: string | null;
  country: string | null;
  continent: string | null;
};

export type CountData = {
  name: string;
  count: number;
};

export type ProductDetailResponse = {
  name: string;
  latest_visits: RequestData[];
  countries: CountData[];
  regions: CountData[];
  cities: CountData[];
};
