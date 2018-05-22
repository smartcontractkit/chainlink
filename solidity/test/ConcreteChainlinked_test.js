'use strict';

require('./support/helpers.js')

contract('ConcreteChainlinked', () => {
  let Link = artifacts.require("LinkToken.sol");
  let Oracle = artifacts.require("Oracle.sol");
  let Chainlinked = artifacts.require("examples/ConcreteChainlinked.sol");
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let specId = "4c7b7ffb66b344fbaa64995af81e355a";
  let cc, gs, oc, link;

  beforeEach(async () => {
    link = await Link.new();
    oc = await Oracle.new(link.address);
    gs = await GetterSetter.new();
    cc = await Chainlinked.new(link.address, oc.address);
  });

  describe("#newRun", () => {
    it("forwards the information to the oracle contract through the link token", async () => {
      let tx = await cc.publicNewRun(
        specId,
        gs.address,
        "requestedBytes32(bytes32,bytes32)");

      assert.equal(1, tx.receipt.logs.length);
      let [jId, cbAddr, cbFId, cborData] = decodeRunABI(tx.receipt.logs[0]);
      let params = await cbor.decodeFirst(cborData);

      assert.equal(specId, jId);
      assert.equal(gs.address, `0x${cbAddr}`);
      assert.equal("ed53e511", toHex(cbFId));
      assert.deepEqual({}, params);
    });
  });

  describe("#chainlinkRequest(Run)", () => {
    it("emits an event from the contract showing the run ID", async () => {
      let tx = await cc.publicRequestRun(specId, cc.address, "fulfillRequest(bytes32,bytes32)", 0);

      let events = await getEvents(cc);
      assert.equal(1, events.length);
      let event = events[0];
      assert.equal(event.event, "ChainlinkRequested");
    });
  });

  describe("#chainlinkRequest(SpecAndRun)", () => {
    it("emits an event from the contract showing the run ID", async () => {
      let tx = await cc.publicRequestSpecAndRun([], cc.address, "fulfillRequest(bytes32,bytes32)", 0);

      let events = await getEvents(cc);
      assert.equal(1, events.length);
      let event = events[0];
      assert.equal(event.event, "ChainlinkRequested");
    });
  });

  describe("#cancelChainlinkRequest", () => {
    let requestId;

    beforeEach(async () => {
      await cc.publicRequestRun(specId, cc.address, "fulfillRequest(bytes32,bytes32)", 0);
      requestId = (await getLatestEvent(cc)).args.id;
    });

    it("emits an event from the contract showing the run was cancelled", async () => {
      let tx = await cc.publicCancelRequest(requestId);

      let events = await getEvents(cc);
      assert.equal(1, events.length);
      let event = events[0];
      assert.equal(event.event, "ChainlinkCancelled");
      assert.equal(requestId, event.args.id);
    });

    context("when the request ID is no longer unfulfilled", () => {
      beforeEach(async () => {
        await cc.publicCancelRequest(requestId);
      });

      it("throws an error", async () => {
        await assertActionThrows(async () => {
          await cc.publicCancelRequest(requestId);
        });
      });
    });
  });

  describe("#checkChainlinkFulfillment(modifier)", () => {
    let internalId, requestId;

    beforeEach(async () => {
      await cc.publicRequestRun(specId, cc.address, "fulfillRequest(bytes32,bytes32)", 0);
      requestId = (await getLatestEvent(cc)).args.id;
      internalId = (await getLatestEvent(oc)).args.internalId;
    });

    it("emits an event marking the request cancelled", async () => {
      await oc.fulfillData(internalId, "hi mom!");

      let events = await getEvents(cc);
      assert.equal(1, events.length);
      let event = events[0];
      assert.equal(event.event, "ChainlinkFulfilled");
      assert.equal(requestId, event.args.id);
    });
  });
});
