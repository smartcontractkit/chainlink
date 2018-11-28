export const assertBigNum = (a, b) => assert(
  a.equals(b),
  `BigNum ${a} is not ${b}`
)
