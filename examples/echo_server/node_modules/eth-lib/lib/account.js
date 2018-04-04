var _slicedToArray = function () { function sliceIterator(arr, i) { var _arr = []; var _n = true; var _d = false; var _e = undefined; try { for (var _i = arr[Symbol.iterator](), _s; !(_n = (_s = _i.next()).done); _n = true) { _arr.push(_s.value); if (i && _arr.length === i) break; } } catch (err) { _d = true; _e = err; } finally { try { if (!_n && _i["return"]) _i["return"](); } finally { if (_d) throw _e; } } return _arr; } return function (arr, i) { if (Array.isArray(arr)) { return arr; } else if (Symbol.iterator in Object(arr)) { return sliceIterator(arr, i); } else { throw new TypeError("Invalid attempt to destructure non-iterable instance"); } }; }();

var Bytes = require("./bytes");
var Nat = require("./nat");
var elliptic = require("elliptic");
var rlp = require("./rlp");
var secp256k1 = new elliptic.ec("secp256k1"); // eslint-disable-line

var _require = require("./hash"),
    keccak256 = _require.keccak256,
    keccak256s = _require.keccak256s;

var create = function create(entropy) {
  var innerHex = keccak256(Bytes.concat(Bytes.random(32), entropy || Bytes.random(32)));
  var middleHex = Bytes.concat(Bytes.concat(Bytes.random(32), innerHex), Bytes.random(32));
  var outerHex = keccak256(middleHex);
  return fromPrivate(outerHex);
};

var toChecksum = function toChecksum(address) {
  var addressHash = keccak256s(address.slice(2));
  var checksumAddress = "0x";
  for (var i = 0; i < 40; i++) {
    checksumAddress += parseInt(addressHash[i + 2], 16) > 7 ? address[i + 2].toUpperCase() : address[i + 2];
  }return checksumAddress;
};

var fromPrivate = function fromPrivate(privateKey) {
  var buffer = new Buffer(privateKey.slice(2), "hex");
  var ecKey = secp256k1.keyFromPrivate(buffer);
  var publicKey = "0x" + ecKey.getPublic(false, 'hex').slice(2);
  var publicHash = keccak256(publicKey);
  var address = toChecksum("0x" + publicHash.slice(-40));
  return {
    address: address,
    privateKey: privateKey
  };
};

var encodeSignature = function encodeSignature(_ref) {
  var _ref2 = _slicedToArray(_ref, 3),
      v = _ref2[0],
      r = _ref2[1],
      s = _ref2[2];

  return Bytes.flatten([r, s, v]);
};

var decodeSignature = function decodeSignature(hex) {
  return [Bytes.slice(64, 65, hex), Bytes.slice(0, 32, hex), Bytes.slice(32, 64, hex)];
};

var makeSign = function makeSign(addToV) {
  return function (hash, privateKey) {
    var signature = secp256k1.keyFromPrivate(new Buffer(privateKey.slice(2), "hex")).sign(new Buffer(hash.slice(2), "hex"), { canonical: true });
    return encodeSignature([Bytes.pad(1, Bytes.fromNumber(addToV + signature.recoveryParam)), Bytes.pad(32, Bytes.fromNat("0x" + signature.r.toString(16))), Bytes.pad(32, Bytes.fromNat("0x" + signature.s.toString(16)))]);
  };
};

var sign = makeSign(27); // v=27|28 instead of 0|1...

var recover = function recover(hash, signature) {
  var vals = decodeSignature(signature);
  var vrs = { v: Bytes.toNumber(vals[0]), r: vals[1].slice(2), s: vals[2].slice(2) };
  var ecPublicKey = secp256k1.recoverPubKey(new Buffer(hash.slice(2), "hex"), vrs, vrs.v < 2 ? vrs.v : 1 - vrs.v % 2); // because odd vals mean v=0... sadly that means v=0 means v=1... I hate that
  var publicKey = "0x" + ecPublicKey.encode("hex", false).slice(2);
  var publicHash = keccak256(publicKey);
  var address = toChecksum("0x" + publicHash.slice(-40));
  return address;
};

var transactionSigningData = function transactionSigningData(tx) {
  return rlp.encode([Bytes.fromNat(tx.nonce), Bytes.fromNat(tx.gasPrice), Bytes.fromNat(tx.gas), tx.to.toLowerCase(), Bytes.fromNat(tx.value), tx.data, Bytes.fromNat(tx.chainId || "0x1"), "0x", "0x"]);
};

var signTransaction = function signTransaction(tx, privateKey) {
  var signingData = transactionSigningData(tx);
  var signature = makeSign(Nat.toNumber(tx.chainId || "0x1") * 2 + 35)(keccak256(signingData), privateKey);
  var rawTransaction = rlp.decode(signingData).slice(0, 6).concat(decodeSignature(signature));
  return rlp.encode(rawTransaction);
};

var recoverTransaction = function recoverTransaction(rawTransaction) {
  var values = rlp.decode(rawTransaction);
  var signature = encodeSignature(values.slice(6, 9));
  var recovery = Bytes.toNumber(values[6]);
  var extraData = recovery < 35 ? [] : [Bytes.fromNumber(recovery - 35 >> 1), "0x", "0x"];
  var signingData = values.slice(0, 6).concat(extraData);
  var signingDataHex = rlp.encode(signingData);
  return recover(keccak256(signingDataHex), signature);
};

module.exports = {
  create: create,
  toChecksum: toChecksum,
  fromPrivate: fromPrivate,
  sign: sign,
  recover: recover,
  signTransaction: signTransaction,
  recoverTransaction: recoverTransaction,
  transactionSigningData: transactionSigningData,
  encodeSignature: encodeSignature,
  decodeSignature: decodeSignature
};