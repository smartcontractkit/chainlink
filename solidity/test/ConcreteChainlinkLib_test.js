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
      "addBytes32",
      "addBytes32Array",
      "fireEvent"
    ]);
  });

  function parseCCLEvent(tx) {
    let data = util.toBuffer(tx.receipt.logs[0].data);
    return abi.rawDecode(["bytes", "bytes", "bytes", "bytes"], data);
  }

  describe("#add", () => {
    it("stores and logs keys and values", async () => {
      await ccl.add("first", "word!!");
      let tx = await ccl.fireEvent();
      let [names, types, values, payload] = parseCCLEvent(tx);

      assert.equal(names.toString(), "first,");
      assert.equal(types.toString(), "string,");
      assert.equal(values.toString(), lPad("\x06") + rPad("word!!"));
    });

    it("handles two entries", async () => {
      await ccl.add("first", "uno");
      await ccl.add("second", "dos");
      let tx = await ccl.fireEvent();
      let [names, types, values, payload] = parseCCLEvent(tx);

      assert.equal(types.toString(), "string,string,");
      assert.equal(names.toString(), "first,second,");
      let val1 = lPadHex("03") + toHex(rPad("uno"));
      let val2 = lPadHex("03") + toHex(rPad("dos"));
      assert.equal(toHex(values.toString()), val1 + val2);
    });

    it("handles multiple entries", async () => {
      await ccl.add("first", "uno");
      await ccl.add("second", "dos");
      await ccl.add("third", "tres");
      let tx = await ccl.fireEvent();
      let [names, types, values, payload] = parseCCLEvent(tx);

      assert.equal(types.toString(), "string,string,string,");
      assert.equal(names.toString(), "first,second,third,");
      let val1 = lPadHex("03") + toHex(rPad("uno"));
      let val2 = lPadHex("03") + toHex(rPad("dos"));
      let val3 = lPadHex("04") + toHex(rPad("tres"));
      assert.equal(toHex(values.toString()), val1 + val2 + val3);
    });
  });

  describe("#addBytes32", () => {
    it("stores and logs keys and values", async () => {
      await ccl.addBytes32("word", "evm 4 LIFE");
      let tx = await ccl.fireEvent();
      let [names, types, values, payload] = parseCCLEvent(tx);

      assert.equal(types.toString(), "bytes32,");
      assert.equal(names.toString(), "word,");
      assert.equal(values.toString(), rPad("evm 4 LIFE"));
    });
  });

  describe("#addBytes32Array", () => {
    it("stores and logs keys and values", async () => {
      await ccl.addBytes32Array("word", ["seinfeld", '"4"', "LIFE"]);
      let tx = await ccl.fireEvent();
      let [names, types, values, payload] = parseCCLEvent(tx);

      assert.equal(types.toString(), "bytes32[],");
      assert.equal(names.toString(), "word,");
      let wantLen = lPad("\x03");
      let wantVals = rPad("seinfeld") + rPad('"4"') + rPad("LIFE");
      assert.equal(values.toString(), wantLen + wantVals);
    });
  });
});
