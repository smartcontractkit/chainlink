var Hash = require("./hash");
var Bytes = require("./bytes");
var Desubits = require("./desubits");

// Bytes -> Bytes
var bytesAddChecksum = function bytesAddChecksum(bytes) {
  var hash = Hash.keccak256(bytes);
  var sum = Bytes.slice(0, 1, hash);
  return Bytes.concat(bytes, sum);
};

// Bytes -> Bool
var bytesChecksum = function bytesChecksum(bytes) {
  var length = Bytes.length(bytes);
  var prefix = Bytes.slice(0, length - 1, bytes);
  return bytesAddChecksum(prefix) === bytes;
};

// () ~> Passphrase
var create = function create() {
  var bytes = bytesAddChecksum(Bytes.random(11));
  var seed = Desubits.fromBytes(bytes);
  var passphrase = seed.replace(/([a-z]{8})/g, "$1 ");
  return passphrase;
};

// Passphrase -> Bytes
var toBytes = function toBytes(passphrase) {
  var seed = passphrase.replace(/ /g, "");
  var bytes = Desubits.toBytes(passphrase);
  return bytes;
};

// Passphrase -> Bool
var checksum = function checksum(passphrase) {
  return bytesChecksum(toBytes(passphrase));
};

// Passphrase -> Bytes
var toMasterKey = function toMasterKey(passphrase) {
  var hash = Hash.keccak256;
  var bytes = toBytes(passphrase);
  for (var i = 0, l = Math.pow(2, 12); i < l; ++i) {
    bytes = hash(bytes);
  }return bytes;
};

module.exports = {
  create: create,
  checksum: checksum,
  toMasterKey: toMasterKey
};