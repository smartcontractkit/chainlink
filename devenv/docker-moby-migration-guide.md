# Migration Guide: github.com/docker/docker → github.com/moby/moby split modules

This guide documents all breaking changes encountered when migrating from the
monolithic `github.com/docker/docker` and `github.com/docker/go-connections/nat`
packages to the new split Moby modules:
- `github.com/moby/moby/api v1.54.x`
- `github.com/moby/moby/client v0.3.x`

---

## 1. go.mod changes

**Before:**
```
github.com/docker/docker v28.x.x+incompatible
github.com/docker/go-connections v0.x.x
```

**After:**
```
github.com/moby/moby/api v1.54.1-0.20260401134807-948d5691a093
github.com/moby/moby/client v0.3.1-0.20260401134807-948d5691a093
github.com/docker/docker v28.x.x+incompatible  // keep as indirect
```

> **Gotcha**: If both `github.com/moby/moby v28.x.x+incompatible` (the old monolith)
> and `github.com/moby/moby/api` (the new split module) are present, `go mod tidy`
> will report "ambiguous import". Keep only the split modules; let `docker/docker` stay
> as an indirect dep for things that still need it transitively.

---

## 2. Port types: nat.Port → network.Port

### Import changes
```go
// Before
import "github.com/docker/go-connections/nat"

// After
import (
    "net/netip"
    "github.com/moby/moby/api/types/network"
)
```

### PortMap / PortBinding
```go
// Before
h.PortBindings = nat.PortMap{
    nat.Port("9000/tcp"): []nat.PortBinding{
        {HostIP: "0.0.0.0", HostPort: "9000"},
    },
}

// After
h.PortBindings = network.PortMap{
    network.MustParsePort("9000/tcp"): []network.PortBinding{
        {HostIP: netip.MustParseAddr("0.0.0.0"), HostPort: "9000"},
    },
}
```

Key differences:
- `nat.PortMap` = `map[nat.Port][]nat.PortBinding` where `nat.Port` is a **string alias**
- `network.PortMap` = `map[network.Port][]network.PortBinding` where `network.Port` is a **struct**
- `network.PortBinding.HostIP` is `netip.Addr` (not `string`)
- Use `network.MustParsePort("NNN/tcp")` for known-valid strings (panics on invalid)
- Use `network.ParsePort("NNN/tcp")` to get `(Port, error)` for runtime values

### network.Port methods
```go
p, _ := network.ParsePort("9000/tcp")
p.Port()   // → "9000"   (string, port number only)
p.Num()    // → 9000     (uint16)
p.Proto()  // → "tcp"    (IPProtocol)
p.String() // → "9000/tcp"
```

### PortMap construction pattern (inside closures/HostConfigModifier)
Pre-compute keys before the closure to handle errors properly:
```go
portKey, err := network.ParsePort(fmt.Sprintf("%s/tcp", portStr))
if err != nil {
    return nil, err
}
req := testcontainers.ContainerRequest{
    HostConfigModifier: func(h *container.HostConfig) {
        h.PortBindings = network.PortMap{
            portKey: []network.PortBinding{{
                HostIP:   netip.MustParseAddr("0.0.0.0"),
                HostPort: portStr,
            }},
        }
    },
}
```

Or use `network.MustParsePort` directly for compile-time known strings:
```go
HostConfigModifier: func(h *container.HostConfig) {
    h.PortBindings = network.PortMap{
        network.MustParsePort(containerPort): []network.PortBinding{{
            HostIP:   netip.MustParseAddr("0.0.0.0"),
            HostPort: in.Port,
        }},
    }
},
```

---

## 3. testcontainers-go API changes (v0.41+)

### ForListeningPort
```go
// Before
wait.ForListeningPort(nat.Port("9000/tcp"))

// After
wait.ForListeningPort("9000/tcp")  // plain string
```

### ForHTTP().WithPort()
```go
// Before
wait.ForHTTP("/health").WithPort(nat.Port("8080/tcp"))

// After
wait.ForHTTP("/health").WithPort("8080/tcp")  // plain string
```

### Container.MappedPort()
```go
// Before: takes nat.Port, returns nat.Port
ep, err := c.MappedPort(ctx, nat.Port("9000/tcp"))
ep.Port()  // → "9000"

// After: takes string, returns network.Port
ep, err := c.MappedPort(ctx, "9000/tcp")
ep.Port()  // → "9000"  (same method, still returns string)
```

---

## 4. Docker Client API changes (moby/client v0.3+)

The client API is a **complete redesign**: every method now takes a single `*Options`
struct and returns a `*Result` struct (or similar).

### Import
```go
// Before
"github.com/docker/docker/client"

// After
"github.com/moby/moby/client"
```

### Ping
```go
// Before
cli.Ping(ctx)

// After
cli.Ping(ctx, client.PingOptions{})
```

### Container Exec
```go
// Before
import "github.com/docker/docker/api/types/container"
execID, err := cli.ContainerExecCreate(ctx, id, container.ExecOptions{...})
resp, err := cli.ContainerExecAttach(ctx, execID.ID, container.ExecStartOptions{})

// After
execID, err := cli.ExecCreate(ctx, id, client.ExecCreateOptions{...})
resp, err := cli.ExecAttach(ctx, execID.ID, client.ExecAttachOptions{})
```

