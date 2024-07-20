import math
from datetime import datetime, UTC, timedelta

import decouple
import pymongo.errors
import requests
from fastapi import APIRouter, Form
from starlette.requests import Request

from db.connection import Connection
from db.models import Products, VisitData, Visit
from routers.api.product.response_models import ProductResponse, ListProductsResponse, RequestData, \
    ProductRequestsResponse, LocationData, ProductLocationsResponse, ProductInfoResponse
from utils.response import CustomResponse

router = APIRouter(
    prefix="/api/product",
    tags=["product"],
    responses={404: {"description": "Not found"}},
)
GENERAL_PRODUCT_LIMIT_PER_HOUR = 60
PASSWORD = decouple.config('PASSWORD')


class IPAPI:
    def __init__(self, data: dict):
        self.status = data.get('status')
        self.continent = data.get('continent')
        self.country = data.get('country')
        self.regionName = data.get('regionName')
        self.city = data.get('city')
        self.zip = data.get('zip')
        self.lat = data.get('lat')
        self.lon = data.get('lon')
        self.lat = str(self.lat) if self.lat else None
        self.lon = str(self.lon) if self.lon else None
        self.timezone = data.get('timezone')
        self.isp = data.get('isp')
        self.org = data.get('org')
        self.as_ = data.get('as')
        self.asname = data.get('asname')
        self.mobile = data.get('mobile', False)
        self.proxy = data.get('proxy', False)
        self.hosting = data.get('hosting', False)

    def to_model_dict(self):
        return {
            'country': self.country if self.country else "Unknown",
            'region': self.regionName if self.regionName else "Unknown",
            'city': self.city if self.city else "Unknown",
            'zipcode': self.zip,
            'lat': self.lat,
            'lon': self.lon,
            'timezone': self.timezone,
            'isp': self.isp if self.isp else "Unknown",
            'org': self.org if self.org else "Unknown",
            'as_': self.as_ if self.as_ else "Unknown",
            'mobile': self.mobile,
            'proxy': self.proxy,
            'hosting': self.hosting
        }

    @staticmethod
    def from_ip_address(ip: str):
        url = f"http://ip-api.com/json/{ip}?fields=status,continent,country,regionName,city,zip,lat,lon,timezone,isp,org,as,asname,mobile,proxy,hosting"
        res = requests.get(url)
        data = res.json()
        if not data.get('error') and data.get('status') == "success":
            return IPAPI(data)
        return IPAPI({})


@router.post('/create/', description="Create a product")
async def create_product(code: str = Form(...), name: str = Form(...), password: str = Form(...)):
    if password != PASSWORD:
        return CustomResponse.get_failure_response("Unauthorized!")
    connection = Connection()
    try:
        await connection.products.insert_one(Products(code=code, name=name).model_dump())
    except pymongo.errors.DuplicateKeyError as e:
        return CustomResponse.get_failure_response("Product already exists.")
    return CustomResponse.get_success_response("Product Created")


@router.get('/{code}/info/', description="Get information about a project")
async def product_info(code: str):
    connection = Connection()
    product = await connection.products.find_one({'code': code})
    if not product:
        return CustomResponse.get_failure_response("No product with the code found!")
    product = Products(**product)
    return CustomResponse.get_success_response(f"{code} info",
                                               data=ProductInfoResponse(name=product.name, code=product.code))


