# dispatch-release-fanout

Tells the internal deployment hub which images a release published. The hub
classifies each tag into a release channel and dispatches the image-bump
workflow in every deployment repo subscribed to that `(stream, channel)`, which
opens the tag-bump PR there.

One dispatch carries every image the build published. This repo holds no
deployment topology: which repo consumes `core` vs `ccip`, and which channel
each environment tracks, is configured in the hub and in the deployment repos.

```yaml
- uses: ./.github/actions/dispatch-release-fanout
  with:
    gati-profile: ${{ secrets.GATI_PROFILE_IMAGE_BUMP_HUB }}
    hub-repo: ${{ secrets.REPO_IMAGE_BUMP_HUB }}
    images: |
      core=2.65.1-rc.0
      ccip=2.65.1-ccip-rc.0
```

## Required secrets

This repository is public and the hub is internal, so neither the hub's name
nor the GATI profile that names it is committed here.

| Secret                         | Value                                                                       |
| ------------------------------ | --------------------------------------------------------------------------- |
| `REPO_IMAGE_BUMP_HUB`          | `owner/name` of the hub repository                                          |
| `GATI_PROFILE_IMAGE_BUMP_HUB`  | GATI v2 profile granting `actions:write` on that repo (a `presets/…` entry) |

The calling job needs `id-token: write` — the profile is minted from this
repository's OIDC identity, and grants nothing beyond starting one workflow in
the hub. It cannot read or write any repository's contents.

## Behaviour

- `images` entries with an empty tag are dropped, so a caller can pass an image
  whose build failed or was skipped without special-casing it. If every entry
  is empty the action warns and exits 0 without dispatching.
- A tag which is unrecognized by the hub for its stream is *not routed*
  and does not fail: the hub warns and still delivers the other images in the
  same dispatch.
- The dispatch is fire-and-forget. It returns as soon as the hub accepts it,
  and reports nothing about what the hub did — read the hub's own run, linked
  by the `correlation-id` in this action's job summary.
