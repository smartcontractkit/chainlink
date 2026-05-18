## Running the Liquidity Manager DON Locally

Start by cloning the chainlink CCIP repo and checking out the branch with the scripts:

```bash
git clone https://github.com/smartcontractkit/ccip.git
cd ccip
```

Change directories to the appropriate directory:

```bash
cd core/scripts/ccip/liquiditymanager
```

Before running the setup script, you need to have a postgres database running. Here’s how I run it:

```bash
# in a new shell
docker run -d --rm --name chainlink-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_HOST_AUTH_METHOD=trust -v $HOME/chainlink-pg-data/:/var/lib/postgresql/data -p 5432:5432 postgres:14 postgres -N 500 -B 1024MB
```

One more thing to do before you can run the setup script is create an env file with the networks you are going to deploy to. Here’s how mine looks like:

```bash
# Arbitrum sepolia RPCs
# Note both RPC_ and WS_ are needed
export RPC_421614=<http-rpc>
export WS_421614=<ws-rpc>
# Sepolia RPCs
export RPC_11155111=<http-rpc>
export WS_11155111=<ws-rpc>

# Private key with funds on both chains above
# It will be used to deploy the contracts
# and fund the chainlink node sending keys
export OWNER_KEY=<private key hex, no leading 0x>
```

If these RPC’s no longer work be sure to check Notion for up to date URLs. You need to get both WebSocket and HTTPS URLs in order for the node to work.

Once the database is up and running, you can go back to the scripts and run the following command:

```bash
# provide the env vars to the current shell
source env_file_from_above.sh
# start the setup
go run . setup-liquiditymanager-nodes
```

The setup script will do the following in order:

1. Deploy the relevant contracts. Here’s a sample output:

