from pydantic import BaseModel
from starlette.responses import Response


class CustomResponse(BaseModel):
    status: str = "success"
    message: str
    data: dict
    has_error: bool = False

    @staticmethod
    def get_failure_response(message: str, response: Response | None = None, status_code: int = 400,
                             data: BaseModel | dict | None = None):
        if response:
            response.status_code = status_code
        return CustomResponse(has_error=True,
                              data=data.model_dump() if isinstance(data, BaseModel) else (data if data else {}),
                              message=message,
                              status="failed").model_dump()

    @staticmethod
    def get_success_response(message: str, data: BaseModel | None = None, response: Response | None = None,
                             status_code: int = 200):
        if response:
            response.status_code = status_code
        return CustomResponse(has_error=False, data=data.model_dump() if data else {}, message=message,
                              status="success").model_dump()
