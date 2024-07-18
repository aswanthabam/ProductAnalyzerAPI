from app import app
from a2wsgi import ASGIMiddleware

# Convert fastapi asgi application into wsgi for deploying it in lambda serverless functions.
app = ASGIMiddleware(app)