```
Executing contract deployment, TX: https://sepolia.etherscan.io/tx/0x7367eb34258c31a430dcd40ce466011a329a41608f79af7d0fb0a722e6d73181
Contract Address: 0xabB73a46050C14Aee721fA3a6c9333709ED806b5
Contract explorer link: https://sepolia.etherscan.io/address/0xabB73a46050C14Aee721fA3a6c9333709ED806b5
Executing contract deployment, TX: https://sepolia.etherscan.io/tx/0xdb2364e2ae1f9418ef62e5b0cb0dfdbb062409d302feaace55f6824f7726ddd3
Contract Address: 0xEeA421a17746A00C2Bc7B5869543cd2A89aC7125
Contract explorer link: https://sepolia.etherscan.io/address/0xEeA421a17746A00C2Bc7B5869543cd2A89aC7125
Executing contract deployment, TX: https://sepolia.etherscan.io/tx/0xe3ad5ae7e75faac8094164b2d80b7fe41cde0a7dd5d602dc6a7a76a1a60a9576
Contract Address: 0x4A0942a3dA5075dFfCeeAa186C618e98CaC25C1c
Contract explorer link: https://sepolia.etherscan.io/address/0x4A0942a3dA5075dFfCeeAa186C618e98CaC25C1c
Executing contract deployment, TX: https://sepolia.etherscan.io/tx/0x225cfd83e352a0dd0d1c71a361aba6a8f00a982f87416b6e342b1eff13b18bf8
Contract Address: 0x28c8D90742fb2cAe6E91d71fa0761f1652F0C347
Contract explorer link: https://sepolia.etherscan.io/address/0x28c8D90742fb2cAe6E91d71fa0761f1652F0C347
Executing TX https://sepolia.etherscan.io/tx/0xbefdb61f5dbc79321f5c86b3bc6ad4286a56abbd889e193352b91348995a10ca [setting rebalancer on token pool]
TX 0xbefdb61f5dbc79321f5c86b3bc6ad4286a56abbd889e193352b91348995a10ca mined.
Block Number: 5144499
Gas Used:  46106
Block hash:  0x325d165d63fd3408ab26e33a2725e28f8bd13b7856acb1b4edd94679efcb6344
Executing contract deployment, TX: https://sepolia.etherscan.io/tx/0x978cee932339721002307e1d6c70a0e2fc067dac0b09da0809a0ce8da9115ca9
Contract Address: 0xC90b91C0b6340e9599096c2EaDCF3e4aE6b098dC
Contract explorer link: https://sepolia.etherscan.io/address/0xC90b91C0b6340e9599096c2EaDCF3e4aE6b098dC
Executing contract deployment, TX: https://sepolia.arbiscan.io//tx/0xc543cf28ff3db63df8c66922c7894bc1fc5e3be56825626344619aab6a7951c4
Contract Address: 0x59CCa90C831a16a8EE007d70498eFE2eeA0e288C
Contract explorer link: https://sepolia.arbiscan.io//address/0x59CCa90C831a16a8EE007d70498eFE2eeA0e288C
Executing contract deployment, TX: https://sepolia.arbiscan.io//tx/0xc8bb080942e85501f676cc6363d16827ab3c7bc046a7c57daa1f2d09578de6f3
Contract Address: 0xc9268244d0F9e3B526cFACf8833fDc21705E5e47
Contract explorer link: https://sepolia.arbiscan.io//address/0xc9268244d0F9e3B526cFACf8833fDc21705E5e47
Executing contract deployment, TX: https://sepolia.arbiscan.io//tx/0xdc1ae184b72a1406194b411b06345f5381eadc2ea2ad914e79cde69b15b7d731
Contract Address: 0x0C6De4B7a9E1A8e5aBCe2D9D7Fc6B4F63311Fef8
Contract explorer link: https://sepolia.arbiscan.io//address/0x0C6De4B7a9E1A8e5aBCe2D9D7Fc6B4F63311Fef8
Executing contract deployment, TX: https://sepolia.arbiscan.io//tx/0x072983d54785ca67471ecc51832d529c0aae029b3156c163b4f05817276dc905
Contract Address: 0x277EfDb0bD67a760c98c997c2716Ecab89CD62fb
Contract explorer link: https://sepolia.arbiscan.io//address/0x277EfDb0bD67a760c98c997c2716Ecab89CD62fb
Executing TX https://sepolia.arbiscan.io//tx/0xcba286924678f63fa2119be1742f7a9f20b229b99b7daac3824d6c32ee64889a [setting rebalancer on token pool]
TX 0xcba286924678f63fa2119be1742f7a9f20b229b99b7daac3824d6c32ee64889a mined.
Block Number: 9262544
Gas Used:  305495
Block hash:  0x93d99b24d15c5bd2097f4a471de96473640c19b75394c59d44042ed5c5e3424c
Executing contract deployment, TX: https://sepolia.arbiscan.io//tx/0x59c801c6df42a8d3e4f97e58f0dac1991fc721e1fdae265af7564c7ac0df3977
Contract Address: 0x6096Bde6eCc9eB9160948Ab414a35FF8aeDf1e7e
Contract explorer link: https://sepolia.arbiscan.io//address/0x6096Bde6eCc9eB9160948Ab414a35FF8aeDf1e7e
Executing TX https://sepolia.etherscan.io/tx/0x4ebab479722a74c474cb16a71849915e0d8e35fa1a73dd272e1ab1dcf61d1752 [setting cross chain rebalancer on L1 rebalancer]
TX 0x4ebab479722a74c474cb16a71849915e0d8e35fa1a73dd272e1ab1dcf61d1752 mined.
Block Number: 5144502
Gas Used:  138017
Block hash:  0x9d16226e552fc843c7628e06d87e3d2e9ef7fc2020191ef352e4696f894dbb06
Executing TX https://sepolia.arbiscan.io//tx/0xf20dc0e1501ac2482720f85052baf10caf1771fd2e969555184520ba4303aa32 [setting cross chain rebalancer on L2 rebalancer]
TX 0xf20dc0e1501ac2482720f85052baf10caf1771fd2e969555184520ba4303aa32 mined.
Block Number: 9262582
Gas Used:  622197
Block hash:  0x07a98b99fbabdf2244212c3238151f29aa5f9fd5ca5bafa4567edb0bb970a921
Deployments complete
 L1 Arm: 0xabB73a46050C14Aee721fA3a6c9333709ED806b5
 L1 Arm Proxy: 0xEeA421a17746A00C2Bc7B5869543cd2A89aC7125
 L1 Token Pool: 0x4A0942a3dA5075dFfCeeAa186C618e98CaC25C1c
 L1 Rebalancer: 0x28c8D90742fb2cAe6E91d71fa0761f1652F0C347
 L1 Bridge Adapter: 0xC90b91C0b6340e9599096c2EaDCF3e4aE6b098dC
 L2 Arm: 0x59CCa90C831a16a8EE007d70498eFE2eeA0e288C
 L2 Arm Proxy: 0xc9268244d0F9e3B526cFACf8833fDc21705E5e47
 L2 Token Pool: 0x0C6De4B7a9E1A8e5aBCe2D9D7Fc6B4F63311Fef8
 L2 Rebalancer: 0x277EfDb0bD67a760c98c997c2716Ecab89CD62fb
 L2 Bridge Adapter: 0x6096Bde6eCc9eB9160948Ab414a35FF8aeDf1e7e
```

