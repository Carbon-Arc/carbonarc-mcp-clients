/**
 * OAuth Helper — browser-based authorization code flow for MCP clients
 * ====================================================================
 *
 * Provides a reusable OAuthClientProvider and callback server
 * so MCP client examples stay focused on the protocol, not on OAuth plumbing.
 */

import { createServer, type IncomingMessage, type ServerResponse } from "node:http";
import { exec } from "node:child_process";
import type {
  OAuthClientProvider,
  OAuthClientMetadata,
  OAuthClientInformationFull,
  OAuthTokens,
} from "@modelcontextprotocol/sdk/client/auth.js";

/**
 * Minimal in-memory OAuth provider that opens the user's browser
 * for the authorization code flow.
 */
export class BrowserOAuthProvider implements OAuthClientProvider {
  private _clientInfo?: OAuthClientInformationFull;
  private _tokens?: OAuthTokens;
  private _codeVerifier?: string;
  private callbackUrl: string;

  constructor(callbackUrl: string) {
    this.callbackUrl = callbackUrl;
  }

  get redirectUrl(): string {
    return this.callbackUrl;
  }

  get clientMetadata(): OAuthClientMetadata {
    return {
      client_name: "CarbonArc MCP Client Example",
      redirect_uris: [this.callbackUrl],
      grant_types: ["authorization_code", "refresh_token"],
      response_types: ["code"],
      token_endpoint_auth_method: "client_secret_post",
    };
  }

  clientInformation() {
    return this._clientInfo;
  }
  saveClientInformation(info: OAuthClientInformationFull) {
    this._clientInfo = info;
  }

  tokens() {
    return this._tokens;
  }
  saveTokens(tokens: OAuthTokens) {
    this._tokens = tokens;
  }

  redirectToAuthorization(url: URL) {
    console.log(`\nOpening browser for login: ${url}\n`);
    exec(`open "${url}"`, (err) => {
      if (err) console.log(`Open this URL manually: ${url}`);
    });
  }

  saveCodeVerifier(v: string) {
    this._codeVerifier = v;
  }
  codeVerifier(): string {
    if (!this._codeVerifier) throw new Error("No code verifier saved");
    return this._codeVerifier;
  }
}

/**
 * Starts a tiny HTTP server and waits for the OAuth callback
 * to deliver the authorization code.
 */
export function waitForAuthCode(port: number): Promise<string> {
  return new Promise((resolve, reject) => {
    const server = createServer(
      (req: IncomingMessage, res: ServerResponse) => {
        if (req.url === "/favicon.ico") {
          res.writeHead(404);
          res.end();
          return;
        }
        const parsed = new URL(req.url || "", `http://localhost:${port}`);
        const code = parsed.searchParams.get("code");
        if (code) {
          res.writeHead(200, { "Content-Type": "text/html" });
          res.end("<h2>Authorized — you can close this tab.</h2>");
          resolve(code);
          setTimeout(() => server.close(), 1000);
        } else {
          res.writeHead(400);
          res.end("Missing authorization code");
          reject(new Error("No code in callback"));
        }
      }
    );
    server.listen(port, () => {
      console.log(`OAuth callback listening on http://localhost:${port}`);
    });
  });
}
