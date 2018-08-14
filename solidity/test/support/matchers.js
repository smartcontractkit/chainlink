export const assertBigNum = (a, b) => assert(
  a.equals(b),
  `payment ${a} is not ${b}`
)
