## Oracle contract

Before deploying your own Oracle contract, make sure you have gone through the instructions to set up your [Chainlink node on Ropsten](./RopstenNode.md).

- In Remix, import the contracts at `chainlink/examples/ropsten/contracts`
- Click on the `Oracle.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/10-29-38.png)

- Change to the Run tab
- Select Oracle from the dropdown in the right panel
- Copy and paste the line below and enter it into the text field next to the Create button <br>
    **0x20fE562d797A42Dcb3399062AE9546cd06f63280**
- Click Create

![create](./images/10-31-04.png)

- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 20 if you don't know what to pick)
- Select Submit

![deploy contracts](./images/07-24-30.png)

- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction

![confirm contract deploy](./images/10-54-23.png)

- Once successful, you should have a new address for the deployed contract

![contract deploy successful](./images/07-25-49.png)

- Keep note of the Oracle contract's address, you will need it for adding a JobSpec to the node.

- In MyCrypto/Mew, go to the Contracts tab
- Paste the deployed Oracle contract address and the [OracleABI](./OracleABI.json), and click Access

![oracle abi](./images/10-59-27.png)

- Select the `transferOwnership` method and paste in the address of your Chainlink node

You can get the address of your node when you start it with `chainlink node`. There will be an `[INFO]` line displayed similar to the one below:

```
2018-05-07T16:01:24Z [INFO]  ETH Balance for 0x5958C587503b40A8576998dB56C3F5ec1f024C3D: 1.962567814000000000 cmd/client.go:71 
```

![transfer owner](./images/11-02-03.png)

- Access your wallet with MetaMask and click Write
- Generate the transaction

![generate transfer tx](./images/11-03-38.png)

- Submit with MetaMask

![send transfer tx](./images/11-04-00.png)

Keep note of your Oracle contract's address, you can now use this address for any of the [Advanced Deployment](./AdvancedDeploy.md) instructions.