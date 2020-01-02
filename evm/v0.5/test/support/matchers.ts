const bigNum = (num: number) => web3.utils.toBN(num)

// Throws if a and b are not equal, as BN's
//@ts-ignore
export const assertBigNum = (a, b, failureMessage) =>
  assert(
    bigNum(a).eq(bigNum(b)),
    `BigNum ${a} is not ${b}` + (failureMessage ? ': ' + failureMessage : ''),
  )
