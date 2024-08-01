# 0.3.0

* Use context to cancel a blocking `Dial`.

  ```
  // With timeout
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := wsrpc.DialWithContext(ctx, "127.0.0.1:1338",
		wsrpc.WithTransportCreds(privKey, serverPubKey),
		wsrpc.WithBlock(),
	)

  // Manual cancel
  ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
    conn, err := wsrpc.DialWithContext(ctx, "127.0.0.1:1338",
		wsrpc.WithTransportCreds(privKey, serverPubKey),
		wsrpc.WithBlock(),
	)}()

  // Something causes the need to cancel.
  cancel()
  ```

* Subscribe to a notification channel to get updates when a client connection is
  established or dropped. You can then retrieve the latest list of keys

  ```
	go func() {
		for {
			notifyCh := s.GetConnectionNotifyChan()
			<-notifyCh

			fmt.Println("Connected to:", s.GetConnectedPeerPublicKeys())
		}
	}()
  ```

# 0.2.0

* Replace metadata public key context with with a peer context.

  **Extracting a public key**
  ```
    // Previously
	pubKey, ok := metadata.PublicKeyFromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract public key")
	}

    // Now
    p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("could not extract peer information")
	}
    pubKey := p.PublicKey
  ```

  **Making a server side RPC call**
  ```
  // Previously
  ctx := context.WithValue(context.Background(), metadata.PublicKeyCtxKey, pubKey)
  res, err := c.Gnip(ctx, &pb.GnipRequest{Body: "Gnip"})

  // Now
  ctx := peer.NewCallContext(context.Background(), pubKey)
  res, err := c.Gnip(ctx, &pb.GnipRequest{Body: "Gnip"})
  ```
* Add a `WithBlock` DialOption which blocks the caller of Dial until the underlying connection is up.

# 0.1.1

## Changed

* Supress logging until we can implement a configurable logging solution.

# 0.1.0

Initial release