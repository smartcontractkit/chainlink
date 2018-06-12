# Advanced Contract Deployment

This section of the guide uses [Remix](https://remix.ethereum.org) in order to compile and deploy an example contract from the `examples/ropsten/contracts` directory.

In Remix, import the `RopstenConsumer.sol` contract at `chainlink/examples/ropsten/contracts`

![contracts](./images/12-29-32.png)

- Click on the `RopstenConsumer.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/12-36-11.png)

- Change to the Run tab
- RopstenConsumer should already be selected
- Click Deploy

![deploy1](./images/12-37-18.png)

- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 20 if you don't know what to pick)
- Select Submit

![deploy contracts](./images/07-24-30.png)

- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction

![confirm contract deploy](./images/07-25-22.png)

- Once successful, you should have a new address for the deployed contract

![contract deploy successful](./images/07-25-49.png)

*You can now reference the [sending Ropsten LINK to the Consumer contract](./README.md#send-ropsten-link-to-the-consumer-contract) section to fund the contract.*

- In Remix, you can interact with and call the requesting functions directly, by supplying a string for the methods that begin with "request".

![contract functions](./images/12-50-55.png)