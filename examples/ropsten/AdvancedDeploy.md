# Advanced Contract Deployment

This section of the guide uses [Remix](https://remix.ethereum.org) in order to compile and deploy contracts in the `examples/ropsten/contracts` directory.

In Remix, import the contracts at `chainlink/examples/ropsten/contracts`

![contracts](./images/10-12-44.png)

You can now skip to the following sections:

- [ConsumerUint256 contract](#consumeruint256-contract)
- [ConsumerInt256 contract](#consumerint256-contract)
- [ConsumerBytes32 contract](#consumerbytes32-contract)

## ConsumerUint256 Contract

- Click on the `ConsumerUint256.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/10-22-54.png)

- Change to the Run tab
- Select ConsumerUint256 from the dropdown in the right panel
- For a standard Uint256 type, copy and paste the line below and enter it into the text field next to the Create button <br>
    **"0x20fE562d797A42Dcb3399062AE9546cd06f63280", "0x5be84B6381d45579Ed04A887B8473F76699E0389", "0x3565616462613666303737653463613861633838636333373065623636613566"**
- For the value to be multiplied by 100 in the node before writing to the blockchain, copy and paste the line below and enter it into the text field next to the Create button <br>
    **"0x20fE562d797A42Dcb3399062AE9546cd06f63280", "0x5be84B6381d45579Ed04A887B8473F76699E0389", "0x6434316130626466393638613433616361383832326366383161326331666137"**
- Click Create

![create](./images/10-24-56.png)

Go to [Remix to Metamask Contract Deployment](#remix-to-metamask-contract-deployment)

## ConsumerInt256 contract

- Click on the `ConsumerInt256.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/10-38-56.png)

- Change to the Run tab
- Select ConsumerInt256 from the dropdown in the right panel
- Copy and paste the line below and enter it into the text field next to the Create button <br>
    **"0x20fE562d797A42Dcb3399062AE9546cd06f63280", "0x5be84B6381d45579Ed04A887B8473F76699E0389", "0xDoes-Not-Exist-Yet-Fill-Me-In"**
- Click Create

![create](./images/10-39-57.png)

Go to [Remix to Metamask Contract Deployment](#remix-to-metamask-contract-deployment)

## ConsumerBytes32 contract

- Click on the `ConsumerBytes32.sol` contract in the left side-bar
- On the Compile tab, click on the "Start to compile" button near the top-right

![compile](./images/10-40-59.png)

- Change to the Run tab
- Select ConsumerBytes32 from the dropdown in the right panel
- Copy and paste the line below and enter it into the text field next to the Create button <br>
    **"0x20fE562d797A42Dcb3399062AE9546cd06f63280", "0x5be84B6381d45579Ed04A887B8473F76699E0389", "0x3236363037363337303734333466333562623538613964623333656531383033"**
- Click Create

![create](./images/10-42-20.png)

Go to [Remix to Metamask Contract Deployment](#remix-to-metamask-contract-deployment)

## Remix to Metamask Contract Deployment

- Metamask will prompt you to Confirm the Transaction
- You will need to choose a Gas Price (use 20 if you don't know what to pick)
- Select Submit

![deploy contracts](./images/07-24-30.png)

- A link to Etherscan will display at the bottom, you can open that in a new tab to keep track of the transaction

![confirm contract deploy](./images/07-25-22.png)

- Once successful, you should have a new address for the deployed contract

![contract deploy successful](./images/07-25-49.png)

You can now go back to [sending Ropsten LINK to the Consumer contract](./README.md#send-ropsten-link-to-the-consumer-contract).