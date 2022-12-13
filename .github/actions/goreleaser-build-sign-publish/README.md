# goreleaser-build-sign-publish

## workflow

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