### Container List
```go
// Before
containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
for _, c := range containers { ... }

// After
containers, err := cli.ContainerList(ctx, client.ContainerListOptions{All: true})
for _, c := range containers.Items { ... }  // note: .Items
```

### Container Inspect
```go
// Before
inspected, err := cli.ContainerInspect(ctx, id)
_ = inspected.Config
_ = inspected.HostConfig
_ = inspected.NetworkSettings.Networks

// After
inspected, err := cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
_ = inspected.Container.Config
_ = inspected.Container.HostConfig
_ = inspected.Container.NetworkSettings.Networks
```

### Container Stop / Remove / Create / Start
```go
// After
_, err = cli.ContainerStop(ctx, name, client.ContainerStopOptions{})
_, err = cli.ContainerRemove(ctx, name, client.ContainerRemoveOptions{RemoveVolumes: false})
createResp, err := cli.ContainerCreate(ctx, client.ContainerCreateOptions{
    Config:           inspected.Container.Config,
    HostConfig:       inspected.Container.HostConfig,
    NetworkingConfig: networkingConfig,
    Name:             name,
})
_, err = cli.ContainerStart(ctx, createResp.ID, client.ContainerStartOptions{})
```

### Copy to Container
```go
// Before
err = cli.CopyToContainer(ctx, id, targetPath, &buf, container.CopyToContainerOptions{...})

// After
_, err = cli.CopyToContainer(ctx, id, client.CopyToContainerOptions{
    DestinationPath:         targetPath,
    Content:                 &buf,
    AllowOverwriteDirWithFile: true,
})
```

### Network Connect
```go
// Before (4 args)
cli.NetworkConnect(ctx, networkName, containerID, &endpointSettings)

// After (single Options struct)
_, err = cli.NetworkConnect(ctx, networkName, client.NetworkConnectOptions{
    Container: containerID,
    EndpointConfig: &networkTypes.EndpointSettings{
        Aliases: []string{alias},
    },
})
```

### Filters (ContainerList, etc.)
The `github.com/docker/docker/api/types/filters` package **does not exist** in the split modules.

```go
// Before
import dfilter "github.com/docker/docker/api/types/filters"
args := dfilter.NewArgs(dfilter.Arg("label", "framework=ctf"))
opts := container.ListOptions{Filters: args}

// After (client.Filters is map[string]map[string]bool with Add method)
filters := make(client.Filters).Add("label", "framework=ctf")
opts := client.ContainerListOptions{All: true, Filters: filters}
```

---

## 5. Mount types

```go
// Before
import "github.com/docker/docker/api/types/mount"

// After
import "github.com/moby/moby/api/types/mount"
```
Types and fields are identical — just the import path changes.

---

## 6. HostConfig

```go
// Before
import "github.com/docker/docker/api/types/container"
// HostConfig is container.HostConfig

// After
import "github.com/moby/moby/api/types/container"
// HostConfig is still container.HostConfig — same package name, different module
```

---

## 7. Iterating over network.PortMap

```go
// Migrating functions that iterate over a port map
// Before: nat.PortMap keys are nat.Port (string alias with .Int() method)
for port, bindings := range portMap {
    portNum := strconv.Itoa(port.Int())
}

// After: network.PortMap keys are network.Port struct
for port, bindings := range portMap {
    portNum := port.Port()   // returns string like "9000"
    // or: strconv.Itoa(int(port.Num())) for uint16
}
```

---

## 8. ContainerInspect result for NetworkSettings

```go
// Before
inspected, _ := cli.ContainerInspect(ctx, id)
networks := inspected.NetworkSettings.Networks  // map[string]*network.EndpointSettings

// After
inspected, _ := cli.ContainerInspect(ctx, id, client.ContainerInspectOptions{})
networks := inspected.Container.NetworkSettings.Networks
```

---

## 9. Things to search for in your codebase

Grep patterns to find all affected code:
```bash
grep -r "docker/go-connections/nat" --include="*.go" .
grep -r "docker/docker/api/types" --include="*.go" .
grep -r "docker/docker/client" --include="*.go" .
grep -r "nat\.Port\(" --include="*.go" .
grep -r "nat\.PortMap" --include="*.go" .
grep -r "nat\.PortBinding" --include="*.go" .
grep -r "ContainerExecCreate\|ContainerExecAttach" --include="*.go" .
grep -r "ContainerList\|ContainerInspect\|NetworkConnect" --include="*.go" .
grep -r "dfilter\.\|filters\.NewArgs" --include="*.go" .
```

---

## 10. Template / code generation

If you generate Go code in string templates, update those too:
- Replace `nat.PortMap{...}` → `network.PortMap{network.MustParsePort(...): ...}`
- Replace `[]nat.PortBinding{{HostIP: "0.0.0.0", ...}}` → `[]network.PortBinding{{HostIP: netip.MustParseAddr("0.0.0.0"), ...}}`
- Update import strings in templates accordingly
