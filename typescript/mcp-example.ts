/**
 * CarbonArc MCP Client — TypeScript Quick Start
 * ==============================================
 *
 * Connect to the CarbonArc MCP server using the official
 * @modelcontextprotocol/sdk, authenticate via OAuth 2.0, and discover
 * available tools.
 *
 * Setup:
 *     npm install
 *
 * Run:
 *     npx tsx mcp-example.ts
 */

import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StreamableHTTPClientTransport } from "@modelcontextprotocol/sdk/client/streamableHttp.js";
import { UnauthorizedError } from "@modelcontextprotocol/sdk/client/auth.js";
import { BrowserOAuthProvider, waitForAuthCode } from "./oauth-helper.js";

const MCP_SERVER_URL = "https://mcp.carbonarc.ai/";
const CALLBACK_PORT = 8090;

async function main() {
  const authProvider = new BrowserOAuthProvider(
    `http://localhost:${CALLBACK_PORT}/callback`
  );
  const serverUrl = new URL(MCP_SERVER_URL);

  const client = new Client(
    { name: "carbonarc-mcp-client-ts", version: "1.0.0" },
    { capabilities: {} }
  );

  // Connect — handle the OAuth redirect if needed
  let transport = new StreamableHTTPClientTransport(serverUrl, { authProvider });
  try {
    await client.connect(transport);
  } catch (err) {
    if (err instanceof UnauthorizedError) {
      console.log("OAuth authorization required…");
      const code = await waitForAuthCode(CALLBACK_PORT);
      await transport.finishAuth(code);
      transport = new StreamableHTTPClientTransport(serverUrl, { authProvider });
      await client.connect(transport);
    } else {
      throw err;
    }
  }

  // Discover tools
  const { tools } = await client.listTools();
  console.log(`\nDiscovered ${tools.length} tool(s):\n`);
  for (const tool of tools) {
    const desc = tool.description?.split("\n")[0] ?? "";
    console.log(`  • ${tool.name} — ${desc}`);
  }

  // Test: search for "Walmart"
  console.log("\nSearching for Walmart\n");
  const result = await client.callTool({
    name: "search_entities",
    arguments: { query: "Walmart" },
  });
  for (const item of result.content as Array<{ type: string; text?: string }>) {
    if (item.type === "text") console.log(item.text);
  }

  await client.close();
}

main().catch(console.error);
