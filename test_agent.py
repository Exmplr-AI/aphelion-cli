#!/usr/bin/env python3
"""
Quick test of the agent to verify API integration
"""
import os
import sys
from agent import AphelionAgent

# Set the token
os.environ["APHELION_TOKEN"] = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6InlsRGFzRFNHMWtlLUhyQ0NmcVpaYyJ9.eyJodHRwczovL2FwaGVsaW9uLWdhdGV3YXkuY29tL3JvbGVzIjpbImFkbWluIiwidXNlciJdLCJpc3MiOiJodHRwczovL2F1dGguYXBoZWxpb24uZXhtcGxyLmFpLyIsInN1YiI6Imdvb2dsZS1vYXV0aDJ8MTA3MTk2NjYwOTk3MjU0MTgyMjM3IiwiYXVkIjpbImh0dHBzOi8vYXBoZWxpb24tZ2F0ZXdheS5jb20iLCJodHRwczovL2Rldi1heTB3NmgycnJlY3NvcHQ4LnVzLmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE3NTAwOTgxODEsImV4cCI6MTc1MDE4NDU4MSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsImF6cCI6IlViWHhwUUJTcjlBc3FwUzJsbjJqem1zbWFtcm9tYUZDIiwicGVybWlzc2lvbnMiOltdfQ.U3e2IUq_mMRWFvnsB3pbNvfhVYftfLwRQ6RBZByWyySdCVRapFmfF5ttPFOxiw5SeDsB3-hdwCaV20YfY-ROHBkoLczqrLUJa06Dpjy2K_xKmpIaZ-U7iy7SCO8pZfVMeTM9iDuntCiSE38AzACYzP4W5oqvGCHChKUPkjs3k3phSwQxzLq4PLvMo0TJUEhLDYl2lFJ8GFCYouURJa-xZ2AfLiMVwvg-eAw22J_3lf0BDquwZFS6e0W4pMQaGdq8-_ILr4YMZ-f5UlaYaQR1z6GrlM7sQd7cDpBJjA5kDT56Mw6HwILeolui2GL73UD7xGlrnR7PXI61p6LWnP3jGA"

def test_agent():
    print("🧪 Testing Aphelion Agent API Integration...")
    
    try:
        # Initialize agent
        agent = AphelionAgent()
        print(f"✅ Agent initialized with session: {agent.session_id}")
        print(f"✅ API URL: {agent.config['gateway']['api_url']}")
        
        # Test 1: Search tools
        print("\n1. Testing search_tools()...")
        tools_result = agent.search_tools("Multiple Sclerosis")
        print(f"✅ Search result: {tools_result}")
        
        # Test 2: Run tool (if tools are available)
        if tools_result.get("tools"):
            print("\n2. Testing run_tool()...")
            # Use first tool found
            first_tool = tools_result["tools"][0]
            tool_name = first_tool.get("name", "echo")
            # Use appropriate parameters based on tool schema
            if "weather" in tool_name.lower():
                params = {"city": "San Francisco"}
            else:
                params = {"message": "Test from agent"}
            result = agent.run_tool(tool_name, params)
            print(f"✅ Tool execution result: {result}")
        else:
            print("\n2. No tools found, testing with example tool...")
            result = agent.run_tool("echo", {"message": "Test from agent"})
            print(f"Tool execution result: {result}")
        
        # Test 3: Save memory
        print("\n3. Testing save_memory()...")
        agent.save_memory(
            "Test memory save", 
            {"test": True, "timestamp": "2025-01-16"}
        )
        print("✅ Memory save completed")
        
        print("\n🎉 All API integrations working!")
        return True
        
    except Exception as e:
        print(f"❌ Error during testing: {e}")
        import traceback
        traceback.print_exc()
        return False

if __name__ == "__main__":
    success = test_agent()
    sys.exit(0 if success else 1)