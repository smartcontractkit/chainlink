# Chainlink Ropsten Instructions

## Tools

- [Metamask](https://metamask.io/)
- [Remix](https://remix.ethereum.org)

## Setup

Add the Ropsten LINK token to Metamask
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
- http...

## Compile Your Consuming Contract

- Update your local repository from [Chainlink](https://github.com/smartcontractkit/chainlink)
- In Remix, import the contracts at `$GOPATH/src/github.com/smartcontractkit/chainlink/examples/ropsten/contracts`
- Click on the `Consumer.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right
- Change to the Run tab
- Select the Consumer contract on the right side-bar
- Copy and paste the line below and enter it into the text field next to the Create button <br>
    "0x20fe562d797a42dcb3399062ae9546cd06f63280", "0x3d4B58A86a0Ee06A99CFCD7AB8abb3d0d1458C9a"
- Click Create
- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 40 if you don't know what to pick)
- Select Submit
- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction
- Once successful, you should have a new address for the Consumer contract

## Send Ropsten LINK to the Consumer Contract

Now that your Consumer contract is deployed to Ropsten, you need to send some Ropsten LINK to it.

- Open your favorite wallet (MEW, MyCrypto, etc.) and connect to the Ropsten network
- You may need to add the Ropsten LINK token to the wallet so that it recognizes it
  - Contract address: 0x20fe562d797a42dcb3399062ae9546cd06f63280
  - Token Symbol: LINK
  - Decimals: 18
- Send LINK to the deployed address of your Consumer contract

## Call the Consumer Contract to Request Data from Chainlink

The Consumer contract should now have some Ropsten LINK on it. Now you can call it to make a request on the network. The examples below use functionallity from MyEtherWallet & MyCrypto.

- Go to the Contracts tab
- Paste your deployed Consumer contract address
- You can get the API / JSON Interface from Remix:
  - In Remix, go to the Compile tab
  - Select the Consumer contract from the drop-down
  - Click the Details button
  - The 4th section should be ABI, click the "Copy value to clipboard" icon
- Paste the ABI and click the Access button
- A new section appears labeled Read / Write Contract
- Click the Select a function drop-down and choose `requestEthereumPrice`
- If not running your own node with a Ropsten Ethereum node, you can enter `something` into the `_jobid` field
- The values accepted for the `_currency` field are `USD`, `EUR`, and `JPY`.

## Verify Data was Written

Back on the Contracts tab of MyEtherWallet or MyCrypto:

- Enter your Consumer contract address and the ABI from Remix
- Select the `currentPrice` function
- The value in `bytes32` should be displayed
- You can go [here](https://www.rapidtables.com/convert/number/ascii-hex-bin-dec-converter.html) to convert that to a readable value

## (Optional) Run your own node against the Ropsten network