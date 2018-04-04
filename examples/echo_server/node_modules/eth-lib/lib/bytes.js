var A = require("./array.js");

var at = function at(bytes, index) {
  return parseInt(bytes.slice(index * 2 + 2, index * 2 + 4), 16);
};

var random = function random(bytes) {
  var rnd = void 0;
  if (typeof window !== "undefined" && window.crypto && window.crypto.getRandomValues) rnd = window.crypto.getRandomValues(new Uint8Array(bytes));else if (typeof require !== "undefined") rnd = require("c" + "rypto").randomBytes(bytes);else throw "Safe random numbers not available.";
  var hex = "0x";
  for (var i = 0; i < bytes; ++i) {
    hex += ("00" + rnd[i].toString(16)).slice(-2);
  }return hex;
};

var length = function length(a) {
  return (a.length - 2) / 2;
};

var flatten = function flatten(a) {
  return "0x" + a.reduce(function (r, s) {
    return r + s.slice(2);
  }, "");
};

var slice = function slice(i, j, bs) {
  return "0x" + bs.slice(i * 2 + 2, j * 2 + 2);
};

var reverse = function reverse(hex) {
  var rev = "0x";
  for (var i = 0, l = length(hex); i < l; ++i) {
    rev += hex.slice((l - i) * 2, (l - i + 1) * 2);
  }
  return rev;
};

var pad = function pad(l, hex) {
  return hex.length === l * 2 + 2 ? hex : pad(l, "0x" + "0" + hex.slice(2));
};

var padRight = function padRight(l, hex) {
  return hex.length === l * 2 + 2 ? hex : padRight(l, hex + "0");
};

var toArray = function toArray(hex) {
  var arr = [];
  for (var i = 2, l = hex.length; i < l; i += 2) {
    arr.push(parseInt(hex.slice(i, i + 2), 16));
  }return arr;
};

var fromArray = function fromArray(arr) {
  var hex = "0x";
  for (var i = 0, l = arr.length; i < l; ++i) {
    var b = arr[i];
    hex += (b < 16 ? "0" : "") + b.toString(16);
  }
  return hex;
};

var toUint8Array = function toUint8Array(hex) {
  return new Uint8Array(toArray(hex));
};

var fromUint8Array = function fromUint8Array(arr) {
  return fromArray([].slice.call(arr, 0));
};

var fromNumber = function fromNumber(num) {
  var hex = num.toString(16);
  return hex.length % 2 === 0 ? "0x" + hex : "0x0" + hex;
};

var toNumber = function toNumber(hex) {
  return parseInt(hex.slice(2), 16);
};

var concat = function concat(a, b) {
  return a.concat(b.slice(2));
};

var fromNat = function fromNat(bn) {
  return bn === "0x0" ? "0x" : bn.length % 2 === 0 ? bn : "0x0" + bn.slice(2);
};

var toNat = function toNat(bn) {
  return bn[2] === "0" ? "0x" + bn.slice(3) : bn;
};

var fromAscii = function fromAscii(ascii) {
  var hex = "0x";
  for (var i = 0; i < ascii.length; ++i) {
    hex += ("00" + ascii.charCodeAt(i).toString(16)).slice(-2);
  }return hex;
};

var toAscii = function toAscii(hex) {
  var ascii = "";
  for (var i = 2; i < hex.length; i += 2) {
    ascii += String.fromCharCode(parseInt(hex.slice(i, i + 2), 16));
  }return ascii;
};

// From https://gist.github.com/pascaldekloe/62546103a1576803dade9269ccf76330
var fromString = function fromString(s) {
  var makeByte = function makeByte(uint8) {
    var b = uint8.toString(16);
    return b.length < 2 ? "0" + b : b;
  };
  var bytes = "0x";
  for (var ci = 0; ci != s.length; ci++) {
    var c = s.charCodeAt(ci);
    if (c < 128) {
      bytes += makeByte(c);
      continue;
    }
    if (c < 2048) {
      bytes += makeByte(c >> 6 | 192);
    } else {
      if (c > 0xd7ff && c < 0xdc00) {
        if (++ci == s.length) return null;
        var c2 = s.charCodeAt(ci);
        if (c2 < 0xdc00 || c2 > 0xdfff) return null;
        c = 0x10000 + ((c & 0x03ff) << 10) + (c2 & 0x03ff);
        bytes += makeByte(c >> 18 | 240);
        bytes += makeByte(c >> 12 & 63 | 128);
      } else {
        // c <= 0xffff
        bytes += makeByte(c >> 12 | 224);
      }
      bytes += makeByte(c >> 6 & 63 | 128);
    }
    bytes += makeByte(c & 63 | 128);
  }
  return bytes;
};

var toString = function toString(bytes) {
  var s = '';
  var i = 0;
  var l = length(bytes);
  while (i < l) {
    var c = at(bytes, i++);
    if (c > 127) {
      if (c > 191 && c < 224) {
        if (i >= l) return null;
        c = (c & 31) << 6 | at(bytes, i) & 63;
      } else if (c > 223 && c < 240) {
        if (i + 1 >= l) return null;
        c = (c & 15) << 12 | (at(bytes, i) & 63) << 6 | at(bytes, ++i) & 63;
      } else if (c > 239 && c < 248) {
        if (i + 2 >= l) return null;
        c = (c & 7) << 18 | (at(bytes, i) & 63) << 12 | (at(bytes, ++i) & 63) << 6 | at(bytes, ++i) & 63;
      } else return null;
      ++i;
    }
    if (c <= 0xffff) s += String.fromCharCode(c);else if (c <= 0x10ffff) {
      c -= 0x10000;
      s += String.fromCharCode(c >> 10 | 0xd800);
      s += String.fromCharCode(c & 0x3FF | 0xdc00);
    } else return null;
  }
  return s;
};

module.exports = {
  random: random,
  length: length,
  concat: concat,
  flatten: flatten,
  slice: slice,
  reverse: reverse,
  pad: pad,
  padRight: padRight,
  fromAscii: fromAscii,
  toAscii: toAscii,
  fromString: fromString,
  toString: toString,
  fromNumber: fromNumber,
  toNumber: toNumber,
  fromNat: fromNat,
  toNat: toNat,
  fromArray: fromArray,
  toArray: toArray,
  fromUint8Array: fromUint8Array,
  toUint8Array: toUint8Array
};