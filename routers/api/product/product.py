import logging
from datetime import datetime

import decouple
import pymongo.errors
import requests
from fastapi import APIRouter, Form
from starlette.requests import Request

from db.connection import Connection
from db.models import Users, ProductVisits, Products, Countries, Regions, Cities
from routers.api.product.response_models import ProductResponse, ListProductsResponse, ProductInfoResponse, RequestData
from utils.response import CustomResponse

router = APIRouter(
    prefix="/api/product",
    tags=["product"],
    responses={404: {"description": "Not found"}},
)
GENERAL_PRODUCT_LIMIT_PER_HOUR = 60
PASSWORD = decouple.config('PASSWORD')


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


@router.get('/{code}/visit/', description="Visit a product")
async def visit(code: str, request: Request):
    connection = Connection()
    # finds the ip address of the user
    x_forwarded_for = request.headers.get("x-forwarded-for")
    if x_forwarded_for:
        client_ip = x_forwarded_for.split(",")[0]
    else:
        client_ip = request.client.host
    if not client_ip:
        return CustomResponse.get_failure_response("Error, trust me.. its not from our side.")
    # find if the ip stored already, if not store it
    user = await connection.users.find_one({'ip': client_ip})
    if not user:
        await connection.users.insert_one(Users(ip=client_ip).model_dump())
        user = await connection.users.find_one({'ip': client_ip})
    if not user:
        return CustomResponse.get_failure_response("Error, its from our side :)")
    user = Users(**user)
    # find the product
    product = await connection.products.find_one({'code': code})
    if not product:
        return CustomResponse.get_failure_response("Product not found.")
    product = Products(**product)
    # find the details about the ip address
    ipapi = f"http://ip-api.com/json/{client_ip}?fields=status,message,continent,continentCode,country,countryCode,region,regionName,city,zip,lat,lon,timezone,offset,currency,isp,org,as,asname,reverse,mobile,proxy,hosting,query"
    res = requests.get(ipapi)
    data = res.json()
    if not data.get('error') and data.get('status') == "success":
        try:
            city_name = data.get('city')
            region_name = data.get('regionName')
            region_code = data.get('region')
            country_name = data.get('country')
            country_code = data.get('countryCode')
            continent = data.get('continent')
            postal = data.get('zip')
            timezone = data.get('timezone')
            isp = data.get('isp')
            latitude = data.get('lat')
            longitude = data.get('lon')
            # Find / Create the country, region and city
            country = await connection.countries.find_one(
                {'code': country_code, 'continent_code': continent, 'name': country_name})
            if not country:
                country = Countries(name=country_name, code=country_code, continent_code=continent)
                await connection.countries.insert_one(country.model_dump())
                country = await connection.countries.find_one(
                    {'code': country_code, 'continent_code': continent, 'name': country_name})
                if not country:
                    return CustomResponse.get_failure_response("Country not found!!")
            country = Countries(**country)
            region = await connection.regions.find_one(
                {'code': region_code, 'country': country.id, 'name': region_name})
            if not region:
                region = Regions(name=region_name, code=region_code, country=country.id)
                await connection.regions.insert_one(region.model_dump())
                region = await connection.regions.find_one(
                    {'code': region_code, 'country': country.id, 'name': region_name})
                if not region:
                    return CustomResponse.get_failure_response("Region not found!!")
            region = Regions(**region)
            city = await connection.cities.find_one({'region': str(region.id), 'name': city_name, 'timezone': timezone})
            if not city:
                city = Cities(region=region.id, name=city_name, timezone=timezone)
                await connection.cities.insert_one(city.model_dump())
                city = await connection.cities.find_one(
                    {'region': str(region.id), 'name': city_name, 'timezone': timezone})
                if not city:
                    return CustomResponse.get_failure_response("City not found!!")
            city = Cities(**city)
        except Exception as e:
            logging.getLogger(__name__).error(e, stack_info=True)
            city = None
            postal = None
            isp = None
            latitude = None
            longitude = None
    else:
        city = None
        postal = None
        isp = None
        latitude = None
        longitude = None
    # Add product visit
    product_visit = ProductVisits(
        product=product.id,
        user=user.id,
        city=city.id if city else None,
        latitude=str(latitude) if latitude else None,
        longitude=str(longitude) if longitude else None,
        postal=str(postal) if postal else None,
        org=isp,
    )
    await connection.product_visits.insert_one(product_visit.model_dump())
    # Update the user visit count
    await connection.users.update_one({'id': user.id}, {'$set': {
        'visit_count': user.visit_count + 1,
        'last_visit': datetime.utcnow()
    }})
    return CustomResponse.get_success_response("OK")