The script will print all the addresses of the contracts at the end as well.

2. The script will then start setting up the nodes, it will migrate the database for each node (by default 4 nodes are configured, however more can be configured if needed).
3. The script will then start the chainlink app, add the job, whether its the bootstrap job or the rebalancer job, and then shut the node down.
4. In the end, if the script ran successfully, you should see something like this:

```
Contract Deployments complete
 L1 Arm: 0xe358F99d8dDDEd4Da4D8c09f1aA1E0a2FF6B04b2
 L1 Arm Proxy: 0xf86a108bF2E4244B7D1a8f265ae9A293dE8F67FD
 L1 Token Pool: 0x17219193eE340312856b10D57F854E8158c0d73D
 L1 Rebalancer: 0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e
 L1 Bridge Adapter: 0xC944b67f64870f06548FAcEA2A26994096DBEfE1
 L2 Arm: 0xbCe07F6E9495eb9c55c3E683b76A88aFb2014667
 L2 Arm Proxy: 0x4BD3F1f1ab559fb02cED41B2800762e0444a5474
 L2 Token Pool: 0xC105069c716c3cA16Eae43442Ba018a29F1CAD3C
 L2 Rebalancer: 0x4aa2525F262C663948499a4C5469a71413fEa303
 L2 Bridge Adapter: 0xf86a108bF2E4244B7D1a8f265ae9A293dE8F67FD
 Node launches complete
 OnChainPublicKeys: cba798ce8fcb1b2b5eca2583786aa3369c415189,43bb9e64771b68d594b2a7c67670d69de5a2484d,889c38708159448d2edf895af622f37ec3d75a7f,d89c486a9c133248a07f3b3c8db924e3a2f4bb93,1555cc58ae2e2c0912d0b1475226274b73a8c4f1
 OffChainPublicKeys: d8ba51dee364881009f9c856defdf68bc91645895189ffc49e1adb672fe64566,2370126da7e518ade265b0191dd79ca001c9cd29a57eca6aad47f2d27b6033d2,c168687a7221e74aedf1ca55959097907bd6c99dd4fc1990a895cf926e65013e,5356ba3a2dc3569414fc8b24c7221b8cb050ce46b46f14291645b1f508ee140f,4eaeac6001132f6728ac282d431f9860a6990f116a1a116f2e70f745cf59ac21
 ConfigPublicKeys: 8678868d1348bdc2bd278464781a1bb9adc21425ac08bc7b48ac6d6d012b363c,533fd6858fbea72330ff1d93332422ef9084a28942fee03d3df7b7880a51173d,3381b38e68c239955698840b71296ffa7aaab6b8bb646111c8056c3e4878a914,85b078ec7981f96bdf5a1f07817275de21128c893a6a335de500a2bd9e17d46b,bbf8d70cdcfb76e294ffe9e1d6284bf14cd4de1b953069d7d3a093bac1ffd77f
 PeerIDs: 12D3KooWSoQa8TiNNCxkecst65M6sDcg8r8LmQfaK3oH83t5Pusq,12D3KooWMTpnQEEqtM2GcMy4X47E2xS8ZVTg8Ddw42FZFd8fnjAp,12D3KooWHurzzRmeMg83DC4K8UtJUPcUygYpj2sSdU7PKJwoAHzV,12D3KooWDrC6QfyMjQmUUawHdxKujgYTsAcEx6RLDUhQPKNPTdVd,12D3KooWNSjjPWhMUMnB7aJYoUW9UBEqYjtHv2gdTckab8icnhns
 Transmitters L1: 0x9Ba13e50aE91B1A3889089acAC536c893d466B71,0xD605F42dC8C7c2fa84C4b88aA30D143bc6f7b471,0xa612e89970BB9879664f761D9f7245274589f399,0x3e61527f40C02FC4B08846B5F5A9290FA98f7662,0x592aa4C9315c3B6e556611Dd27D0196417588288
 Transmitters L2: 0x557219795E6dbA209DbF8476C0EA17848385a76d,0xa663265Ff6BbcA3769d2d888F1821aAB54e7C38b,0xc4Bf32b8f859ac28C99332719C47cc0F5Ebd2eA8,0xf02727ECADc3731fB5E4Bb454A7705bF320242Aa,0x155d4Aef3d24866da06406173f9AAe9BD7364241

Set config command:
 go run . set-config -l1-chain-id 11155111 -l2-chain-id 421614 -l1-liquiditymanager-address 0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e -l2-liquiditymanager-address 0x4aa2525F262C663948499a4C5469a71413fEa303 -signers cba798ce8fcb1b2b5eca2583786aa3369c415189,43bb9e64771b68d594b2a7c67670d69de5a2484d,889c38708159448d2edf895af622f37ec3d75a7f,d89c486a9c133248a07f3b3c8db924e3a2f4bb93,1555cc58ae2e2c0912d0b1475226274b73a8c4f1 -offchain-pubkeys d8ba51dee364881009f9c856defdf68bc91645895189ffc49e1adb672fe64566,2370126da7e518ade265b0191dd79ca001c9cd29a57eca6aad47f2d27b6033d2,c168687a7221e74aedf1ca55959097907bd6c99dd4fc1990a895cf926e65013e,5356ba3a2dc3569414fc8b24c7221b8cb050ce46b46f14291645b1f508ee140f,4eaeac6001132f6728ac282d431f9860a6990f116a1a116f2e70f745cf59ac21 -config-pubkeys 8678868d1348bdc2bd278464781a1bb9adc21425ac08bc7b48ac6d6d012b363c,533fd6858fbea72330ff1d93332422ef9084a28942fee03d3df7b7880a51173d,3381b38e68c239955698840b71296ffa7aaab6b8bb646111c8056c3e4878a914,85b078ec7981f96bdf5a1f07817275de21128c893a6a335de500a2bd9e17d46b,bbf8d70cdcfb76e294ffe9e1d6284bf14cd4de1b953069d7d3a093bac1ffd77f -peer-ids 12D3KooWSoQa8TiNNCxkecst65M6sDcg8r8LmQfaK3oH83t5Pusq,12D3KooWMTpnQEEqtM2GcMy4X47E2xS8ZVTg8Ddw42FZFd8fnjAp,12D3KooWHurzzRmeMg83DC4K8UtJUPcUygYpj2sSdU7PKJwoAHzV,12D3KooWDrC6QfyMjQmUUawHdxKujgYTsAcEx6RLDUhQPKNPTdVd,12D3KooWNSjjPWhMUMnB7aJYoUW9UBEqYjtHv2gdTckab8icnhns -l1-transmitters 0x9Ba13e50aE91B1A3889089acAC536c893d466B71,0xD605F42dC8C7c2fa84C4b88aA30D143bc6f7b471,0xa612e89970BB9879664f761D9f7245274589f399,0x3e61527f40C02FC4B08846B5F5A9290FA98f7662,0x592aa4C9315c3B6e556611Dd27D0196417588288 -l2-transmitters 0x557219795E6dbA209DbF8476C0EA17848385a76d,0xa663265Ff6BbcA3769d2d888F1821aAB54e7C38b,0xc4Bf32b8f859ac28C99332719C47cc0F5Ebd2eA8,0xf02727ECADc3731fB5E4Bb454A7705bF320242Aa,0x155d4Aef3d24866da06406173f9AAe9BD7364241

Funding command:
 go run . fund-contracts -l1-chain-id 11155111 -l2-chain-id 421614 -l1-liquiditymanager-address 0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e -l2-liquiditymanager-address 0x4aa2525F262C663948499a4C5469a71413fEa303 -l1-token-address 0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9 -l2-token-address 0x980B62Da83eFf3D4576C647993b0c1D7faf17c73 -l1-token-pool-address 0x17219193eE340312856b10D57F854E8158c0d73D -l2-token-pool-address 0xC105069c716c3cA16Eae43442Ba018a29F1CAD3C
```

