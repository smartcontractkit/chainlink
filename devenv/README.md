<div align="center">

# Chainlink Developer Environment

`NodeSet` + `Anvil` + `Fake Server` + `OCR2 Product Orchestration`

</div>

## Run the Environment

We are using [Justfile](https://github.com/casey/just?tab=readme-ov-file#cross-platform)

Change directory to `devenv` then

```bash
brew install just # click the link above if you are not on OS X
just build-fakes && just cli && cl sh
```

Then start the observability stack and run the soak test, perform actions in `ccv sh` console

```bash
up
obs up -f # Click OCR2 Dashboard link to open in another tab
test load # Run the load test, you'll see OCR2 rounds stats
```

## Run with custom CL image

Use `up env.toml,env-cl-rebuild.toml` to rebuild custom CL image from your local `chainlink` repository.

## Updating Fakes

Fake represent a controlled External Adapter that returns feed values.

```bash
just build-fakes <aws_registry> # use SDLC registry
just push-fakes <aws_registry> # use SDLC registry
```

## Adding Products

To extend the environment all you need to do is to implement the [interface](interface.go) and add a switch clause in [environment](environment.go)
