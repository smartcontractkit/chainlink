const bigNum = web3.utils.toBN

// Throws if a and b are not equal, as BN's
export const assertBigNum = (a, b, failureMessage) => assert(
  bigNum(a).eq(bigNum(b)),
  `BigNum ${bigNum(a)} is not ${bigNum(b)}` + (failureMessage ? ': ' + failureMessage : '')
)
