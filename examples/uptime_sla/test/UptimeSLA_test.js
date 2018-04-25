'use strict';

require('./support/helpers.js')

contract('UptimeSLA', () => {
  let Link = artifacts.require("../../../solidity/contracts/LinkToken.sol");
  let Oracle = artifacts.require("../../../solidity/contracts/Oracle.sol");
  let SLA = artifacts.require("UptimeSLA.sol");
  let jobId = "4c7b7ffb66b344fbaa64995af81e355a";
  let deposit = 1000000000;
  let link, oc, sla, client, serviceProvider, startAt;

  beforeEach(async () => {
    client = newAddress()
    serviceProvider = newAddress();
    link = await Link.new();
    oc = await Oracle.new(link.address, {from: oracleNode});
    sla = await SLA.new(client, serviceProvider, link.address, oc.address, jobId, {
      value: deposit
    });
    link.transfer(sla.address, web3.toWei(1, 'ether'));
    startAt = await getLatestTimestamp();
  });

  describe("before updates", () => {
    it("does not release money to anyone", async () => {
      assert.equal(await eth.getBalance(sla.address), deposit);
      assert.equal(await eth.getBalance(client), 0);
      assert.equal(await eth.getBalance(serviceProvider), 0);
    });
  });

  describe("#updateUptime", () => {
    it("triggers a log event in the Oracle contract", async () => {
      let tx = await sla.updateUptime("0");

      let events = await getEvents(oc);
      assert.equal(1, events.length)

      let event = events[0]
      assert.equal(web3.toUtf8(event.args.jobId), jobId);

      let decoded = cbor.decodeFirstSync(util.toBuffer(event.args.data));
      assert.deepEqual(
        decoded,
        {"url":"https://status.heroku.com/api/ui/availabilities","path":["data","0","attributes","calculation"]}
      )
    });
  });

  describe("#fulfillData", () => {
    let response = "0x00000000000000000000000000000000000000000000000000000000000f8c4c";
    let requestId;

    beforeEach(async () => {
      await sla.updateUptime("0");
      let event = await getLatestEvent(oc);
      requestId = event.args.id
    });

    context("when the value is below 9999", async () => {
      let response = "0x000000000000000000000000000000000000000000000000000000000000270e";

      it("sends the deposit to the client", async () => {
        await oc.fulfillData(requestId, response, {from: oracleNode})

        assert.equal(await eth.getBalance(sla.address), 0);
        assert.equal(await eth.getBalance(client), deposit);
        assert.equal(await eth.getBalance(serviceProvider), 0);
      });
    });

    context("when the value is 9999 or above", () => {
      let response = "0x000000000000000000000000000000000000000000000000000000000000270f";

      it("does not move the money", async () => {
        await oc.fulfillData(requestId, response, {from: oracleNode})

        assert.equal(await eth.getBalance(sla.address), deposit);
        assert.equal(await eth.getBalance(client), 0);
        assert.equal(await eth.getBalance(serviceProvider), 0);
      });

      context("and a month has passed", () => {
        beforeEach(async () => {
          await fastForwardTo(startAt + days(30));
        });

        it("gives the money back to the service provider", async () => {
          await oc.fulfillData(requestId, response, {from: oracleNode})

          assert.equal(await eth.getBalance(sla.address), 0);
          assert.equal(await eth.getBalance(client), 0);
          assert.equal(await eth.getBalance(serviceProvider), deposit);
        });
      });
    });

    context("when the consumer does not recognize the request ID", () => {
      beforeEach(async () => {

        let fid = functionSelector("fulfill(uint256,bytes32)");
        let args = requestDataBytes(jobId, sla.address, fid, "xid", "");
        await requestDataFrom(oc, link, 0, args);
        let event = await getLatestEvent(oc);
        requestId = event.args.id;
      });

      it("does not accept the data provided", async () => {
        let tx = await sla.updateUptime("0");

        await assertActionThrows(async () => {
          await oc.fulfillData(requestId, response, {from: oracleNode})
        });
      });
    });

    context("when called by anyone other than the oracle contract", () => {
      it("does not accept the data provided", async () => {
        await assertActionThrows(async () => {
          await sla.report(requestId, response, {from: oracleNode})
        });
      });
    });
  });
});
