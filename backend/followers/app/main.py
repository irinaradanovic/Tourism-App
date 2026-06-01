from fastapi import FastAPI
from api.endpoints import router as followers_router
from fastapi.openapi.utils import get_openapi
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(title="Followers Service")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.include_router(followers_router, prefix="/api/followers")

@app.get("/")
async def root():
    return {"message": "Welcome to the Followers Service!"}
