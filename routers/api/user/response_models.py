from pydantic import BaseModel

class AuthenticateResponse(BaseModel):
    access_token: str
