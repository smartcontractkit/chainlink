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
    env:
      MACOS_SDK_VERSION: 12.3
    steps:
      - name: Checkout repository
        uses: actions/checkout@b4ffde65f46336ab88eb53be808477a3936bae11 # v4.1.1
      - name: Configure aws credentials
        uses: aws-actions/configure-aws-credentials@010d0da01d0b5a38af31e9c3470dbfdabdecca3a # v4.0.1
        with:
          role-to-assume: ${{ secrets.aws-role-arn }}
          role-duration-seconds: ${{ secrets.aws-role-dur-sec }}
          aws-region: ${{ secrets.aws-region }}
      - name: Cache macos sdk
        id: sdk-cache
        uses: actions/cache@v3
        with:
          path: ${{ format('MacOSX{0}.sdk', env.MAC_SDK_VERSION) }}
          key: ${{ runner.OS }}-${{ env.MAC_SDK_VERSION }}-macos-sdk-cache-${{ hashFiles('**/SDKSettings.json') }}
          restore-keys: |
            ${{ runner.OS }}-${{ env.MAC_SDK_VERSION }}-macos-sdk-cache-
      - name: Get macos sdk
        if: steps.sdk-cache.outputs.cache-hit != 'true'
        run: |
          curl -L https://github.com/joseluisq/macosx-sdks/releases/download/${MACOS_SDK_VERSION}/MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz > MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz
          tar -xf MacOSX${MACOS_SDK_VERSION}.sdk.tar.xz
      - name: Build, sign, and publish
        uses: ./.github/actions/goreleaser-build-sign-publish
        with:
          enable-docker-publish: "true"
          enable-goreleaser-snapshot: "false"
          docker-registry: ${{ secrets.aws-ecr-registry }}
          goreleaser-exec: goreleaser
          goreleaser-config: .goreleaser.yaml
          macos-sdk-dir: ${{ format('MacOSX{0}.sdk', env.MAC_SDK_VERSION) }}
        env:
          GITHUB_TOKEN: ${{ secrets.gh-token }}
```

### snapshot release

```yaml
- name: Build, sign, and publish image
  uses: ./.github/actions/goreleaser-build-sign-publish
  with:
    enable-docker-publish: "true"
    enable-goreleaser-snapshot: "true"
    docker-registry: ${{ secrets.aws-ecr-registry }}
    goreleaser-exec: goreleaser
    goreleaser-config: .goreleaser.yaml
```

### image signing

```yaml
- name: Build, sign, and publish
  uses: ./.github/actions/goreleaser-build-sign-publish
  with:
    enable-docker-publish: "true"
    enable-goreleaser-snapshot: "false"
    enable-cosign: "true"
    docker-registry: ${{ secrets.aws-ecr-registry }}
    goreleaser-exec: goreleaser
    goreleaser-config: .goreleaser.yaml
    cosign-password: ${{ secrets.cosign-password }}
    cosign-public-key: ${{ secrets.cosign-public-key }}
    cosign-private-key: ${{ secrets.cosign-private-key }}
    macos-sdk-dir: MacOSX12.3.sdk
```

## customizing

### inputs

Following inputs can be used as `step.with` keys

| Name                         | Type   | Default            | Description                                                             |
| ---------------------------- | ------ | ------------------ | ----------------------------------------------------------------------- |
| `goreleaser-version`         | String | `1.13.1`           | `goreleaser` version                                                    |
| `zig-version`                | String | `0.10.0`           | `zig` version                                                           |
| `cosign-version`             | String | `v1.13.1`          | `cosign` version                                                        |
| `macos-sdk-dir`              | String | `MacOSX12.3.sdk`   | MacOSX sdk directory                                                    |
| `enable-docker-publish`      | Bool   | `true`             | Enable publishing of Docker images / manifests                          |
| `docker-registry`            | String | `localhost:5001`   | Docker registry                                                         |
| `enable-goreleaser-snapshot` | Bool   | `false`            | Enable goreleaser build / release snapshot                              |
| `goreleaser-exec`            | String | `goreleaser`       | The goreleaser executable, can invoke wrapper script                    |
| `goreleaser-config`          | String | `.goreleaser.yaml` | The goreleaser configuration yaml                                       |
| `enable-cosign`              | Bool   | `false`            | Enable signing of Docker images                                         |
| `cosign-public-key`          | String | `""`               | The public key to be used with cosign for verification                  |
| `cosign-private-key`         | String | `""`               | The private key to be used with cosign to sign the image                |
| `cosign-password-key`        | String | `""`               | The password to decrypt the cosign private key needed to sign the image |

## testing

- bring up local docker registry

```sh
docker run -d --restart=always -p "127.0.0.1:5001:5000" --name registry registry:2
```

- run snapshot release, publish to local docker registry

```sh
GORELEASER_EXEC="<goreleaser-wrapper" \
GORELEASER_CONFIG=".goreleaser.yaml" \
ENABLE_GORELEASER_SNAPSHOT=true \
ENABLE_DOCKER_PUBLISH=true \
DOCKER_MANIFEST_EXTRA_ARGS="--insecure" \
./.github/actions/goreleaser-build-sign-publish/action_utils goreleaser_release
```
