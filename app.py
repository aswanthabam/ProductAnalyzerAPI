import asyncio
from fastapi import Request
from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware
from starlette.responses import JSONResponse, HTMLResponse, RedirectResponse
from starlette.staticfiles import StaticFiles
from starlette.status import HTTP_401_UNAUTHORIZED

from utils.response import CustomResponse
from db.connection import Connection
from routers.api.product import product

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(product.router)

app.mount("/dashboard", StaticFiles(directory="frontend/dashboard/build", html=True), name="home-dashboard")

asyncio.run(Connection().initialize_database())


async def custom_auth_exception_handler(request: Request, exc: Exception):
    return JSONResponse(CustomResponse.get_failure_response(message="Unauthorized!"))


app.add_exception_handler(HTTP_401_UNAUTHORIZED, custom_auth_exception_handler)


@app.get("/")
def root():
    return RedirectResponse(url="/dashboard")