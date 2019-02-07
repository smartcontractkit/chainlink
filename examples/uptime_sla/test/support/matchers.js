const bigNum = number => web3.utils.toBN(number)

// Throws if a and b are not equal, as BN's
export const assertBigNum = (a, b, failureMessage) => assert(
  bigNum(a).eq(bigNum(b)),
  `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : '')
)
