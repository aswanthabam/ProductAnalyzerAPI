import asyncio
from fastapi import Request
from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware
from starlette.responses import JSONResponse, HTMLResponse
from starlette.status import HTTP_401_UNAUTHORIZED

from utils.response import CustomResponse
from db.connection import Connection
from routers.product import product
app = FastAPI()

app.include_router(product.router)
asyncio.run(Connection().initialize_database())

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


async def custom_auth_exception_handler(request: Request, exc: Exception):
    return JSONResponse(CustomResponse.get_failure_response(message="Unauthorized!"))


app.add_exception_handler(HTTP_401_UNAUTHORIZED, custom_auth_exception_handler)


@app.get("/")
def root():
    return HTMLResponse(content="<html><h1>Haa shit! My Code is working.</h1></html>")