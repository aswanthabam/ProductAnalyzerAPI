export type ProductInfo = {
  name: string;
  code: string;
  total_visits: number;
  monthly_visits: number;
};

export type RequestData = {
  ip: string;
  path: string;
  method: string;
  time: string;
  isp: string;
  postal: string;
  timezone: string;
  location: string;
  user_agent: string;
};

export type LocationData = {
  location: string;
  count: number;
};

export type ProductLocationsResponse = {
  top_cities: LocationData[];
  top_regions: LocationData[];
  top_countries: LocationData[];
};

export type ProductRequestsResponse = {
  page: number;
  total_pages: number;
  total: number;
  page_size: number;
  requests: RequestData[];
};
