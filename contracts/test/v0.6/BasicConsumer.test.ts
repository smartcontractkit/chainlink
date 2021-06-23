import { ethers } from "hardhat";
import { toWei, increaseTime5Minutes, toHex } from "../test-helpers/helpers";
import { assert, expect } from "chai";
import { BigNumber, constants, Contract, ContractFactory } from "ethers";
import { Roles, getUsers } from "../test-helpers/setup";
import { bigNumEquals, evmRevert } from "../test-helpers/matchers";
import { convertFufillParams, decodeRunRequest, encodeOracleRequest, RunRequest } from "../test-helpers/oracle";
import cbor from "cbor";
import { makeDebug } from "../test-helpers/debug";

const d = makeDebug("BasicConsumer");
let basicConsumerFactory: ContractFactory;
let oracleFactory: ContractFactory;
let linkTokenFactory: ContractFactory;

let roles: Roles;

before(async () => {
  roles = (await getUsers()).roles;
  basicConsumerFactory = await ethers.getContractFactory(
    "src/v0.6/tests/BasicConsumer.sol:BasicConsumer",
    roles.defaultAccount,
  );
  oracleFactory = await ethers.getContractFactory("src/v0.6/Oracle.sol:Oracle", roles.oracleNode);
  linkTokenFactory = await ethers.getContractFactory("LinkToken", roles.defaultAccount);
});

describe("BasicConsumer", () => {
  const specId = "0x4c7b7ffb66b344fbaa64995af81e355a".padEnd(66, "0");
  const currency = "USD";
  const payment = toWei("1");
  let link: Contract;
  let oc: Contract;
  let cc: Contract;

  beforeEach(async () => {
    link = await linkTokenFactory.connect(roles.defaultAccount).deploy();
    oc = await oracleFactory.connect(roles.oracleNode).deploy(link.address);
    cc = await basicConsumerFactory.connect(roles.defaultAccount).deploy(link.address, oc.address, specId);
  });

  it("has a predictable gas price", async () => {
    const rec = await ethers.provider.getTransactionReceipt(cc.deployTransaction.hash ?? "");
    assert.isBelow(rec.gasUsed?.toNumber() ?? -1, 1750000);
  });

  describe("#requestEthereumPrice", () => {
    describe("without LINK", () => {
      it("reverts", async () => await expect(cc.requestEthereumPrice(currency, payment)).to.be.reverted);
    });

    describe("with LINK", () => {
      beforeEach(async () => {
        await link.transfer(cc.address, toWei("1"));
      });

      it("triggers a log event in the Oracle contract", async () => {
        const tx = await cc.requestEthereumPrice(currency, payment);
        const receipt = await tx.wait();

        const log = receipt?.logs?.[3];
        assert.equal(log?.address.toLowerCase(), oc.address.toLowerCase());

        const request = decodeRunRequest(log);
        const expected = {
          path: ["USD"],
          get: "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY",
        };

        assert.equal(toHex(specId), request.specId);
        bigNumEquals(toWei("1"), request.payment);
        assert.equal(cc.address.toLowerCase(), request.requester.toLowerCase());
        assert.equal(1, request.dataVersion);
        assert.deepEqual(expected, cbor.decodeFirstSync(request.data));
      });

      it("has a reasonable gas cost", async () => {
        const tx = await cc.requestEthereumPrice(currency, payment);
        const receipt = await tx.wait();

        assert.isBelow(receipt?.gasUsed?.toNumber() ?? -1, 131515);
      });
    });
  });

  describe("#fulfillOracleRequest", () => {
    const response = ethers.utils.formatBytes32String("1,000,000.00");
    let request: RunRequest;

    beforeEach(async () => {
      await link.transfer(cc.address, toWei("1"));
      const tx = await cc.requestEthereumPrice(currency, payment);
      const receipt = await tx.wait();

      request = decodeRunRequest(receipt?.logs?.[3]);
    });

    it("records the data given to it by the oracle", async () => {
      await oc.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));

      const currentPrice = await cc.currentPrice();
      assert.equal(currentPrice, response);
    });

    it("logs the data given to it by the oracle", async () => {
      const tx = await oc.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(request, response));
      const receipt = await tx.wait();

      assert.equal(2, receipt?.logs?.length);
      const log = receipt?.logs?.[1];

      assert.equal(log?.topics[2], response);
    });

    describe("when the consumer does not recognize the request ID", () => {
      let otherRequest: RunRequest;

      beforeEach(async () => {
        // Create a request directly via the oracle, rather than through the
        // chainlink client (consumer). The client should not respond to
        // fulfillment of this request, even though the oracle will faithfully
        // forward the fulfillment to it.
        const args = encodeOracleRequest(
          toHex(specId),
          cc.address,
          basicConsumerFactory.interface.getSighash("fulfill"),
          43,
          constants.HashZero,
        );
        const tx = await link.transferAndCall(oc.address, 0, args);
        const receipt = await tx.wait();

        otherRequest = decodeRunRequest(receipt?.logs?.[2]);
      });

      it("does not accept the data provided", async () => {
        d("otherRequest %s", otherRequest);
        await oc.connect(roles.oracleNode).fulfillOracleRequest(...convertFufillParams(otherRequest, response));

        const received = await cc.currentPrice();

        assert.equal(ethers.utils.parseBytes32String(received), "");
      });
    });

    describe("when called by anyone other than the oracle contract", () => {
      it("does not accept the data provided", async () => {
        await evmRevert(cc.connect(roles.oracleNode).fulfill(request.requestId, response));

        const received = await cc.currentPrice();
        assert.equal(ethers.utils.parseBytes32String(received), "");
      });
    });
  });

  describe("#cancelRequest", () => {
    const depositAmount = toWei("1");
    let request: RunRequest;

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount);
      const tx = await cc.requestEthereumPrice(currency, payment);
      const receipt = await tx.wait();

      request = decodeRunRequest(receipt.logs?.[3]);
    });

    describe("before 5 minutes", () => {
      it("cant cancel the request", () =>
        evmRevert(
          cc
            .connect(roles.consumer)
            .cancelRequest(oc.address, request.requestId, request.payment, request.callbackFunc, request.expiration),
        ));
    });

    describe("after 5 minutes", () => {
      it("can cancel the request", async () => {
        await increaseTime5Minutes(ethers.provider);

        await cc
          .connect(roles.consumer)
          .cancelRequest(oc.address, request.requestId, request.payment, request.callbackFunc, request.expiration);
      });
    });
  });

  describe("#withdrawLink", () => {
    const depositAmount = toWei("1");

    beforeEach(async () => {
      await link.transfer(cc.address, depositAmount);
      const balance = await link.balanceOf(cc.address);
      bigNumEquals(balance, depositAmount);
    });

    it("transfers LINK out of the contract", async () => {
      await cc.connect(roles.consumer).withdrawLink();
      const ccBalance = await link.balanceOf(cc.address);
      const consumerBalance = BigNumber.from(await link.balanceOf(await roles.consumer.getAddress()));
      bigNumEquals(ccBalance, 0);
      bigNumEquals(consumerBalance, depositAmount);
    });
  });
});
