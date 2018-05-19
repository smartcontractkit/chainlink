'use strict';

require('./support/helpers.js')

contract('Oracle', () => {
  let Oracle = artifacts.require("Oracle.sol");
  let LinkToken = artifacts.require("LinkToken.sol");
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let fHash = functionSelector("requestedBytes32(bytes32,bytes32)");;
  let specId = "4c7b7ffb66b344fbaa64995af81e355a";
  let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
  let link, oc;

  beforeEach(async () => {
    link = await LinkToken.new();
    oc = await Oracle.new(link.address, {from: oracleNode});
  });

  it("has a limited public interface", () => {
    checkPublicABI(Oracle, [
      "cancel",
      "fulfillData",
      "onTokenTransfer",
      "owner",
      "requestData",
      "specAndRun",
      "transferOwnership",
      "withdraw",
    ]);
  });

  describe("#transferOwnership", () => {
    context("when called by the owner", () => {
      beforeEach( async () => {
        await oc.transferOwnership(stranger, {from: oracleNode});
      });

      it("can change the owner", async () => {
        let owner = await oc.owner.call();
        assert.isTrue(web3.isAddress(owner));
        assert.equal(stranger, owner);
      });
    });

    context("when called by a non-owner", () => {
      it("cannot change the owner", async () => {
        await assertActionThrows(async () => {
          await oc.transferOwnership(stranger, {from: stranger});
        });
      });
    });
  });

  describe("#onTokenTransfer", () => {
    context("when called from the LINK token", () => {
      it("triggers the intended method", async () => {
        let callData = requestDataBytes(specId, to, fHash, "id", "");

        let tx = await link.transferAndCall(oc.address, 0, callData);
        assert.equal(3, tx.receipt.logs.length)
      });

      context("with no data", () => {
        it("reverts", async () => {
          await assertActionThrows(async () => {
            await link.transferAndCall(oc.address, 0, "");
          });
        });
      });
    });

    context("when called from any address but the LINK token", () => {
      it("triggers the intended method", async () => {
        let callData = requestDataBytes(specId, to, fHash, "id", "");

        await assertActionThrows(async () => {
          let tx = await oc.onTokenTransfer(oracleNode, 0, callData);
        });
      });
    });
  });

  describe("#requestData", () => {
    context("when called through the LINK token", () => {
      let log, tx;
      beforeEach(async () => {
        let args = requestDataBytes(specId, to, fHash, "id", "");
        tx = await requestDataFrom(oc, link, 0, args);
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2];
      });

      it("logs an event", async () => {
        assert.equal(specId, web3.toUtf8(log.topics[2]));
      });

      it("uses the expected event signature", async () => {
        // If updating this test, be sure to update services.RunLogTopic.
        let eventSignature = "0x3fab86a1207bdcfe3976d0d9df25f263d45ae8d381a60960559771a2b223974d";
        assert.equal(eventSignature, log.topics[0]);
      });
    });

    context("when not called through the LINK token", () => {
      it("reverts", async () => {
        await assertActionThrows(async () => {
          await oc.requestData(1, specId, to, fHash, "id", "", {from: oracleNode});
        });
      });
    });
  });

  describe("#specAndRun", () => {
    context("when called through the LINK token", () => {
      let log, tx, amount;
      beforeEach(async () => {
        amount = 1337;
        let args = specAndRunBytes(to, fHash, "requestId", "");
        tx = await link.transferAndCall(oc.address, amount, args);
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2];
      });

      it("logs an event", async () => {
        assert.equal(amount, web3.toBigNumber(log.topics[2]));
      });

      it("uses the expected event signature", () => {
        // If updating this test, be sure to update services.SpecAndRunTopic.
        let eventSignature = "0x40a86f3bd301164dcd67d63d081ecb2db540ac73bafb27eea27d65b3a2694f39";
        assert.equal(eventSignature, log.topics[0]);
      });
    });

    context("when not called through the LINK token", () => {
      it("reverts", async () => {
        await assertActionThrows(async () => {
          await oc.specAndRun(1, to, fHash, "id", "", {from: oracleNode});
        });
      });
    });
  });

  describe("#fulfillData", () => {
    let mock, internalId;
    let requestId = "XID";

    beforeEach(async () => {
      mock = await GetterSetter.new();
      let fHash = functionSelector("requestedBytes32(bytes32,bytes32)");
      let args = requestDataBytes(specId, mock.address, fHash, requestId, "");
      let req = await requestDataFrom(oc, link, 0, args);
      internalId = req.receipt.logs[2].topics[1];
    });

    context("when called by a non-owner", () => {
      it("raises an error", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(internalId, "Hello World!", {from: stranger});
        });
      });
    });

    context("when called by an owner", () => {
      it("raises an error if the request ID does not exist", async () => {
        await assertActionThrows(async () => {
          await oc.fulfillData(0xdeadbeef, "Hello World!", {from: oracleNode});
        });
      });

      it("sets the value on the requested contract", async () => {
        await oc.fulfillData(internalId, "Hello World!", {from: oracleNode});

        let mockRequestId = await mock.requestId.call();
        assert.equal(requestId.toString(), web3.toUtf8(mockRequestId));

        let currentValue = await mock.getBytes32.call();
        assert.equal("Hello World!", web3.toUtf8(currentValue));
      });

      it("does not allow a request to be fulfilled twice", async () => {
        await oc.fulfillData(internalId, "First message!", {from: oracleNode});
        await assertActionThrows(async () => {
          await oc.fulfillData(internalId, "Second message!!", {from: oracleNode});
        });
      });
    });
  });

  describe("#withdraw", () => {
    context("without reserving funds via requestData", () => {
      it("does nothing", async () => {
        let balance = await link.balanceOf(oracleNode);
        assert.equal(0, balance);
        await oc.withdraw(oracleNode, {from: oracleNode});
        balance = await link.balanceOf(oracleNode);
        assert.equal(0, balance);
      });
    });

    context("reserving funds via requestData", () => {
      let log, tx, mock, internalId, amount;
      beforeEach(async () => {
        amount = 15;
        mock = await GetterSetter.new();
        let args = requestDataBytes(specId, mock.address, fHash, "id", "");
        tx = await requestDataFrom(oc, link, amount, args);
        assert.equal(3, tx.receipt.logs.length)

        log = tx.receipt.logs[2];
        internalId = log.topics[1];
      });

      context("but not freeing funds w fulfillData", () => {
        it("does not transfer funds", async () => {
          await oc.withdraw(oracleNode, {from: oracleNode});
          let balance = await link.balanceOf(oracleNode);
          assert.equal(0, balance);
        });
      });

      context("and freeing funds", () => {
        beforeEach(async () => {
          await oc.fulfillData(internalId, "Hello World!", {from: oracleNode});
        });

        it("allows transfer of funds by owner to specified address", async () => {
          await oc.withdraw(stranger, {from: oracleNode});
          let balance = await link.balanceOf(stranger);
          assert.equal(amount, balance);
        });

        it("does not allow a transfer of funds by non-owner", async () => {
          await assertActionThrows(async () => {
            await oc.withdraw(stranger, {from: stranger});
          });
          let balance = await link.balanceOf(stranger);
          assert.equal(0, balance);
        });
      });
    });
  });

  describe("#cancel", () => {
    context("with no pending requests", () => {
      it("fails", async () => {
        await assertActionThrows(async () => {
          await oc.cancel(1337, {from: stranger});
        });
      });
    });

    context("with a pending request", () => {
      let log, tx, mock, requestAmount, startingBalance;
      let requestId = "requestId";
      beforeEach(async () => {
        startingBalance = 100;
        requestAmount = 20;

        mock = await GetterSetter.new({from: consumer});
        await link.transfer(consumer, startingBalance);

        let args = requestDataBytes(specId, consumer, fHash, requestId, "");
        tx = await link.transferAndCall(oc.address, requestAmount, args, {from: consumer});
        assert.equal(3, tx.receipt.logs.length)
      });

      it("has correct initial balances", async () => {
        let oracleBalance = await link.balanceOf(oc.address);
        assert.equal(requestAmount, oracleBalance);

        let consumerAmount = await link.balanceOf(consumer);
        assert.equal(startingBalance - requestAmount, consumerAmount);
      });

      context("from a stranger", () => {
        it("fails", async () => {
          await assertActionThrows(async () => {
            await oc.cancel(requestId, {from: stranger});
          });
        });
      });

      context("from the requester", () => {
        it("refunds the correct amount", async () => {
          await oc.cancel(requestId, {from: consumer});
          let balance = await link.balanceOf(consumer);
          assert.equal(startingBalance, balance); // 100
        });

        context("canceling twice", () => {
          it("fails", async () => {
            await oc.cancel(requestId, {from: consumer});
            await assertActionThrows(async () => {
              await oc.cancel(requestId, {from: consumer});
            });
          });
        });
      });
    });
  });
});
