import * as presenters from 'core/store/presenters'

export const accountBalances = (
  accountBalances: Partial<presenters.AccountBalance>[],
) => {
  return {
    data: accountBalances.map((balance) => ({
      id: balance.address ?? '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
      type: 'eTHKeys',
      attributes: {
        address:
          balance.address ?? '0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f',
        ethBalance: balance.ethBalance ?? '0',
        linkBalance: balance.linkBalance ?? '0',
        isFunding: balance.isFunding ?? false,
        createdAt: balance.createdAt ?? new Date().toISOString(),
      },
    })),
  }
}
