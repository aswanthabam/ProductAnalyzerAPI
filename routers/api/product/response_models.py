from datetime import datetime

from pydantic import BaseModel


class ProductInfoResponse(BaseModel):
    name: str
    code: str


class ProductResponse(BaseModel):
    name: str
    code: str
    created_at: datetime
    total_visits: int
    monthly_visits: int


class ListProductsResponse(BaseModel):
    products: list[ProductResponse]


class RequestData(BaseModel):
    ip: str | None
    path: str | None
    method: str | None
    time: datetime | None
    isp: str | None
    postal: str | None
    timezone: str | None
    location: str | None
    user_agent: str | None


class ProductRequestsResponse(BaseModel):
    page: int
    total_pages: int
    total: int
    page_size: int
    requests: list[RequestData]


class LocationData(BaseModel):
    location: str
    count: int

class ProductLocationsResponse(BaseModel):
    top_cities: list[LocationData]
    top_countries: list[LocationData]
    top_regions: list[LocationData]
