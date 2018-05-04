# Chainlink Ropsten Instructions

This guide will allow you to create and deploy a consuming contract to fulfill a data request using our deployed oracle contract. You do not need to run a node yourself in order to follow these instructions.

## Tools

This guide requires the following tools:

- [Metamask](https://metamask.io/)
- [Remix](https://remix.ethereum.org)

## Setup

Add the Ropsten LINK token to Metamask:

- Switch to the "Ropsten Test Net" network in Metamask
- Click on the Tokens tab
- Click Add Token button
- Paste the contract address 0x20fe562d797a42dcb3399062ae9546cd06f63280
- The rest should fill in, if it doesn't the Token Symbol is LINK and use 18 for Decimals

You should now see your Ropsten LINK

### Faucets

Ropsten ETH
- http://faucet.ropsten.be:3001/
- https://faucet.metamask.io/

Ropsten LINK
- Let the team know your Ropsten Ethereum address on [Gitter](https://gitter.im/smartcontractkit-chainlink/Lobby) and Ropsten LINK will be sent to you.

## Compile Your Consuming Contract

- Update your local repository from [Chainlink](https://github.com/smartcontractkit/chainlink) or [download](https://github.com/smartcontractkit/chainlink/archive/master.zip) a zip.
- In Remix, import the contracts at `chainlink/examples/ropsten/contracts`
- Click on the `Consumer.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right
- Change to the Run tab
- Select Consumer from the dropdown in the right panel
- Copy and paste the line below and enter it into the text field next to the Create button <br>
    <mark>"0x20fE562d797A42Dcb3399062AE9546cd06f63280", "0x4d40982F8408e496F3dEEfE72550F23680013872", "2e7a2bb478374fbd9542cbb7f5f30fa5"</mark>
- Click Create
- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 20 if you don't know what to pick)
- Select Submit
- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction
- Once successful, you should have a new address for the Consumer contract

## Send Ropsten LINK to the Consumer Contract

Now that your Consumer contract is deployed to Ropsten, you need to send some Ropsten LINK to it.

- Open your favorite wallet (MEW, MyCrypto, etc.) and connect to the Ropsten network
- Go to the Send tab in MEW or MyCrypto
- Access your wallet using Metamask
- You may need to add the Ropsten LINK token to the wallet so that it recognizes it
  - Contract address: 0x20fe562d797a42dcb3399062ae9546cd06f63280
  - Token Symbol: LINK
  - Decimals: 18
- Send LINK to the deployed address of your Consumer contract (1 LINK is enough)

## Call the Consumer Contract to Request Data from Chainlink

The Consumer contract should now have some Ropsten LINK on it. Now you can call it to make a request on the network. The examples below use functionallity from MyEtherWallet & MyCrypto.

- Go to the Contracts tab
- Paste your deployed Consumer contract address
- You can get the API / JSON Interface from the [`ConsumerABI.json`](./ConsumerABI.json) file or from Remix:
  - In Remix, go to the Compile tab
  - Select the Consumer contract from the drop-down
  - Click the Details button
  - The 4th section should be ABI, click the "Copy value to clipboard" icon
- Paste the ABI and click the Access button
- A new section appears labeled Read / Write Contract
- Click the Select a function drop-down and choose `requestEthereumPrice`
- The values accepted for the `_currency` field are `USD`, `EUR`, and `JPY`.
- Access your wallet using Metamask again, and send the transaction leaving the Amount to Send as 0, and Gas Limit as the default.

## Verify Data was Written

(You may need to refresh MEW/MyCrypto first)

Back on the Contracts tab of MyEtherWallet or MyCrypto:

- Enter your Consumer contract address and the ABI from Remix
- Select the `currentPrice` function
- The encoded hex value should be displayed
- You can go [here](https://adibas03.github.io/online-ethereum-abi-encoder-decoder/#/decode) to convert that to a readable value
  - Select "decode" for the Action
  - Enter `uint256` for the Argument Type
  - Paste the hex value into the Encoded data field and click the `DECODE` button
  - The decoded value will display below
