---
"chainlink": minor
---

#added Added OIDC login support to the CLI via the OAuth 2.0 device authorization grant (RFC 8628). When a node is configured with `AuthenticationMethod = 'oidc'`, `chainlink admin login` (with no credentials file) now performs a browser-based device flow against the identity provider, brokered by the node. The local email/password path is unchanged and remains available as a break-glass admin.

#changed The OIDC authorization-code (operator UI) flow now always uses PKCE (RFC 7636). The OIDC `ClientSecret` is now optional: confidential clients still send it, while public clients (required for the device flow) rely on PKCE instead. Existing confidential-client deployments are unaffected.
