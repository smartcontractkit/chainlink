import { decamelizeKeys } from 'humps'

export default (ethBalance, linkBalance) => {
  return decamelizeKeys({
    data: {
      id: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
      type: 'accountBalances',
      attributes: {
        address: '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
        ethBalance: ethBalance,
        linkBalance: linkBalance,
      },
    },
  })
}
