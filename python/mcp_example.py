"""
CarbonArc MCP Client — Python Quick Start
==========================================

Connect to the CarbonArc MCP server using langchain-mcp-adapters,
authenticate via OAuth 2.0, and discover available tools.

Run:
    uv run python mcp_example.py
"""

import asyncio

from fastmcp import Client
from fastmcp.client.auth import OAuth
from langchain_mcp_adapters.tools import load_mcp_tools


MCP_SERVER_URL = "https://mcp.carbonarc.co/"


async def main():
    # Connect to CarbonArc with OAuth (browser opens automatically for login)
    oauth = OAuth(MCP_SERVER_URL, client_name="CarbonArc MCP Client Example")

    async with Client(MCP_SERVER_URL, auth=oauth) as client:
        # Load tools as LangChain-compatible tools
        tools = await load_mcp_tools(client.session)

        print(f"\nDiscovered {len(tools)} tool(s):\n")
        for tool in tools:
            desc = tool.description.split("\n")[0]
            print(f"  • {tool.name} — {desc}")

        # Test: search for "Walmart"
        search_tool = next(t for t in tools if t.name == "search_entities")
        print('\nSearching for Walmart\n')
        result = await search_tool.ainvoke({"query": "Walmart"})
        print(result)


if __name__ == "__main__":
    asyncio.run(main())