from fastapi import APIRouter, HTTPException
from model.schemas import FollowRequest, FollowerInfo, RecommendedFollower
from service.follower_service import FollowerService
from typing import List
from utils.auth import get_current_user
from fastapi import Depends

router = APIRouter()
service = FollowerService()

@router.post("/follow")
async def follow_user(
    followedId: int, 
    current_user_id: int = Depends(get_current_user)
):
    # current_user_id we get from the token
    data = FollowRequest(followerId=current_user_id, followedId=followedId)
    try:
        return service.process_follow(data)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))

@router.delete("/unfollow/{followedId}")
async def unfollow_user(
    followedId: int,
    current_user_id: int = Depends(get_current_user)
):
    data = FollowRequest(followerId=current_user_id, followedId=followedId)
    try:
        return service.process_unfollow(data)
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    
@router.get("/recommendations")
async def get_recommendations(
    current_user_id: int = Depends(get_current_user)
):
    return await service.get_recommendations(current_user_id)

@router.get("/my-followings", response_model=List[FollowerInfo])
async def get_following_with_data(current_user_id: int = Depends(get_current_user)):
    return await service.get_following_with_data(current_user_id)

@router.get("/my-followers", response_model=List[FollowerInfo])
async def get_followers_with_data(current_user_id: int = Depends(get_current_user)):
    return await service.get_followers_with_data(current_user_id)

@router.get("/{userId}/followings", response_model=List[FollowerInfo])
async def get_user_following(userId: int):
    return await service.get_following_with_data(userId)

@router.get("/{userId}/followers", response_model=List[FollowerInfo])
async def get_user_followers(userId: int):
    return await service.get_followers_with_data(userId)