import { ethers } from "hardhat";
import { toBytes32String, toWei } from "../test-helpers/helpers";
import { Contract, ContractFactory } from "ethers";
import { getUsers, Roles } from "../test-helpers/setup";
import { convertFufillParams, convertFulfill2Params, decodeRunRequest } from "../test-helpers/oracle";
import { gasDiffLessThan } from "../test-helpers/matchers";

let operatorFactory: ContractFactory;
let oracleFactory: ContractFactory;
let basicConsumerFactory: ContractFactory;
let linkTokenFactory: ContractFactory;

let roles: Roles;

before(async () => {
  const users = await getUsers();

  roles = users.roles;
  operatorFactory = await ethers.getContractFactory("Operator", roles.defaultAccount);
  oracleFactory = await ethers.getContractFactory("src/v0.6/Oracle.sol:Oracle", roles.defaultAccount);
  basicConsumerFactory = await ethers.getContractFactory(
    "src/v0.6/tests/BasicConsumer.sol:BasicConsumer",
    roles.defaultAccount,
  );
  linkTokenFactory = await ethers.getContractFactory("LinkToken", roles.defaultAccount);
});

describe("Operator Gas Tests", () => {
  const specId = "0x4c7b7ffb66b344fbaa64995af81e355a00000000000000000000000000000000";
  let link: Contract;
  let oracle1: Contract;
  let operator1: Contract;
  let operator2: Contract;

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy();

    operator1 = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, await roles.defaultAccount.getAddress());
    await operator1.setAuthorizedSenders([await roles.oracleNode.getAddress()]);

    operator2 = await operatorFactory
      .connect(roles.defaultAccount)
      .deploy(link.address, await roles.defaultAccount.getAddress());
    await operator2.setAuthorizedSenders([await roles.oracleNode.getAddress()]);

    oracle1 = await oracleFactory.connect(roles.defaultAccount).deploy(link.address);
    await oracle1.setFulfillmentPermission(await roles.oracleNode.getAddress(), true);
  });

  // Test Oracle.fulfillOracleRequest vs Operator.fulfillOracleRequest
  describe("v0.6/Oracle vs v0.7/Operator #fulfillOracleRequest", () => {
    const response = "Hi Mom!";
    let basicConsumer1: Contract;
    let basicConsumer2: Contract;

    let request1: ReturnType<typeof decodeRunRequest>;
    let request2: ReturnType<typeof decodeRunRequest>;

    beforeEach(async () => {
      basicConsumer1 = await basicConsumerFactory.connect(roles.consumer).deploy(link.address, oracle1.address, specId);
      basicConsumer2 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator1.address, specId);

      const paymentAmount = toWei("1");
      const currency = "USD";

      await link.transfer(basicConsumer1.address, paymentAmount);
      const tx1 = await basicConsumer1.requestEthereumPrice(currency, paymentAmount);
      const receipt1 = await tx1.wait();
      request1 = decodeRunRequest(receipt1.logs?.[3]);

      await link.transfer(basicConsumer2.address, paymentAmount);
      const tx2 = await basicConsumer2.requestEthereumPrice(currency, paymentAmount);
      const receipt2 = await tx2.wait();
      request2 = decodeRunRequest(receipt2.logs?.[3]);
    });

    it("uses acceptable gas", async () => {
      const tx1 = await oracle1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...convertFufillParams(request1, response));
      const tx2 = await operator1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...convertFufillParams(request2, response));
      const receipt1 = await tx1.wait();
      const receipt2 = await tx2.wait();
      // 38014 vs 40260
      gasDiffLessThan(2500, receipt1, receipt2);
    });
  });

  // Test Operator1.fulfillOracleRequest vs Operator2.fulfillOracleRequest2
  // with single word response
  describe("Operator #fulfillOracleRequest vs #fulfillOracleRequest2", () => {
    const response = "Hi Mom!";
    let basicConsumer1: Contract;
    let basicConsumer2: Contract;

    let request1: ReturnType<typeof decodeRunRequest>;
    let request2: ReturnType<typeof decodeRunRequest>;

    beforeEach(async () => {
      basicConsumer1 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator1.address, specId);
      basicConsumer2 = await basicConsumerFactory
        .connect(roles.consumer)
        .deploy(link.address, operator2.address, specId);

      const paymentAmount = toWei("1");
      const currency = "USD";

      await link.transfer(basicConsumer1.address, paymentAmount);
      const tx1 = await basicConsumer1.requestEthereumPrice(currency, paymentAmount);
      const receipt1 = await tx1.wait();
      request1 = decodeRunRequest(receipt1.logs?.[3]);

      await link.transfer(basicConsumer2.address, paymentAmount);
      const tx2 = await basicConsumer2.requestEthereumPrice(currency, paymentAmount);
      const receipt2 = await tx2.wait();
      request2 = decodeRunRequest(receipt2.logs?.[3]);
    });

    it("uses acceptable gas", async () => {
      const tx1 = await operator1
        .connect(roles.oracleNode)
        .fulfillOracleRequest(...convertFufillParams(request1, response));

      const responseTypes = ["bytes32"];
      const responseValues = [toBytes32String(response)];
      const tx2 = await operator2
        .connect(roles.oracleNode)
        .fulfillOracleRequest2(...convertFulfill2Params(request2, responseTypes, responseValues));

      const receipt1 = await tx1.wait();
      const receipt2 = await tx2.wait();
      gasDiffLessThan(1240, receipt1, receipt2);
    });
  });
});
