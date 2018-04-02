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
      let params = abi.rawDecode(["bytes", "bytes", "bytes"], util.toBuffer(log.data));
      let [type, name, value] = params;

      assert.equal(type.toString(), "string,");
      assert.equal(name.toString(), "first,");
      assert.equal(value.toString(), lPad("\x04") + rPad("word"));
    });

    it("handles multiple entries", async () => {
      await args.add("first", "uno");
      await args.add("second", "dos");
      let tx = await args.fireEvent();
      let log = tx.receipt.logs[0];

      let params = abi.rawDecode(["bytes", "bytes", "bytes"], util.toBuffer(log.data));
      let [types, names, values] = params;

      assert.equal(types.toString(), "string,string,");
      assert.equal(names.toString(), "first,second,");
      let val1 = lPad("\x03") + rPad("uno");
      let val2 = lPad("\x03") + rPad("dos");
      assert.equal(values.toString(), val1 + val2);
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