Before we are able to launch the nodes, we need to create their config TOML files. Here’s a sample for Arbitrum Sepolia + Sepolia:

```toml
# Arbitrum Sepolia
[[EVM]]
ChainID = "421614"
GasEstimator.LimitDefault = 3_500_000
# This is needed so that the log poller doesn't start
# from a super old block. It takes a while to catch up to
# the tip.
FinalityTagEnabled = false

[[EVM.Nodes]]
HTTPURL = "<http-rpc>"
Name = "arbitrum_sepolia_1"
WSURL = "<ws-rpc>"

# Sepolia
[[EVM]]
ChainID = "11155111"
GasEstimator.LimitDefault = 3_500_000

[[EVM.Nodes]]
HTTPURL = "<http-rpc>"
Name = "sepolia_1"
WSURL = "<ws-rpc>"

[EVM.Transactions]
ForwardersEnabled = false

[Feature]
LogPoller = true

[OCR2]
Enabled = true
ContractPollInterval = "15s"

[OCR]
Enabled = false

[P2P.V2]
# Different per node
ListenAddresses = ["127.0.0.1:8001"]

[WebServer]
# Different per node
# Adjust accordingly in your local setup
HTTPPort = 6689

[WebServer.TLS]
CertPath = ''
ForceRedirect = false
Host = ''
HTTPSPort = 0
KeyPath = ''
```