@router.get('/list', description="List available products")
async def list():
    connection = Connection()
    products = await connection.products.find().to_list(length=None)
    ps = []
    now = datetime.utcnow()
    start_of_month = datetime(now.year, now.month, 1)
    if now.month == 12:
        start_of_next_month = datetime(now.year + 1, 1, 1)
    else:
        start_of_next_month = datetime(now.year, now.month + 1, 1)
    for product in products:
        product = Products(**product)
        visit_count = await connection.product_visits.count_documents({'product': product.id})
        monthly_visit = await connection.product_visits.count_documents({
            'time': {'$gte': start_of_month, '$lt': start_of_next_month},
            'product': product.id
        })
        ps.append(
            ProductResponse(name=product.name, code=product.code, created_at=product.created_at,
                            total_visits=visit_count, monthly_visits=monthly_visit))
    return CustomResponse.get_success_response("", data=ListProductsResponse(products=ps))


@router.get('/{code}/info', description="Get information about a project")
async def product_info(code: str):
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
            '$lookup': {
                'from': 'cities',
                'localField': 'city',
                'foreignField': 'id',
                'as': 'city_info'
            }
        },
        {
            '$unwind': '$city_info'
        },
        {
            '$lookup': {
                'from': 'regions',
                'localField': 'city_info.region',
                'foreignField': 'id',
                'as': 'region_info'
            }
        },
        {
            '$unwind': '$region_info'
        },
        {
            '$lookup': {
                'from': 'countries',
                'localField': 'region_info.country',
                'foreignField': 'id',
                'as': 'country_info'
            }
        },
        {
            '$unwind': '$country_info'
        },
        {
            '$group': {
                '_id': {
                    'city': '$city_info.id',
                    'city_name': '$city_info.name',
                    'region': '$region_info.id',
                    'region_name': '$region_info.name',
                    'country': '$country_info.id',
                    'country_name': '$country_info.name'
                },
                'count': {'$sum': 1}
            }
        },
        {
            '$sort': {'count': -1}
        },
        {
            '$project': {
                '_id': 0,
                'city_id': '$_id.city',
                'city_name': '$_id.city_name',
                'region_id': '$_id.region',
                'region_name': '$_id.region_name',
                'country_id': '$_id.country',
                'country_name': '$_id.country_name',
                'count': 1
            }
        }
    ]

    results = connection.product_visits.aggregate(pipeline)
    results = await results.to_list(length=None)

    cities = {}
    regions = {}
    countries = {}
    for row in results:
        city_id = row.get('city_id')
        region_id = row.get('region_id')
        country_id = row.get("country_id")
        if not city_id or not region_id or not country_id:
            if cities.get("unknown"):
                cities['unknown']['count'] += row.get('count')
            else:
                cities['unknown'] = {
                    'name': 'Unknown',
                    'count': row.get('count')
                }
            continue
        cities[city_id] = {
            "name": row.get("city_name"),
            "count": row.get("count")
        }
        if not regions.get(region_id):
            regions[region_id] = {
                "name": row.get("region_name"),
                "count": row.get("count")
            }
        else:
            regions[region_id]['count'] += row.get('count')
        if not countries.get(country_id):
            countries[country_id] = {
                "name": row.get("country_name"),
                "count": row.get("count")
            }
        else:
            countries[country_id]['count'] += row.get('count')
    cities = sorted(cities.values(), key=lambda x: x['count'], reverse=True)[:10]
    regions = sorted(regions.values(), key=lambda x: x['count'], reverse=True)[:10]
    countries = sorted(countries.values(), key=lambda x: x['count'], reverse=True)[:10]
    pipeline = [
        {'$match': {'product': product.id}},
        {'$sort': {'time': -1}},
        {'$limit': 5},
        {'$lookup': {
            'from': 'cities',
            'localField': 'city',
            'foreignField': 'id',
            'as': 'city_info'
        }},
        {'$unwind': '$city_info'},
        {'$lookup': {
            'from': 'regions',
            'localField': 'city_info.region',
            'foreignField': 'id',
            'as': 'region_info'
        }},
        {'$unwind': '$region_info'},
        {'$lookup': {
            'from': 'countries',
            'localField': 'region_info.country',
            'foreignField': 'id',
            'as': 'country_info'
        }},
        {'$unwind': '$country_info'},
        {'$lookup': {
            'from': 'users',
            'localField': 'user',
            'foreignField': 'id',
            'as': 'user_info'
        }},
        {'$unwind': '$user_info'},
        {'$project': {

            'time': 1,
            'isp': 1,
            'postal': 1,
            'ip': '$user_info.ip',
            'timezone': '$city_info.timezone',
            'city': '$city_info.name',
            'region': '$region_info.name',
            'country': '$country_info.name',
            'continent': '$country_info.continent'
        }}
    ]
    results = await connection.product_visits.aggregate(pipeline).to_list(None)

    latest = [
        RequestData(
            ip=x.get('ip'),
            time=x.get('time'),
            isp=x.get('isp'),
            postal=x.get('postal'),
            timezone=x.get('timezone'),
            city=x.get('city'),
            region=x.get('region'),
            country=x.get('country'),
            continent=x.get('continent')
        )
        for x in results
    ]
    return CustomResponse.get_success_response(f"{code} info",
                                               data=ProductInfoResponse(
                                                   name=product.name,
                                                   latest_visits=latest,
                                                   cities=cities, countries=countries,
                                                   regions=regions))
