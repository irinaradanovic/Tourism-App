from repository.follower_repo import FollowerRepo
from model.schemas import FollowRequest, FollowerInfo, RecommendedFollower
import httpx
import os
from typing import List

class FollowerService:
    def __init__(self):
        self.repo = FollowerRepo()
        self.stakeholders_url = os.getenv("STAKEHOLDERS_SERVICE_URL")

    def process_follow(self, data: FollowRequest):
        if data.followerId == data.followedId:
            raise ValueError("You cannot follow yourself.")
        
        self.repo.follow_user(data.followerId, data.followedId)
        return {"message": "Successfully followed"}

    def process_unfollow(self, data: FollowRequest):
        if data.followerId == data.followedId:
            raise ValueError("You cannot unfollow yourself.")
        
        self.repo.unfollow_user(data.followerId, data.followedId)
        return {"message": "Successfully unfollowed"}
    
    async def get_recommendations(self, userId: int):
        raw_recommendations = self.repo.get_recommendations(userId)
        
        if not raw_recommendations:
            return []

        ids = [item['userId'] for item in raw_recommendations]

        users_data = await self._populate_usernames(ids)
 
        recommended_with_usernames = []
        for i, user_info in enumerate(users_data):
            recommended_with_usernames.append(RecommendedFollower(
                userId=user_info.userId,
                username=user_info.username,
                mutualFollowersCount=raw_recommendations[i]['mutualFollowersCount']
            ))
        
        return recommended_with_usernames

    async def get_following_with_data(self, userId: int):
        ids = self.repo.get_following_ids(userId)
        return await self._populate_usernames(ids)

    async def get_followers_with_data(self, userId: int):
        ids = self.repo.get_followers_ids(userId)
        return await self._populate_usernames(ids)
    

    #util
    async def _populate_usernames(self, user_ids: List[int]) -> List[FollowerInfo]:
        populated_list = []
        async with httpx.AsyncClient() as client:
            for uid in user_ids:
                try:
                    response = await client.get(f"{self.stakeholders_url}/{uid}")
                    if response.status_code == 200:
                        data = response.json()
                        populated_list.append(FollowerInfo(
                            userId=uid,
                            username=data.get('username', 'Unknown')
                        ))
                    else:
                        populated_list.append(FollowerInfo(userId=str(uid), username="Unknown"))
                except Exception as e:
                    print(f"Error calling Stakeholders: {e}")
                    populated_list.append(FollowerInfo(userId=str(uid), username="Service Unavailable"))
        return populated_list