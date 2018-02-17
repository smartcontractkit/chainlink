'use strict';

require('./support/helpers.js')

contract('GetterSetter', () => {
  let GetterSetter = artifacts.require("examples/GetterSetter.sol");
  let requestId = 5432;
  let bytes32 = "Hi Mom!";
  let uint256 = 645746535432;
  let gs;

  beforeEach(async () => {
    gs = await GetterSetter.new();
  });

  describe("#setBytes32Val", () => {
    it("updates the bytes32 value", async () => {
      await gs.setBytes32(bytes32, {from: stranger});

      let currentBytes32 = await gs.getBytes32.call();
      assert.equal(web3.toUtf8(currentBytes32), bytes32);
    });
  });

  describe("#requestedBytes32", () => {
    it("updates the request ID and value", async () => {
      await gs.requestedBytes32(requestId, bytes32, {from: stranger});

      let currentRequestId = await gs.requestId.call();
      assert.equal(currentRequestId, requestId);

      let currentBytes32 = await gs.getBytes32.call();
      assert.equal(web3.toUtf8(currentBytes32), bytes32);
    });
  });

  describe("#setUint256", () => {
    it("updates uint256 value", async () => {
      await gs.setUint256(uint256, {from: stranger});

      let currentUint256 = await gs.getUint256.call();
      assert.equal(currentUint256, uint256);
    });
  });

  describe("#requestedUint256", () => {
    it("updates the request ID and value", async () => {
      await gs.requestedUint256(requestId, uint256, {from: stranger});

      let currentRequestId = await gs.requestId.call();
      assert.equal(currentRequestId, requestId);

      let currentUint256 = await gs.getUint256.call();
      assert.equal(currentUint256, uint256);
    });
  });
});
