import jwt
import os
from fastapi import HTTPException, Security
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from jwt.exceptions import ExpiredSignatureError, InvalidTokenError

security = HTTPBearer()

JWT_SECRET = os.getenv("JWT_SECRET")

async def get_current_user(auth: HTTPAuthorizationCredentials = Security(security)):
    if not auth:
        raise HTTPException(status_code=403, detail="Not authenticated")
    
    token = auth.credentials 
    try:
        secret_bytes = JWT_SECRET.encode('utf-8')

        payload = jwt.decode(token, secret_bytes, algorithms=["HS256"])
        
        user_id = payload.get("sub")
        if user_id is None:
            raise HTTPException(status_code=403, detail="Token missing subject (sub)")
            
        return int(user_id)
        
    except ExpiredSignatureError:
        raise HTTPException(status_code=403, detail="Token has expired")
    except InvalidTokenError as e:
        print(f"JWT Error: {str(e)}")
        raise HTTPException(status_code=403, detail=f"Invalid token: {str(e)}")
    except Exception as e:
        print(f"Unexpected error: {str(e)}")
        raise HTTPException(status_code=500, detail="Auth error")