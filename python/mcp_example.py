"""
CarbonArc MCP Client — Python Quick Start
==========================================

Connect to the CarbonArc MCP server using fastmcp,
authenticate via OAuth 2.0, and discover available tools.

Run:
    uv run python mcp_example.py
"""

import asyncio

from fastmcp import Client
from fastmcp.client.auth import OAuth


MCP_SERVER_URL = "https://mcp.carbonarc.ai/"


async def main():
    # Connect to CarbonArc with OAuth (browser opens automatically for login)
    oauth = OAuth(MCP_SERVER_URL, client_name="CarbonArc MCP Client Example")

    async with Client(MCP_SERVER_URL, auth=oauth) as client:
        # Discover tools
        tools = await client.list_tools()

        print(f"\nDiscovered {len(tools)} tool(s):\n")
        for tool in tools:
            desc = (tool.description or "").split("\n")[0]
            print(f"  • {tool.name} — {desc}")

        # Test: search for "Walmart"
        print('\nSearching for Walmart\n')
        result = await client.call_tool("search_entities", {"query": "Walmart"})

        # `.content` is the wire format every MCP SDK returns — a list of
        # content blocks. Use it when you want the human-readable rendering.
        for block in result.content:
            if block.type == "text":
                print(block.text)

        # `.data` is the same response deserialized into Python objects, and is
        # None for tools that declare no structured output. Prefer it in code.
        if result.data is not None:
            print(f"\nStructured: {result.data!r}")


if __name__ == "__main__":
    asyncio.run(main())