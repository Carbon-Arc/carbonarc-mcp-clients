# CarbonArc MCP Client Examples

Minimal, self-contained MCP client examples showing how to connect to the
[CarbonArc](http://docs.carbonarc.ai/) MCP server from every popular language
using the **official Model Context Protocol SDKs**.

Each example follows the same pattern:

1. Connect to `https://mcp.carbonarc.ai/` over Streamable HTTP
2. Authenticate via OAuth 2.0 (browser opens automatically where supported)
3. Discover available tools
4. Call `search_entities` with `"Walmart"` as a quick smoke test

---

## Client Examples at a Glance

| Language | Directory | SDK | Transport |
|---|---|---|---|
| **Python** | [`python/`](python/) | `fastmcp` + `langchain-mcp-adapters` | Streamable HTTP |
| **TypeScript / Node.js** | [`typescript/`](typescript/) | `@modelcontextprotocol/sdk` | `StreamableHTTPClientTransport` |
| **Go** | [`go/`](go/) | `modelcontextprotocol/go-sdk` | `StreamableClientTransport` |

---

## Python

```bash
cd python
uv run python mcp_example.py
```

## TypeScript / Node.js

Uses the [official MCP TypeScript SDK](https://github.com/modelcontextprotocol/typescript-sdk)
(`@modelcontextprotocol/sdk`). Includes a full OAuth 2.0 authorization-code
flow with a local callback server.

```bash
cd typescript
npm install
npx tsx mcp-example.ts
```

## Go

Uses the [official MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
(`v1.5.0+`). Requires **Go 1.25+**.

```bash
cd go
go mod tidy
go run -tags mcp_go_client_oauth .
```

---

## Authentication

The CarbonArc MCP server requires OAuth 2.0 authentication. Most examples
handle this by:

1. **Discovering** the authorization endpoint via the server's
   `/.well-known/oauth-protected-resource` metadata.
2. **Registering** a client dynamically (RFC 7591) if needed.
3. **Opening a browser** for the user to log in.
4. **Receiving the auth code** on a local callback server.
5. **Exchanging** the code for an access token.

All three client examples implement the full browser-based OAuth flow.

## Security

These examples are **minimal quickstarts** designed to demonstrate MCP connectivity. They store OAuth tokens in memory for the duration of the script and do **not** include:

- Token encryption at rest
- Audit logging
- Rate limiting
- Secret scanning or rotation

For production-grade security patterns — including AES-256-GCM token encryption, SOC 2 controls mapping, and reference implementations for AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, and HashiCorp Vault — see the [Carbon Arc Partnership Examples](https://github.com/Carbon-Arc/carbonarc-partnership-examples) repository, specifically [embedded-ca-user/SECURITY.md](https://github.com/Carbon-Arc/carbonarc-partnership-examples/blob/main/embedded-ca-user/SECURITY.md).

## Further Reading

- [Model Context Protocol Specification](https://modelcontextprotocol.io/specification/latest/)
- [MCP SDKs Overview](https://modelcontextprotocol.io/docs/sdk)
- [CarbonArc Documentation](http://docs.carbonarc.ai/)
