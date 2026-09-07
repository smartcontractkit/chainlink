# vault-dkg-decode

Decodes an exported vault DKG **result package** and prints the TDH2 master
public key in the exact form CapReg stores it (the `VaultPublicKey`
`stringValue`).

Run this on every vault DKG **reshare** (or fresh dealing). On a reshare the
master key (`G_bar`/`H`) is preserved but the per-recipient verification shares
(`HArray`) are redistributed to the new committee — so CapReg's `VaultPublicKey`
must be refreshed to the value this tool prints, or `secrets.get`
decryption-share verification fails (`failed to verify decryption share`).

## Get a result package

Export it from any node that holds it (member of the instance's committee), via
the node API. Authenticate first, then export and pull out `hexDKGResultPackage`:

```bash
# 1. authenticate to a cookie jar (node API email/password)
curl -s -c cookies.txt -X POST "$NODE_URL/sessions" \
  -H 'Content-Type: application/json' \
  -d '{"email":"'"$NODE_EMAIL"'","password":"'"$NODE_PASSWORD"'"}'

# 2. export and extract the hex (adjust the export path to your build)
curl -s -b cookies.txt "$NODE_URL/.../export?instanceID=<instanceID>" \
  | jq -r .hexDKGResultPackage > resultpkg.hex
```

## Decode

```bash
cd core/scripts/cre/vault-dkg-decode

# full human-readable summary + the CapReg stringValue at the end
go run . ./resultpkg.hex

# or pipe straight from the export
curl ... | jq -r .hexDKGResultPackage | go run . -

# just the CapReg VaultPublicKey stringValue (what you paste into the
# capability_registry_update_don input)
go run . --capreg ./resultpkg.hex

# just the decoded pubkey JSON
go run . --json ./resultpkg.hex
```

## Field meaning

| Field    | Meaning                                              | Reshare?              |
|----------|------------------------------------------------------|-----------------------|
| `Group`  | curve name (`P256`)                                  | never changes         |
| `H`      | master public key — what secrets are encrypted to    | **stable**            |
| `G_bar`  | deterministically-derived generator (from `H`)       | **stable**            |
| `HArray` | per-recipient public verification shares (one/member)| **changes on reshare**|

`G_bar` + `H` drive encryption (`secrets.create`/`update`), so they stay valid
after a reshare. `HArray` verifies decryption shares (`secrets.get`), so a stale
`HArray` in CapReg is what breaks reads after a reshare — republish the value
this tool prints, then reboot the gateways to drop their cached master key.
