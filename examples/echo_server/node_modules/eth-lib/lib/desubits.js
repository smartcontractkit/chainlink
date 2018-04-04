var Bytes = require("./bytes");

var inis = "pbtdkgxjfvlrmnsz".split("");

var mids = "auie".split("");

var back = [inis, mids].map(function (chars) {
  var map = {};
  for (var i = 0; i < chars.length; ++i) {
    map[chars[i]] = i;
  }return map;
});

var syllableFromB64 = function syllableFromB64(b64) {
  var ini = b64 >>> 2 & 15;
  var mid = b64 >>> 0 & 3;
  return inis[ini] + mids[mid];
};

var syllableToB64 = function syllableToB64(syllable) {
  var ini = back[0][syllable[0]];
  var mid = back[1][syllable[1]];
  return ini * 4 + mid;
};

var b64sFromBytes = function b64sFromBytes(bytes) {
  //BoooooBoooooBoooooBooooo
  //BoooooooBoooooooBooooooo
  var b64s = [],
      b64 = void 0;
  for (var i = 0, l = Math.ceil(bytes.length * 8 / 6); i < l; ++i) {
    var j = i / 8 * 6 | 0;
    b64s.push(i % 4 === 0 ? bytes[j + 0] / 4 | 0 : i % 4 === 1 ? bytes[j + 0] % 4 * 16 + (bytes[j + 1] / 16 | 0) : i % 4 === 2 ? bytes[j + 0] % 16 * 4 + (bytes[j + 1] / 64 | 0) : bytes[j + 0] % 64 * 1);
  }
  return b64s;
};

var b64sToBytes = function b64sToBytes(b64s) {
  var bytes = [];
  for (var i = 0, l = Math.floor(b64s.length * 6 / 8); i < l; ++i) {
    var j = i / 6 * 8 | 0;
    bytes.push(i % 3 === 0 ? b64s[j + 0] % 64 * 4 + (b64s[j + 1] / 16 | 0) : i % 3 === 1 ? b64s[j + 0] % 16 * 16 + (b64s[j + 1] / 4 | 0) : b64s[j + 0] % 4 * 64 + (b64s[j + 1] / 1 | 0));
  }
  return bytes;
};

var fromBytes = function fromBytes(bytes) {
  return b64sFromBytes(Bytes.toArray(bytes)).map(syllableFromB64).join("");
};

var toBytes = function toBytes(syllables) {
  return Bytes.fromArray(b64sToBytes(syllables.match(/\w\w/g).map(syllableToB64)));
};

module.exports = {
  fromBytes: fromBytes,
  toBytes: toBytes
};