from datetime import datetime, time, date, UTC
from typing import Annotated
from uuid import uuid4
from pydantic import BaseModel, BeforeValidator, Field

PyObjectId = Annotated[str, BeforeValidator(str)]


class VisitData(BaseModel):
    path: str
    method: str
    time: str


class Visit(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    ip: str = Field(...)
    product: str = Field(...)
    city: str | None = Field(None)
    region: str | None = Field(None)
    country: str | None = Field(None)
    zipcode: str | None = Field(None)
    timezone: str | None = Field(None)
    lat: str | None = Field(None)
    lon: str | None = Field(None)
    isp: str | None = Field(None)
    org: str | None = Field(None)
    as_: str | None = Field(None)
    hosting: bool = Field(False)
    proxy: bool = Field(False)
    mobile: bool = Field(False)
    user_agent: str | None = Field(None)
    visits: list[VisitData] = Field([])
    date_: datetime = Field(..., default_factory=lambda: datetime.now(UTC))


class AdminUsers(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    email: str = Field(...)
    password: str = Field(...)


class Products(BaseModel):
    id: str = Field(default_factory=lambda: str(uuid4()))
    name: str = Field(...)
    code: str = Field(...)
    created_at: datetime = Field(..., default_factory=lambda: datetime.now(UTC))
