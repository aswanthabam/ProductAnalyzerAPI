from datetime import datetime

from pydantic import BaseModel


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
    time: datetime | None
    isp: str | None
    postal: str | None
    timezone: str | None
    city: str | None
    region: str | None
    country: str | None
    continent: str | None


class ProductInfoResponse(BaseModel):
    name: str
    latest_visits: list[RequestData]
    countries: list[dict]
    regions: list[dict]
    cities: list[dict]