You need a separate config file for each node. The main differences between each node’s configuration will be the `HTTPPort` field, which cannot be the same for all nodes since they will be trying to listen on the same port, and the `P2P.V2.ListenAddresses` field for a similar reason. The secrets.toml file will also be different due to different databases being used for each node.

Create 4 secrets files (or however many you need, you need 1 for each node). I named them `secrets_0.toml` up until `secrets_3.toml` (for 4 nodes). Here are the contents:

```toml
[Database]
# notice rebalancer-test-0 as the database name
# set this to be the correct name for each node
# databases are named rebalancer-test-i where i is the node
# index, starting from 0. So for 4 nodes, i will go from 0 to 3.
URL = "postgres://postgres:postgres_password_padded_for_security@localhost:5432/rebalancer-test-0?sslmode=disable"

[Password]
# This is the test password in tools/secrets/password.txt
Keystore = 'T.tLHkcmwePT/p,]sYuntjwHKAsrhm#4eRs4LuKHwvHejWYAC2JP4M8HimwgmbaZ'
```

One more thing to do before we run the nodes is to fund the relevant contracts in the setup. These contracts are the L1 rebalancer, L1 token pool, and L2 token pool. The L2 rebalancer doesn’t need to be funded (usually) because L2 → L1 transfers are typically free and require no native except for gas paid. In the same shell you used to run the setup script from, run the `fund-contracts` command:

```bash
go run . fund-contracts -l1-chain-id 11155111 -l2-chain-id 421614 -l1-liquiditymanager-address 0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e -l2-liquiditymanager-address 0x4aa2525F262C663948499a4C5469a71413fEa303 -l1-token-address 0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9 -l2-token-address 0x980B62Da83eFf3D4576C647993b0c1D7faf17c73 -l1-token-pool-address 0x17219193eE340312856b10D57F854E8158c0d73D -l2-token-pool-address 0xC105069c716c3cA16Eae43442Ba018a29F1CAD3C
```

Note that your command will probably have different rebalancer addresses (and maybe different token addresses if you’re not using WETH).

At this point, with the above files saved, you can create 4 new shell sessions. In each session:

```bash
cd <ccip-repo>
# build the ./chainlink binary
make chainlink-dev
# run the node
./chainlink -s  <path-to-secrets_<i>.toml> -c <path-to-config.toml> node start --api tools/secrets/apicredentials
```

Make sure to reference the correct secrets file.

When the nodes start running they will log that they config digest is 0 onchain, which is correct. You have to set the configuration onchain using the command that is printed out at the end of the script, e.g:

```bash
go run . set-config -l1-chain-id 11155111 -l2-chain-id 421614 -l1-liquiditymanager-address 0x2CBfee4A397Dc75F47b10E89b2e435a3eC940073 -l2-liquiditymanager-address 0xdE9497B46f2B046A852ec582fc629c8279416Eb6 -signers 2db30dd13f2dc2f58c9a20009e695c03a68ee839,897b2db516417818d75adfa976049487526f955d,6acd64a139fd8423977f2ab0a08f09e2fa8e2268,49827c809ab89e5f0ece2ba0e9ffcca8d3633965,d7080d0c94fdcbd287465f6ce9a502d4b2412239 -offchain-pubkeys 96d151b49557ed57e4a6b968f1f4bfe1c0a8114e76d68502baaccf5bec9c8e8a,95c76eed08139c3f3bcb28e5f703b44b2559f6d90950f984c6ebd773c73ab73c,eb9b63099c267c89175beda8c451263db7997c718f7f70ba8d9a8427ddd1841d,4134611fa18bccff167d959da84954f6fca560e28219545769206dcb10edcedb,004b72c8283e6137328e38a79f719a0c488bc8e77ac637f5f0140fa406b57561 -config-pubkeys 132119a037a37d80e1307ac3d960c77bb128fa27ec116397ac755714e4d7850f,ec56836dc5371d14ea9afbfc0a16b08d5bf9f8d36b3007999aa95cafc998370f,a6166fd070f7fc299abdaebc9439bb9dbaefe444a5b26ba9e22e569e9aeeca7a,af677a193c3c2a36ca55df4a0e45ab8f47fe97947b8f82bdda0aec4b42b3bd2e,5ad7492f0cc0f2043631f874b693a26f2cb629bee89b4010f3bd185596ee1205 -peer-ids 12D3KooWArASehuA7bXs4Y2Gk9uMsHqW28wzHv8iojy5XpSWUk76,12D3KooWJ5uW1SFF6v8HPMhja6SACtG5dcwiYc9Zms88DqGT5jrd,12D3KooWJfftNZLA8aPLXinzuJym16EhecS3z4WMoDL9xbtLr2YE,12D3KooWLvYmwR7LEbvWDe6sxeZfiFWSVuFEsPLzz45HLzgnn6cg,12D3KooWExWYTuUrR2zf35s9jB113TXe7SmUgeAtrticH8mSBX8w -l1-transmitters 0x20E323185F0116F3930f272cC70127d80D15ACC0,0x4b7fF6592dA98f5162656A74704F2378b0E37448,0x6612BDC4a6d70E59a7DD457F8B4df0aa3Ddc1e5D,0x93f0a01Cf89032Cc645863366555fDa9467dFF26,0x73Ba1041e35528e0cA92cCBc71125086a3EDa635 -l2-transmitters 0xeDdaBfa6DCFFCF8955d6FB3cd2648039FAa71A16,0xaFE78225E20c8e210Ea737455E05ba6035C21B5F,0x8258C5956dd1E20f978fF95B979191D19700fa1b,0x2c88db3162E0a34C5de4837831f85CE8B503Eb7B,0x27bB8e111785aaBf0013F0aeA1ADBE7eaE378965
```

Once you run this command, the nodes will pick up the `ConfigSet` event and discover each other. Note that the regular (non-bootstrap) nodes will have to pick up both the set config event on L2 and on L1.

Seeing the logs below:

```
2024-01-25T17:35:42.383+0200 [INFO]  TrackConfig: returning config                      managed/track_config.go:64       configDigest=0001ba1c68e44e9e2e5603040e9ff394e81cac620d6bde5b8f8aaf0a570996ad contractID=0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e feedID=<nil> jobID=1 jobName=bootstrap-chainID-11155111 logger=OCRBootstrap version=2.7.1-ccip1.2.1-beta.0@37a8ec8
2024-01-25T17:35:42.383+0200 [INFO]  runWithContractConfig: switching between configs   managed/run_with_contract_config.go:73 contractID=0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e feedID=<nil> jobID=1 jobName=bootstrap-chainID-11155111 logger=OCRBootstrap newConfigDigest=0001ba1c68e44e9e2e5603040e9ff394e81cac620d6bde5b8f8aaf0a570996ad oldConfigDigest=0000000000000000000000000000000000000000000000000000000000000000 version=2.7.1-ccip1.2.1-beta.0@37a8ec8
2024-01-25T17:35:42.383+0200 [INFO]  runWithContractConfig: winding down old configuration managed/run_with_contract_config.go:114 contractID=0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e feedID=<nil> jobID=1 jobName=bootstrap-chainID-11155111 logger=OCRBootstrap newConfigDigest=0001ba1c68e44e9e2e5603040e9ff394e81cac620d6bde5b8f8aaf0a570996ad oldConfigDigest=0000000000000000000000000000000000000000000000000000000000000000 version=2.7.1-ccip1.2.1-beta.0@37a8ec8
2024-01-25T17:35:42.383+0200 [INFO]  runWithContractConfig: closed old configuration    managed/run_with_contract_config.go:120 contractID=0x1f3210Ade3b167c6C83aa3d5972aFD363b8abb2e feedID=<nil> jobID=1 jobName=bootstrap-chainID-11155111 logger=OCRBootstrap newConfigDigest=0001ba1c68e44e9e2e5603040e9ff394e81cac620d6bde5b8f8aaf0a570996ad oldConfigDigest=0000000000000000000000000000000000000000000000000000000000000000 version=2.7.1-ccip1.2.1-beta.0@37a8ec8
```

