---
"chainlink": minor
---

#added Added OIDC login support to the CLI via the OAuth 2.0 device authorization grant (RFC 8628). When a node is configured with `AuthenticationMethod = 'oidc'`, `chainlink admin login` (with no credentials file) now performs a browser-based device flow against the identity provider, brokered by the node. The local email/password path is unchanged and remains available as a break-glass admin.

#changed The OIDC authorization-code (operator UI) flow now always uses PKCE (RFC 7636). The OIDC `ClientSecret` is now optional: confidential clients still send it, while public clients (required for the device flow) rely on PKCE instead. Existing confidential-client deployments are unaffected.

#changed Hardened OIDC device and auth-code flows: server-side PKCE verifier and OIDC nonce (the session cookie holds only the anti-CSRF state), email claim validation, orphan-session cleanup if cookie bind fails, and a stricter CLI `/oidc-enabled` probe. Browser login state and device-flow results are stored in Postgres so any replica can finish a login. The provider device code stays on the replica that started the flow and is not written to the database. Device-flow lifetime is capped at five minutes. A verification URI whose host does not match the configured provider is rejected. Groups the node does not know are ignored; the highest configured role group picks the role. A denied device login tells the CLI why, including when none of the user's groups map to a node role. OIDC API tokens cannot be created or used. Debug request logs redact the OIDC authorization code and state.

#added `WebServer.OIDC.DeviceFlowAllowedCIDRs` lists the client networks (for example the operator VPN or bastion range) that may use the CLI device login. It defaults to empty, which turns the device login off; requests from outside the list get 403. Browser login is not affected.

#added `WebServer.OIDC.TrustedProxyCIDRs` lists the load balancers and proxies in front of the node. The OIDC login routes believe `X-Forwarded-For` only from these hops and take the first address outside them as the client. It defaults to empty, which uses the socket peer and ignores the header.

#changed The node no longer passes the identity provider's `verification_uri_complete` to the CLI, and the CLI no longer prints a pre-filled approval link. The operator types the user code at the identity provider, which removes the one-click approval a forwarded link would allow.

#changed The CLI refuses device login unless `--remote-node-url` uses https (plain http only for localhost) and TLS verification is on, and refuses to show a verification URI that is not https. Note the CLI does not check the verification URI host itself; the node checks it against the configured provider.

#changed Device-flow abuse limits: dedicated pre-auth rate limits on `/oidc-device/start`, `/oidc-device/poll` and `/oidc-login`, keyed per IPv4 address or IPv6 /64. A device flow the CLI stops polling for 30 seconds (or three poll intervals, if longer) releases its slot and stops polling the identity provider. `/oidc-login` caps pending sign-ins at 50 per client and 10,000 in total, and answers 503 with `Retry-After` past either.

#added `AUTH_DEVICE_FLOW_STARTED` audit event with the client IP, User-Agent and a handle fingerprint. Device-flow success and failure events carry the same fingerprint and the starting client IP; success also records the redeeming IP and `client_ip_mismatch` when the two differ.
