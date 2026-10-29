---
"chainlink": minor
---

#added Added OIDC login support to the CLI via the OAuth 2.0 device authorization grant (RFC 8628). When a node is configured with `AuthenticationMethod = 'oidc'`, `chainlink admin login` (with no credentials file) now performs a browser-based device flow against the identity provider, brokered by the node. The local email/password path is unchanged and remains available as a break-glass admin.

#changed The OIDC authorization-code (operator UI) flow now always uses PKCE (RFC 7636). The OIDC `ClientSecret` is now optional: confidential clients still send it, while public clients (required for the device flow) rely on PKCE instead. Existing confidential-client deployments are unaffected.

#changed Hardened OIDC device and auth-code flows: dedicated pre-auth rate limits and per-IP concurrency caps on device endpoints, server-side PKCE verifier and OIDC nonce (the session cookie holds only the anti-CSRF state), email claim validation, orphan-session cleanup if cookie bind fails, and a stricter CLI `/oidc-enabled` probe. Browser login state and device-flow results are stored in Postgres so any replica can finish a login. The provider device code stays on the replica that started the flow and is not written to the database. Device-flow lifetime is capped at five minutes. A verification URI whose host does not match the configured provider is rejected. When the socket peer is a private proxy, device-flow rate limits use the rightmost X-Forwarded-For address; a public peer is used as-is. Denied device flows are audited. A groups claim that contains any value other than the configured role groups is rejected. OIDC API tokens cannot be created or used. The CLI refuses device login when TLS verification is disabled.