@router.get('/{code}/visit/', description="Visit a product")
async def visit(code: str, request: Request, path: str = "/", method: str = "GET"):
    connection = Connection()

    # check if the product exists

    product = await connection.products.find_one({'code': code})
    if not product:
        return CustomResponse.get_failure_response("Product not found.")
    product = Products(**product)

    # finds the ip address of the user

    x_forwarded_for = request.headers.get("x-forwarded-for")
    if x_forwarded_for:
        client_ip = x_forwarded_for.split(",")[0]
    else:
        client_ip = request.client.host
    if not client_ip:
        return CustomResponse.get_failure_response("Error, trust me.. its not from our side.")

    # Finds the user agent of the user

    user_agent = request.headers.get("user-agent", "Unknown")

    # check if the user has visited the product the same day, if visited update the visit time
    # else creates a new visit

    time = datetime.now(UTC)
    date = time.date()

    start_of_today = datetime.combine(datetime.today(), datetime.min.time())
    start_of_tomorrow = start_of_today + timedelta(days=1)
    prev_visits = await connection.visits.find_one({'ip': client_ip, 'date_': {
        '$gte': start_of_today,
        '$lt': start_of_tomorrow
    }, 'product': product.id})  # Finds visit of the user in the same day
    if prev_visits:
        cur_visit = VisitData(
            path=path,
            method=method,
            time=time.time().strftime("%H:%M:%S")
        )
        await connection.visits.update_one({"id": prev_visits.get('id')}, {'$push': {'visits': cur_visit.model_dump()}})
    else:
        data = {}
        if ip_info := IPAPI.from_ip_address(client_ip):
            data = ip_info.to_model_dict()
        data.update({
            'ip': client_ip,
            'product': product.id,
            'date_': time,
            'user_agent': user_agent,
            'visits': [VisitData(path=path, method=method, time=time.time().strftime("%H:%M:%S")).model_dump()]
        })
        data = Visit(**data)
        await connection.visits.insert_one(data.model_dump())
    return CustomResponse.get_success_response("OK")


@router.get('/list', description="List available products")
async def list_products():
    connection = Connection()
    products = await connection.products.find().to_list(length=None)
    ps = []
    now = datetime.now(UTC)
    start_of_month = datetime(now.year, now.month, 1)
    if now.month == 12:
        start_of_next_month = datetime(now.year + 1, 1, 1)
    else:
        start_of_next_month = datetime(now.year, now.month + 1, 1)
    for product in products:
        product = Products(**product)
        pipeline = [
            {
                '$match': {
                    'product': product.id
                }
            },
            {
                '$project': {
                    'v_len': {
                        '$size': '$visits'
                    }
                }
            },
            {
                '$group': {
                    '_id': None,
                    'total_count': {
                        '$sum': '$v_len'
                    }
                }
            }
        ]
        visit_count = (await connection.visits.aggregate(pipeline).to_list(None))[0].get('total_count', 0)
        pipeline[0]['$match']['date_'] = {'$gte': start_of_month, '$lt': start_of_next_month}
        monthly_visit = (await connection.visits.aggregate(pipeline).to_list(None))[0].get('total_count', 0)
        ps.append(
            ProductResponse(name=product.name, code=product.code, created_at=product.created_at,
                            total_visits=visit_count, monthly_visits=monthly_visit))
    return CustomResponse.get_success_response("", data=ListProductsResponse(products=ps))


@router.get('/{code}/requests', description="Get information about a project")
async def product_requests(code: str, page: int = 1, page_size: int = 10):
    connection = Connection()
    product = await connection.products.find_one({'code': code})
    if not product:
        return CustomResponse.get_failure_response("No product with the code found!")
    product = Products(**product)
    skip_count = (page - 1) * page_size
    pipeline = [
        {
            '$match': {
                'product': product.id
            }
        },
        {
            '$sort': {
                'date_': -1
            }
        },
        {
            '$unwind': '$visits'
        },
        {
            '$addFields': {
                'formated_date': {
                    '$dateFromParts': {
                        'year': {'$year': '$date_'},
                        'month': {'$month': '$date_'},
                        'day': {'$dayOfMonth': '$date_'},
                        'hour': {'$toInt': {'$arrayElemAt': [{'$split': ['$visits.time', ':']}, 0]}},
                        'minute': {'$toInt': {'$arrayElemAt': [{'$split': ['$visits.time', ':']}, 1]}},
                        'second': {'$toInt': {'$arrayElemAt': [{'$split': ['$visits.time', ':']}, 2]}}
                    }
                }
            }
        },
        {
            '$sort': {
                'formated_date': -1
            }
        },
        {'$facet': {
            'totalCount': [
                {
                    '$count': 'count'
                }
            ],
            'paginatedResults': [
                {
                    '$skip': skip_count
                },
                {
                    '$limit': page_size
                },
                {
                    '$project': {
                        'visits': 1,
                        'city': 1,
                        'region': 1,
                        'country': 1,
                        'timezone': 1,
                        'postal': 1,
                        'ip': 1,
                        'isp': 1,
                        'org': 1,
                        'as_': 1,
                        'hosting': 1,
                        'proxy': 1,
                        'mobile': 1,
                        'user_agent': 1,
                        'formated_date': 1
                    }
                }
            ]
        }, }
    ]
    results = await connection.visits.aggregate(pipeline).to_list(None)
    total_documents = results[0]['totalCount'][0]['count'] if results[0]['totalCount'] else 0
    total_pages = math.ceil(total_documents / page_size)

    paginated_results = results[0]['paginatedResults']
    latest = [
        RequestData(
            ip=x.get('ip'),
            time=x.get('formated_date'),
            path=x.get('visits').get('path'),
            method=x.get('visits').get('method'),
            timezone=x.get('timezone'),
            isp=f"{x.get('isp')} | {x.get('org')} | {x.get('as_')}",
            postal=x.get('postal'),
            location=f"{x.get('city')}, {x.get('region')}, {x.get('country')}",
            user_agent=x.get('user_agent')
        )
        for x in paginated_results
    ]
    return CustomResponse.get_success_response(f"{code} requests",
                                               data=ProductRequestsResponse(
                                                   requests=latest,
                                                   page=page,
                                                   total_pages=total_pages,
                                                   total=total_documents,
                                                   page_size=page_size
                                               ))


