from pydantic import BaseModel
from typing import Optional, List

class FollowRequest(BaseModel):
    followerId: int
    followedId: int

class FollowerInfo(BaseModel):
    userId: int
    username: Optional[str] = "Unknown"

class RecommendedFollower(BaseModel):
    userId: int
    username: Optional[str] = "Unknown"
    mutualFollowersCount: int 