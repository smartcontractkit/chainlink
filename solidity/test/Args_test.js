'use strict';

require('./support/helpers.js')

contract('Args', () => {
  let Args = artifacts.require("Args.sol");
  let args;

  beforeEach(async () => {
    args = await Args.new({});
  });

  it("has a limited public interface", () => {
    checkPublicABI(Args, [ "add", "addBytes32", "addBytes32Array", "fireEvent", ]);
  });

  describe("#add", () => {
    it("stores and logs keys and values", async () => {
      await args.add("first", "word");
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];
      let params = abi.rawDecode(["bytes", "uint16[]", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, valueLength, name, value] = params;

      assert.equal(type.toString(), "string,");
      assert.equal(name.toString(), "first,");
      assert.equal(valueLength, 4);
      assert.equal(value.toString(), "word");
    });

    it("handles multiple entries", async () => {
      await args.add("first", "uno");
      await args.add("second", "dos");
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];

      let params = abi.rawDecode(["bytes", "uint16[]", "bytes", "bytes"], util.toBuffer(log.data));
      let [types, valueLengths, names, values] = params;

      assert.equal(types.toString(), "string,string,");
      assert.equal(valueLengths.toString(), "3,3");
      assert.equal(names.toString(), "first,second,");
      assert.equal(values.toString(), ["unodos"]);
    });
  });

  describe("#addBytes32", () => {
    it("stores and logs keys and values", async () => {
      await args.addBytes32("word", "bytes32 4 LIFE");
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];
      let params = abi.rawDecode(["bytes", "uint16[]", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, valueLength, name, value] = params;

      assert.equal(type.toString(), "bytes32,");
      assert.equal(name.toString(), "word,");
      assert.equal(value.toString(), rPadWord("bytes32 4 LIFE"));
    });
  });

  describe("#addBytes32Array", () => {
    it("stores and logs keys and values", async () => {
      await args.addBytes32Array("word", ["seinfeld", '"4"', "LIFE"]);
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];
      let params = abi.rawDecode(["bytes", "uint16[]", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, valueLength, name, value] = params;

      assert.equal(type.toString(), "bytes32,");
      assert.equal(name.toString(), "word,");
      let wantLen = lPadWord("\x03");
      let wantVals = rPadWord("seinfeld") + rPadWord('"4"') + rPadWord("LIFE");
      assert.equal(value.toString(), wantLen + wantVals);
    });
  });
});
