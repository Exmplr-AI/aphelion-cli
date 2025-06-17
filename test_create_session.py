#!/usr/bin/env python3
"""
Test session creation directly
"""
import os
import requests

# Set the token
os.environ["APHELION_TOKEN"] = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InlsRGFzRFNHMWtlLUhyQ0NmcVpaYyJ9.eyJodHRwczovL2FwaGVsaW9uLWdhdGV3YXkuY29tL3JvbGVzIjpbImFkbWluIiwidXNlciJdLCJpc3MiOiJodHRwczovL2F1dGguYXBoZWxpb24uZXhtcGxyLmFpLyIsInN1YiI6Imdvb2dsZS1vYXV0aDJ8MTA3MTk2NjYwOTk3MjU0MTgyMjM3IiwiYXVkIjpbImh0dHBzOi8vYXBoZWxpb24tZ2F0ZXdheS5jb20iLCJodHRwczovL2Rldi1heTB3NmgycnJlY3NvcHQ4LnVzLmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3NTAwOTgxODEsImV4cCI6MTc1MDE4NDU4MSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImF6cCI6IlViWHhwUUJTcjlBc3FwUzJsbjJqem1zbWFtcm9tYUZDIiwicGVybWlzc2lvbnMiOltdfQ.U3e2IUq_mMRWFvnsB3pbNvfhVYftfLwRQ6RBZByWyySdCVRapFmfF5ttPFOxiw5SeDsB3-hdwCaV20YfY-ROHBkoLczqrLUJa06Dpjy2K_xKmpIaZ-U7iy7SCO8pZfVMeTM9iDuntCiSE38AzACYzP4W5oqvGCHChKUPkjs3k3phSwQxzLq4PLvMo0TJUEhLDYl2lFJ8GFCYouURJa-xZ2AfLiMVwvg-eAw22J_3lf0BDquwZFS6e0W4pMQaGdq8-_ILr4YMZ-f5UlaYaQR1z6GrlM7sQd7cDpBJjA5kDT56Mw6HwILeolui2GL73UD7xGlrnR7PXI61p6LWnP3jGA"

def test_session_creation():
    print("🧪 Testing Session Creation...")
    
    token = os.getenv("APHELION_TOKEN")
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json"
    }
    url = "https://api.aphelion.exmplr.ai/v1/agents"
    payload = {
        "subscribed_services": []
    }
    
    try:
        print("Creating new session...")
        response = requests.post(url, headers=headers, json=payload)
        print(f"Status: {response.status_code}")
        print(f"Response: {response.text}")
        
        if response.status_code == 200 or response.status_code == 201:
            result = response.json()
            session_id = result.get("session_id")
            print(f"✅ New session created: {session_id}")
            
            # Test tool execution with real session
            if session_id:
                print(f"\n🧪 Testing tool execution with session {session_id}...")
                exec_url = f"https://api.aphelion.exmplr.ai/v1/agents/{session_id}/execute"
                exec_payload = {
                    "tool": "echo",
                    "parameters": {"message": "Hello from real session!"}
                }
                
                exec_response = requests.post(exec_url, headers=headers, json=exec_payload)
                print(f"Tool execution status: {exec_response.status_code}")
                print(f"Tool execution response: {exec_response.text}")
                
        else:
            print(f"❌ Failed to create session: {response.status_code}")
            print(response.text)
            
    except Exception as e:
        print(f"❌ Error: {e}")

if __name__ == "__main__":
    test_session_creation()