var BN = require("bn.js");
var Bytes = require("./bytes");

var fromBN = function fromBN(bn) {
  return "0x" + bn.toString("hex");
};

var toBN = function toBN(str) {
  return new BN(str.slice(2), 16);
};

var fromString = function fromString(str) {
  var bn = "0x" + (str.slice(0, 2) === "0x" ? new BN(str.slice(2), 16) : new BN(str, 10)).toString("hex");
  return bn === "0x0" ? "0x" : bn;
};

var toEther = function toEther(wei) {
  return toNumber(div(wei, fromString("10000000000"))) / 100000000;
};

var fromEther = function fromEther(eth) {
  return mul(fromNumber(Math.floor(eth * 100000000)), fromString("10000000000"));
};

var toString = function toString(a) {
  return toBN(a).toString(10);
};

var fromNumber = function fromNumber(a) {
  return typeof a === "string" ? /^0x/.test(a) ? a : "0x" + a : "0x" + new BN(a).toString("hex");
};

var toNumber = function toNumber(a) {
  return toBN(a).toNumber();
};

var toUint256 = function toUint256(a) {
  return Bytes.pad(32, a);
};

var bin = function bin(method) {
  return function (a, b) {
    return fromBN(toBN(a)[method](toBN(b)));
  };
};

var add = bin("add");
var mul = bin("mul");
var div = bin("div");
var sub = bin("sub");

module.exports = {
  toString: toString,
  fromString: fromString,
  toNumber: toNumber,
  fromNumber: fromNumber,
  toEther: toEther,
  fromEther: fromEther,
  toUint256: toUint256,
  add: add,
  mul: mul,
  div: div,
  sub: sub
};