@router.get('/{code}/locations', description="Get information about a project")
async def product_locations(code: str):
    connection = Connection()
    product = await connection.products.find_one({'code': code})
    if not product:
        return CustomResponse.get_failure_response("No product with the code found!")
    product = Products(**product)
    pipeline = [
        {
            '$match': {
                'product': product.id
            }
        },
        {
            '$unwind': '$visits'
        },
        {
            '$facet': {
                'top_cities': [
                    {
                        '$group': {
                            '_id': '$city',
                            'total_ci_sum': {
                                '$sum': 1
                            },
                            'region': {
                                '$first': '$region'
                            },
                            'country': {
                                '$first': '$country'
                            }
                        }
                    },
                    {
                        '$sort': {
                            'total_ci_sum': -1
                        }
                    },
                    {
                        '$limit': 10
                    },
                    {
                        '$project': {
                            '_id': 0,
                            'city': '$_id',
                            'total_ci_sum': 1,
                            'region': 1,
                            'country': 1
                        }
                    },
                ],
                'top_regions': [
                    {
                        '$group': {
                            '_id': '$region',
                            'total_r_sum': {
                                '$sum': 1
                            },
                            'country': {
                                '$first': '$country'
                            }
                        }
                    },
                    {
                        '$sort': {
                            'total_r_sum': -1
                        }
                    },
                    {
                        '$limit': 10
                    },
                    {
                        '$project': {
                            '_id': 0,
                            'region': '$_id',
                            'total_r_sum': 1,
                            'country': 1
                        }
                    },
                ],
                'top_countries': [
                    {
                        '$group': {
                            '_id': '$country',
                            'total_co_sum': {
                                '$sum': 1
                            }
                        }
                    },
                    {
                        '$sort': {
                            'total_co_sum': -1
                        }
                    },
                    {
                        '$limit': 10
                    },
                    {
                        '$project': {
                            '_id': 0,
                            'country': '$_id',
                            'total_co_sum': 1
                        }
                    },
                ]
            }
        },
    ]

    results = await connection.visits.aggregate(pipeline).to_list(None)
    top_cities = results[0].get('top_cities', [])
    top_cities = [LocationData(
        location=f"{x.get('city')}, {x.get('region')}, {x.get('country')}",
        count=x.get('total_ci_sum', 0)
    ) for x in top_cities]
    top_regions = results[0].get('top_regions', [])
    top_regions = [LocationData(
        location=f"{x.get('region')}, {x.get('country')}",
        count=x.get('total_r_sum', 0)
    ) for x in top_regions]
    top_countries = results[0].get('top_countries', [])
    top_countries = [LocationData(
        location=f"{x.get('country')}",
        count=x.get('total_co_sum', 0)
    ) for x in top_countries]
    return CustomResponse.get_success_response(f"{code} locations",
                                               data=ProductLocationsResponse(
                                                   top_cities=top_cities,
                                                   top_countries=top_countries,
                                                   top_regions=top_regions
                                               ))
