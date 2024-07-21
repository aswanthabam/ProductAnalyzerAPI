from datetime import timedelta
from typing import Annotated

from fastapi import APIRouter, Form, Depends

from db.connection import Connection
from db.models import Users
from routers.api.user.response_models import AuthenticateResponse
from utils.auth import verify_password, get_password_hash, create_access_token, get_current_user, get_admin_user
from utils.response import CustomResponse

router = APIRouter(
    prefix="/api/user",
    tags=["product"],
    responses={404: {"description": "Not found"}},
)


@router.post('/register', description='Register user')
async def register(email: str = Form(...), password: str = Form(...)):
    connection = Connection()
    user = await connection.users.find_one({'email': email})
    if user is not None:
        return CustomResponse.get_failure_response('User already exists')
    await connection.users.insert_one({'email': email, 'password': get_password_hash(password), 'is_admin': False})
    token = create_access_token({'email': email}, expires_delta=timedelta(days=10))
    return CustomResponse.get_success_response('User registered successfully',
                                               data=AuthenticateResponse(access_token=token))


@router.post('/authenticate', description='Authenticate user')
async def authenticate(email: str = Form(...), password: str = Form(...)):
    connection = Connection()
    user = await connection.users.find_one({'email': email})
    if user is None:
        return CustomResponse.get_failure_response('User not found')
    if not verify_password(password, user.get('password')):
        return CustomResponse.get_failure_response('Invalid password')
    token = create_access_token({'email': email}, expires_delta=timedelta(days=10))
    return CustomResponse.get_success_response('User authenticated successfully',
                                               data=AuthenticateResponse(access_token=token))


@router.get('/is_admin', description='Check if user is admin')
async def me(current_user: Annotated[Users, Depends(get_admin_user)]):
    return current_user.is_admin