Indicates that the node has picked up the new configuration, and will switch to the new configuration.

Shortly thereafter, you should start seeing transmissions.

## Running Bridge Transfers Through the Adapter Contracts

### Optimism

First, you have to deploy the L1 and L2 bridge adapters:

```shell
# Switch into the liquiditymanager scripts dir
cd core/scripts/ccip/liquiditymanager
# Uses sepolia chain id, switch to 1 for mainnet, and sepolia weth contract address (same as arbitrum)
go run . deploy-op-l1-adapter -l1-chain-id 11155111
# Uses sepolia chain id, switch to 420 for OP mainnet
go run . deploy-op-l2-adapter -l2-chain-id 11155420
```

Now you're ready to do some cross-chain transfers.

#### L1 -> L2

In order to send tokens from L1 to L2, pick a token ([FaucetTestingToken](https://sepolia.etherscan.io/address/0x5589bb8228c07c4e15558875faf2b859f678d129) is easiest).

> NOTE: if you're using the [FaucetTestingToken](https://sepolia.etherscan.io/address/0x5589bb8228c07c4e15558875faf2b859f678d129) you can get some tokens by interacting
with the [faucet](https://sepolia.etherscan.io/address/0x5589bb8228c07c4e15558875faf2b859f678d129#writeContract#F2) method using MetaMask.

```shell
# Uses sepolia chain id, switch to 1 for mainnet
# All values are in the lowest denomination (i.e wei)
go run . op-send-to-l2 -l1-chain-id $SEPOLIA_CHAIN_ID \
    -l1-bridge-adapter-address <l1-adapter-address> \
    -l2-to-address <to-address-on-L2> \
    -l1-token-address 0x5589BB8228C07c4e15558875fAf2B859f678d129 \ # This is the L1 FaucetTestingToken address
    -l2-token-address 0xD08a2917653d4E460893203471f0000826fb4034 \ # This is the L2 FaucetTestingToken address
    -amount 1
```

The easiest way to see the status of the transaction is to visit the [tokentxns](https://sepolia-optimism.etherscan.io/address/0x77ffC73eD3B2614D21B3398fe368E989f318b412#tokentxns)
page on Optimism Sepolia's Etherscan site (replace with your address) to see if the coins have been minted on L2.

#### L2 -> L1

In order to withdraw from L2 to L1, pick a token ([FaucetTestingToken](https://sepolia.etherscan.io/address/0x5589bb8228c07c4e15558875faf2b859f678d129) is easiest)
and make sure you have enough balance.

Then invoke the following function:

```shell
# Uses sepolia chain id, switch to 420 for mainnet
# All values are in the lowest denomination (i.e wei)
go run . op-withdraw-from-l2 -l2-chain-id 11155420 \
    -l2-bridge-adapter-address <l2-adapter-address> \
    -amount 1 \
    -l1-to-address <l1-to-address> \
    -l2-token-address 0xD08a2917653d4E460893203471f0000826fb4034 # This is the L2 FaucetTestingToken
```

Once the `op-withdraw-from-l2` executes successfully, you have to wait a few minutes (usually around 10) to prove the withdrawal on L1.

#### Prove Withdrawal on L1

In order to prove the withdrawal on L1, make sure you have the L2 Withdrawal Transaction Hash. This is printed by the `op-withdraw-from-l2` command if you used that to initiate a withdrawal.

In order to be sure that the withdrawal can be proven, you can visit the transaction hash on the explorer, e.g [0x1569872053c27de1cffd11dc1951c49e03a61dba8131afe947c6bc5abe352c20](https://sepolia-optimism.etherscan.io/tx/0x1569872053c27de1cffd11dc1951c49e03a61dba8131afe947c6bc5abe352c20). If you see "L1 State Batch Index:" then that means the L2 batch that includes the withdrawal has been submitted to L1 and can be proven to be
part of that batch.

```shell
go run . op-prove-withdrawal-l1 -l1-chain-id 11155111 \
    -l2-chain-id 11155420 \
    -l2-tx-hash <withdrawal-tx-hash> \
    -l1-bridge-adapter-address <l1-bridge-adapter>
```

`op-prove-withdrawal-l1` will print some important debugging information and then send the withdrawal proof transaction to the L1 OptimismPortal contract.

#### Finalize Withdrawal on L1

In order to finalize a withdrawal from L2, you must first prove the withdrawal (see above) and wait a sufficient amount of time (usually 10+ minutes).

The finalize command is as follows:

```shell
go run . op-finalize-l1 -l1-chain-id 11155111 \
    -l2-chain-id 11155420 \
    -l2-tx-hash <withdrawal-tx-hash> \
    -l1-bridge-adapter-address <l1-bridge-adapter>
```

Once this transaction is mined you should have your funds back on L1. See the transaction's logs for more details or query your balance in the appropriate token contract.

### Arbitrum

First, you have to deploy the L1 and L2 bridge adapters:

```shell
# Switch into the liquiditymanager scripts dir
cd core/scripts/ccip/liquiditymanager
# Uses sepolia chain id, switch to 1 for mainnet
go run . deploy-arb-l1-adapter -l1-chain-id 11155111
# Uses sepolia chain id, switch to 42161 for arb mainnet
go run . deploy-arb-l2-adapter -l2-chain-id 421614
```

Now you're ready to do some cross-chain transfers.

#### L1 -> L2

In order to send tokens from L1 to L2, pick a token (WETH is easiest) and make sure you have enough balance (the script will check).
Then invoke the following function:

```shell
# Uses sepolia chain id, switch to 42161 for arb mainnet
# All values are in the lowest denomination (i.e wei)
go run . arb-send-to-l2 -l1-chain-id 11155111 -l2-chain-id 421614 \
    -l1-bridge-adapter-address <l1-adapter-address> \
    -amount 1 -l2-to-address <to-address-on-L2> \                # This is the address that will receive the funds - use an EOA you own to easily check status on Arbitrum's Bridge UI
    -l1-token-address 0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9 # This is the L1 WETH token
```

You can go to the [Arbitrum Bridge](https://bridge.arbitrum.io) and connect with the receiver address on L2 (via MetaMask) to see the status of the transaction. It should
be automatically deposited on L2 in around 10 minutes.

#### L2 -> L1

In order to withdraw from L2 to L1, pick a token (WETH is probably the easiest) and make sure you have enough balance.
Then invoke the following function:

```shell
# Uses sepolia chain id, switch to 42161 for arb mainnet
# All values are in the lowest denomination (i.e wei)
go run . arb-withdraw-from-l2 -l2-chain-id 421614 -l2-bridge-adapter-address <adapter-address-from-deploy-cmd> \
    -amount 1 -l1-to-address <receiver-address-on-L1> \
    -l2-token-address 0x980B62Da83eFf3D4576C647993b0c1D7faf17c73 \ # This is the L2 WETH token
    -l1-token-address 0x7b79995e5f793A07Bc00c21412e50Ecae098E7f9   # This is the L1 WETH token
```

If you're using WETH, you can get some WETH from native this way:

```shell
# Uses sepolia chain id, switch to 42161 for arb mainnet
# All values are in the lowest denomination (i.e wei)
go run . deposit-weth -amount 1 -weth-address 0x980B62Da83eFf3D4576C647993b0c1D7faf17c73 -chain-id 421614
```

Once the `arb-withdraw-from-l2` command executes successfully, you will have to wait some time until you can claim the funds on L1 (see next section).

#### Finalize Withdrawal on L1

In order to finalize an L2 withdrawal on L1, you will need to first have a transaction on L2 that is not yet finalized.

See the previous section for doing a withdrawal and then come back when you have a successful transaction on L2.

Save that transaction hash and use it in the following command:

```shell
# Uses sepolia and arb sepolia!
go run . arb-finalize-l1 -l1-chain-id 11155111 -l2-chain-id 421614 -l2-tx-hash <tx-hash-here> -l1-bridge-adapter-address <l1-adapter-address>
```

This will print some execution information and then finalize the transaction.

[!NOTE] Finalization WILL FAIL if the batch that the withdrawal tx was in was not successfully submitted to L1! This command can be used as an indicator if a tx is ready to finalize.
