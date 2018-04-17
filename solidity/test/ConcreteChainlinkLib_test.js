'use strict';

require('./support/helpers.js')

contract('ConcreteChainlinkLib', () => {
  let CCL = artifacts.require("examples/ConcreteChainlinkLib.sol");
  let ccl;

  beforeEach(async () => {
    ccl = await CCL.new({});
  });

  it("has a limited public interface", () => {
    checkPublicABI(CCL, [
      "add",
      "addStringArray",
      "closeEvent"
    ]);
  });

  function parseCCLEvent(tx) {
    let data = util.toBuffer(tx.receipt.logs[0].data);
    return abi.rawDecode(["bytes"], data);
  }

  describe("#close", () => {
    it("handles empty payloads", async () => {
      let tx = await ccl.closeEvent();
      let [payload] = parseCCLEvent(tx);
      var decoded = await cbor.decodeFirst(payload);
      assert.deepEqual(decoded, {});
    });
  });

  describe("#add", () => {
    it("stores and logs keys and values", async () => {
      await ccl.add("first", "word!!");
      let tx = await ccl.closeEvent();
      let [payload] = parseCCLEvent(tx);
      var decoded = await cbor.decodeFirst(payload);
      assert.deepEqual(decoded, { "first": "word!!" });
    });

    it("handles two entries", async () => {
      await ccl.add("first", "uno");
      await ccl.add("second", "dos");
      let tx = await ccl.closeEvent();
      let [payload] = parseCCLEvent(tx);
      var decoded = await cbor.decodeFirst(payload);

      assert.deepEqual(decoded, {
        "first": "uno",
        "second": "dos"
      });
    });
  });

  describe("#addStringArray", () => {
    it("stores and logs keys and values", async () => {
      await ccl.addStringArray("word", ["seinfeld", '"4"', "LIFE"]);
      let tx = await ccl.closeEvent();
      let [payload] = parseCCLEvent(tx);
      var decoded = await cbor.decodeFirst(payload);

      assert.deepEqual(decoded, { "word": ["seinfeld", '"4"', "LIFE"] });
    });
  });
});
