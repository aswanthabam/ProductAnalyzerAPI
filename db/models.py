from datetime import datetime
from typing import Optional, Annotated
from uuid import uuid4

from pydantic import BaseModel, BeforeValidator, Field

DEFAULT_DAILY_LIMIT = (100 * 1024 * 1024)  # 100 MB

PyObjectId = Annotated[str, BeforeValidator(str)]


class Products(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    name: str = Field(...)
    code: str = Field(...)
    created_at: datetime = Field(..., default_factory=datetime.utcnow)


class Users(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    ip: str = Field(...)
    visit_count: int = Field(default=1)
    last_visit: datetime = Field(..., default_factory=datetime.utcnow)


class ProductVisits(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    product: str = Field(...)
    user: str = Field(...)
    city: str | None = Field(None)
    latitude: str | None = Field(None)
    longitude: str | None = Field(None)
    postal: str | None = Field(None)
    isp: str | None = Field(None)
    time: datetime = Field(..., default_factory=datetime.utcnow)


class Cities(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    name: str = Field(...)
    region: str = Field(...)
    timezone: str = Field(...)


class Regions(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    name: str = Field(...)
    code: str = Field(...)
    country: str = Field(...)


class Countries(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    name: str = Field(...)
    code: str = Field(...)
    continent_code: str = Field(...)
