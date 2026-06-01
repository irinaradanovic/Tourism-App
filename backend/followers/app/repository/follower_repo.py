import os
from neo4j import GraphDatabase

class FollowerRepo:
    def __init__(self):
        uri = os.getenv("NEO4J_URI")
        user = os.getenv("NEO4J_USER")
        password = os.getenv("NEO4J_PASSWORD")
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def follow_user(self, fId, tId):
        with self.driver.session() as session:
            query = """
            MERGE (a:User {userId: $fId})
            MERGE (b:User {userId: $tId})
            MERGE (a)-[:FOLLOWS]->(b)
            """
            session.run(query, fId=int(fId), tId=int(tId))

    def unfollow_user(self, fId, tId):
        with self.driver.session() as session:
            query = """
            MATCH (a:User {userId: $fId})-[r:FOLLOWS]->(b:User {userId: $tId})
            DETACH DELETE r
            """
            session.run(query, fId=int(fId), tId=int(tId))

    def get_following_ids(self, uId):
        with self.driver.session() as session:
            query = """
            MATCH (a:User {userId: $uId})-[:FOLLOWS]->(b)
            RETURN DISTINCT b.userId as id
            """
            result = session.run(query, uId=int(uId))
            return [record["id"] for record in result]


    def get_followers_ids(self, uId):
        with self.driver.session() as session:
            query = """
            MATCH (a)-[:FOLLOWS]->(b:User {userId: $uId})
            RETURN a.userId as id
            """
            result = session.run(query, uId=int(uId))
            return [record["id"] for record in result]
        
    def get_recommendations(self, uId):
        with self.driver.session() as session:
            query = """
            MATCH (me:User {userId: $uId})-[:FOLLOWS]->(friend:User)
            MATCH (friend)-[:FOLLOWS]->(foaf:User)
            WITH me, foaf, count(friend) as mutualCount
            WHERE NOT (me)-[:FOLLOWS]->(foaf) AND me.userId <> foaf.userId
            RETURN foaf.userId as userId, mutualCount as mutualFollowersCount
            ORDER BY mutualFollowersCount DESC
            LIMIT 5
            """
            result = session.run(query, uId=int(uId))
            return [dict(record) for record in result]