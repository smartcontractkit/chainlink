import { ROPSTEN_ID, MAINNET_ID } from 'utils/'

export const networkName = networkId => {
  switch (networkId) {
    case MAINNET_ID:
      return 'mainnet'
    case ROPSTEN_ID:
      return 'ropsten'
    default:
      return 'mainnet'
  }
}
