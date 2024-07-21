import asyncio
from fastapi import Request
from fastapi import FastAPI
from fastapi.exceptions import RequestValidationError, HTTPException
from starlette.middleware.cors import CORSMiddleware
from starlette.responses import JSONResponse, RedirectResponse
from starlette.status import HTTP_401_UNAUTHORIZED, HTTP_422_UNPROCESSABLE_ENTITY

from utils.response import CustomResponse
from db.connection import Connection
from routers.api.product import product
from routers.api.user import user

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(product.router)
app.include_router(user.router)

asyncio.run(Connection().initialize_database())


async def custom_auth_exception_handler(request: Request, exc: HTTPException):
    return JSONResponse(
        CustomResponse.get_failure_response(message="Unauthorized! You need to login with an authorized credentials.",
                                            data={"error": exc.detail}),
        status_code=HTTP_401_UNAUTHORIZED)


async def exception_unprocessable_entry(request: Request, exc: RequestValidationError):
    return JSONResponse(
        CustomResponse.get_failure_response(message="Unprocessable entry!", data={"errors": exc.errors()}),
        status_code=HTTP_422_UNPROCESSABLE_ENTITY)


app.add_exception_handler(HTTP_401_UNAUTHORIZED, custom_auth_exception_handler)
app.add_exception_handler(RequestValidationError, exception_unprocessable_entry)


@app.get("/")
def root():
    return RedirectResponse(url="/dashboard")
