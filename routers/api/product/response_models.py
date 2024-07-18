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
