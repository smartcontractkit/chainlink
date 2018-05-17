import { decamelizeKeys } from 'humps'

export default (ethBalance, linkBalance) => {
  return decamelizeKeys({
    data: {
      attributes: {
        ethBalance: ethBalance,
        linkBalance: linkBalance
      }
    }
  })
}
