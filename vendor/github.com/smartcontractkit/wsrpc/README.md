# Websockets RPC

Establishes a persistent bi-directional communication channel using mTLS and websockets.

## Set up

In order to generate a service definition you will need the wsrpc protoc plugin.

Build the protoc plugin

```
cd cmd/protoc-gen-go-wsrpc
go build
```

Place the resulting binary in your GOPATH.

In the directory containing your protobuf service definition, run:

```
protoc --go_out=. --go_opt=paths=source_relative \
--go-wsrpc_out=. \
--go-wsrpc_opt=paths=source_relative yourproto.proto
```

This will generate the service definitions in *_wsrpc.pb.go

## Usage

### Client to Server RPC

Implement handlers for the server
```go
type pingServer struct {}

func (s *pingServer) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	// Extracts the connection client's public key.
	// You can use this to identify the client
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}
	pubKey := p.PublicKey

	fmt.Println(pubKey)

	return &pb.PingResponse{
		Body: "Pong",
	}, nil
}
```

Initialize a server with the server's private key and a slice of allowable public keys.

```go
lis, err := net.Listen("tcp", "127.0.0.1:1337")
if err != nil {
	log.Fatalf("[MAIN] failed to listen: %v", err)
}
s := wsrpc.NewServer(wsrpc.Creds(privKey, pubKeys))
// Register the ping server implementation with the wsrpc server
pb.RegisterPingServer(s, &pingServer{})

s.Serve(lis)
```

Initialize a client with the client's private key and the server's public key

```go
conn, err := wsrpc.Dial("127.0.0.1:1338", wsrpc.WithTransportCreds(privKey, serverPubKey))
if err != nil {
	log.Fatalln(err)
}
defer conn.Close()

// Initialize a new wsrpc client caller
// This is used to called RPC methods on the server
c := pb.NewPingClient(conn)

c.Ping(context.Background(), &pb.Ping{Body: "Ping"})
```

### Server to Client RPC

Implement handlers for the client

```go
type pingClient struct{}

func (c *pingClient) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Body: "Pong",
	}, nil
}
```

Initialize a server with the server's private key and a slice of allowable public keys.

```go
lis, err := net.Listen("tcp", "127.0.0.1:1337")
if err != nil {
	log.Fatalf("[MAIN] failed to listen: %v", err)
}
s := wsrpc.NewServer(wsrpc.Creds(privKey, pubKeys))
c := pb.NewPingClient(s)

s.Serve(lis)

// Call the RPC method with the pub key so we know which connection to send it to
// otherwise it will error.
ctx := peer.NewCallContext(context.Background(), pubKey)
c.Ping(ctx, &pb.PingRequest{Body: "Ping"})
```

Initialize a client with the client's private key and the server's public key

```go
conn, err := wsrpc.Dial("127.0.0.1:1337", wsrpc.WithTransportCreds(privKey, serverPubKey))
if err != nil {
	log.Fatalln(err)
}
defer conn.Close()

// Initialize RPC call handlers on the client connection
pb.RegisterPingServer(conn, &pingClient{})
```

## Example

You can run a simple example where both the client and server implement a Ping service, and perform RPC calls to each other every 5 seconds.

1. Run the server in `examples/simple/server` with `go run main.go`
2. Run a client (Alice) in `examples/simple/server` with `go run main.go 0`
3. Run a client (Bob) in `examples/simple/server` with `go run main.go 1`
4. Run a invalid client (Charlie) in `examples/simple/server` with `go run main.go 2`. The server will reject this connection.

While the client's are connected, kill the server and see the client's enter a backoff retry loop. Start the server again and they will reconnect.

## TODO

- [ ] Improve Tests
- [ ] Return a response status
- [x] Add a Blocking DialOption