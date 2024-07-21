import motor.motor_asyncio, motor.core
import pymongo
from bson import SON
from decouple import config

DB_URL = config("DB_URL")
DB_NAME = config("DB_NAME")


class Connection:
    client: motor.motor_asyncio.AsyncIOMotorClient
    db: motor.motor_asyncio.AsyncIOMotorDatabase
    products: motor.motor_asyncio.AsyncIOMotorCollection
    visits: motor.motor_asyncio.AsyncIOMotorCollection
    users: motor.motor_asyncio.AsyncIOMotorCollection

    def __init__(self):
        self.client = motor.motor_asyncio.AsyncIOMotorClient(DB_URL)
        self.db = self.client.get_database(name=DB_NAME)
        self.products = self.db.get_collection('products')
        self.users = self.db.get_collection('users')
        self.visits = self.db.get_collection('visits')

    @staticmethod
    async def unique_index(collection: motor.motor_asyncio.AsyncIOMotorCollection, key: str):
        indexes: list[SON] = await collection.list_indexes().to_list(length=None)
        index_exists = False
        for index in indexes:
            if index.get('key').get(key) and index.get('unique'):
                index_exists = True
        if not index_exists:
            await collection.create_index([(key, pymongo.ASCENDING)], unique=True)

    async def initialize_database(self):
        await self.unique_index(self.products, 'name')
        await self.unique_index(self.products, 'code')
        await self.unique_index(self.visits, 'id')
        await self.unique_index(self.users, 'email')

    def __getitem__(self, collection_name: str) -> motor.motor_asyncio.AsyncIOMotorCollection | None:
        return self.db.get_collection(collection_name)
