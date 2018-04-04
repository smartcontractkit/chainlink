'use strict';

require('./support/helpers.js')

contract('Args', () => {
  let Args = artifacts.require("Args.sol");
  let args;

  beforeEach(async () => {
    args = await Args.new({});
  });

  it("has a limited public interface", () => {
    checkPublicABI(Args, [
      "add",
      "addBytes32",
      "addBytes32Array",
      "fireEvent"
    ]);
  });

  function parseArgsEvent(tx) {
    let data = util.toBuffer(tx.receipt.logs[0].data);
    return abi.rawDecode(["bytes", "bytes", "bytes"], data);
  }
  describe("#add", () => {
    it("stores and logs keys and values", async () => {
      await args.add("first", "word!!");
      let tx = await args.fireEvent();
      let [types, names, values] = parseArgsEvent(tx);

      assert.equal(types.toString(), "string,");
      assert.equal(names.toString(), "first,");
      assert.equal(values.toString(), lPad("\x06") + rPad("word!!"));
    });

    it("handles two entries", async () => {
      await args.add("first", "uno");
      await args.add("second", "dos");
      let tx = await args.fireEvent();
      let [types, names, values] = parseArgsEvent(tx);

      assert.equal(types.toString(), "string,string,");
      assert.equal(names.toString(), "first,second,");
      let val1 = lPadHex("03") + toHex(rPad("uno"));
      let val2 = lPadHex("03") + toHex(rPad("dos"));
      assert.equal(toHex(values.toString()), val1 + val2);
    });

    it("handles multiple entries", async () => {
      await args.add("first", "uno");
      await args.add("second", "dos");
      await args.add("third", "tres");
      let tx = await args.fireEvent();
      let [types, names, values] = parseArgsEvent(tx);

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
      await args.addBytes32("word", "bytes32 4 LIFE");
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];
      let params = abi.rawDecode(["bytes", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, name, value] = params;

      assert.equal(type.toString(), "bytes32,");
      assert.equal(name.toString(), "word,");
      assert.equal(value.toString(), rPad("bytes32 4 LIFE"));
    });
  });

  describe("#addBytes32Array", () => {
    it("stores and logs keys and values", async () => {
      await args.addBytes32Array("word", ["seinfeld", '"4"', "LIFE"]);
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];
      let params = abi.rawDecode(["bytes", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, name, value] = params;

      assert.equal(type.toString(), "bytes32,");
      assert.equal(name.toString(), "word,");
      let wantLen = lPad("\x03");
      let wantVals = rPad("seinfeld") + rPad('"4"') + rPad("LIFE");
      assert.equal(value.toString(), wantLen + wantVals);
    });
  });
});
