from datetime import datetime

import decouple
import pymongo.errors
import requests
from fastapi import APIRouter, Form
from starlette.requests import Request

from db.connection import Connection
from db.models import Users, ProductVisits, Products, Countries, Regions, Cities
from utils.response import CustomResponse

router = APIRouter(
    prefix="/product",
    tags=["product"],
    responses={404: {"description": "Not found"}},
)
GENERAL_PRODUCT_LIMIT_PER_HOUR = 60
PASSWORD = decouple.config('PASSWORD')


@router.post('/create-product/', description="Create a product")
async def create_product(code: str = Form(...), name: str = Form(...), password: str = Form(...)):
    if password != PASSWORD:
        return CustomResponse.get_failure_response("Unauthorized!")
    connection = Connection()
    try:
        await connection.products.insert_one(Products(code=code, name=name).model_dump())
    except pymongo.errors.DuplicateKeyError as e:
        return CustomResponse.get_failure_response("Product already exists.")
    return CustomResponse.get_success_response("Product Created")


@router.get('/visit/{code}', description="Visit a product")
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
    if not data.get('error'):
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
                city = await connection.cities.find_one({'region': str(region.id), 'name': city_name, 'timezone': timezone})
                if not city:
                    return CustomResponse.get_failure_response("City not found!!")
            city = Cities(**city)
        except Exception as e:
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
    await connection.users.update_one({'_id': user.id}, {'$set': {
        'visit_count': user.visit_count + 1,
        'last_visit': datetime.utcnow()
    }})
    return CustomResponse.get_success_response("OK")
