const Bytes = require("./bytes");
const Nat = require("./nat");
const elliptic = require("elliptic");
const rlp = require("./rlp");
const secp256k1 = new (elliptic.ec)("secp256k1"); // eslint-disable-line
const {keccak256, keccak256s} = require("./hash");

const create = entropy => {
  const innerHex = keccak256(Bytes.concat(Bytes.random(32), entropy || Bytes.random(32)));
  const middleHex = Bytes.concat(Bytes.concat(Bytes.random(32), innerHex), Bytes.random(32));
  const outerHex = keccak256(middleHex);
  return fromPrivate(outerHex);
}

const toChecksum = address => {
  const addressHash = keccak256s(address.slice(2));
  let checksumAddress = "0x";
  for (let i = 0; i < 40; i++)
    checksumAddress += parseInt(addressHash[i + 2], 16) > 7
      ? address[i + 2].toUpperCase()
      : address[i + 2];
  return checksumAddress;
}

const fromPrivate = privateKey => {
  const buffer = new Buffer(privateKey.slice(2), "hex");
  const ecKey = secp256k1.keyFromPrivate(buffer);
  const publicKey = "0x" + ecKey.getPublic(false, 'hex').slice(2);
  const publicHash = keccak256(publicKey);
  const address = toChecksum("0x" + publicHash.slice(-40));
  return {
    address: address,
    privateKey: privateKey
  }
}

const encodeSignature = ([v, r, s]) =>
  Bytes.flatten([r,s,v]);

const decodeSignature = (hex) =>
  [Bytes.slice(64,65,hex), Bytes.slice(0,32,hex), Bytes.slice(32,64,hex)];

const makeSign = addToV => (hash, privateKey) => {
  const signature = secp256k1
    .keyFromPrivate(new Buffer(privateKey.slice(2), "hex"))
    .sign(new Buffer(hash.slice(2), "hex"), {canonical: true});
  return encodeSignature([
    Bytes.pad(1, Bytes.fromNumber(addToV + signature.recoveryParam)),
    Bytes.pad(32, Bytes.fromNat("0x" + signature.r.toString(16))),
    Bytes.pad(32, Bytes.fromNat("0x" + signature.s.toString(16)))]);
}

const sign = makeSign(27); // v=27|28 instead of 0|1...

const recover = (hash, signature) => {
  const vals = decodeSignature(signature);
  const vrs = {v: Bytes.toNumber(vals[0]), r:vals[1].slice(2), s:vals[2].slice(2)};
  const ecPublicKey = secp256k1.recoverPubKey(new Buffer(hash.slice(2), "hex"), vrs, vrs.v < 2 ? vrs.v : 1 - (vrs.v % 2)); // because odd vals mean v=0... sadly that means v=0 means v=1... I hate that
  const publicKey = "0x" + ecPublicKey.encode("hex", false).slice(2);
  const publicHash = keccak256(publicKey);
  const address = toChecksum("0x" + publicHash.slice(-40));
  return address;
}

const transactionSigningData = tx =>
  rlp.encode([
    Bytes.fromNat(tx.nonce),
    Bytes.fromNat(tx.gasPrice),
    Bytes.fromNat(tx.gas),
    tx.to.toLowerCase(),
    Bytes.fromNat(tx.value),
    tx.data,
    Bytes.fromNat(tx.chainId || "0x1"),
    "0x",
    "0x"]);

const signTransaction = (tx, privateKey) => {
  const signingData = transactionSigningData(tx);
  const signature = makeSign(Nat.toNumber(tx.chainId || "0x1") * 2 + 35)(keccak256(signingData), privateKey);
  const rawTransaction = rlp.decode(signingData).slice(0,6).concat(decodeSignature(signature));
  return rlp.encode(rawTransaction);
}

const recoverTransaction = (rawTransaction) => {
  const values = rlp.decode(rawTransaction);
  const signature = encodeSignature(values.slice(6,9));
  const recovery = Bytes.toNumber(values[6]);
  const extraData = recovery < 35 ? [] : [Bytes.fromNumber((recovery - 35) >> 1), "0x", "0x"]
  const signingData = values.slice(0,6).concat(extraData);
  const signingDataHex = rlp.encode(signingData);
  return recover(keccak256(signingDataHex), signature);
}

module.exports = { 
  create,
  toChecksum,
  fromPrivate,
  sign,
  recover,
  signTransaction,
  recoverTransaction,
  transactionSigningData,
  encodeSignature,
  decodeSignature
}
