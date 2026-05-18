# goreleaser-build-sign-publish

> goreleaser wrapper action

## workflows

### build publish

```yaml
name: goreleaser

on:
  push:
    tags:
      - "v*"

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    environment: release
    permissions:
      id-token: write
      contents: read
    steps:
      - name: Checkout repository
        uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1
      - name: Configure aws credentials
        uses: aws-actions/configure-aws-credentials@010d0da01d0b5a38af31e9c3470dbfdabdecca3a # v4.0.1
        with:
          role-to-assume: ${{ secrets.aws-role-arn }}
          role-duration-seconds: ${{ secrets.aws-role-dur-sec }}
          aws-region: ${{ secrets.aws-region }}
      - name: Build, sign, and publish
        uses: ./.github/actions/goreleaser-build-sign-publish
        with:
          docker-registry: ${{ secrets.aws-ecr-registry }}
          goreleaser-config: .goreleaser.yaml
        env:
          GITHUB_TOKEN: ${{ secrets.gh-token }}
```

### snapshot release

```yaml
- name: Build, sign, and publish image
  uses: ./.github/actions/goreleaser-build-sign-publish
  with:
    docker-registry: ${{ secrets.aws-ecr-registry }}
    goreleaser-config: .goreleaser.yaml
```

## customizing

### inputs

Following inputs can be used as `step.with` keys

| Name                         | Type   | Default            | Description                                                             |
| ---------------------------- | ------ | ------------------ | ----------------------------------------------------------------------- |
| `goreleaser-version`         | String | `~> v2`            | `goreleaser` version                                                    |
| `docker-registry`            | String | `localhost:5001`   | Docker registry                                                         |
| `docker-image-tag`           | String | `develop`          | Docker image tag                                                        |
| `goreleaser-config`          | String | `.goreleaser.yaml` | The goreleaser configuration yaml                                       |

## testing

- bring up local docker registry

```sh
docker run -d --restart=always -p "127.0.0.1:5001:5000" --name registry registry:2
```

- run snapshot release, publish to local docker registry

```sh
GORELEASER_CONFIG=".goreleaser.yaml" \
DOCKER_MANIFEST_EXTRA_ARGS="--insecure" \
./.github/actions/goreleaser-build-sign-publish/action_utils goreleaser_release
```
