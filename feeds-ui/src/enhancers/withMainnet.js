import React from 'react'
import { Redirect } from 'react-router-dom'
import { withRouter } from 'react-router'

const DEFAULT_OPTIONS = {
  contractVersion: 2,
  network: 'mainnet',
  history: true,
  bollinger: false,
}

const MAINNET_CONTRACTS = [
  {
    contractAddress: '0x9ab9ee9a5ac8aa4be93345877966b727777501b3',
    name: 'AUD / USD aggregation',
    valuePrefix: '$',
    answerName: 'AUD',
    counter: 3600,
    path: 'aud-usd',
  },
  {
    contractAddress: '0xa18e0bf319a1316984c9fc9aba850e2a915cce0e',
    name: 'EUR / USD aggregation',
    valuePrefix: '$',
    answerName: 'EUR',
    counter: 3600,
    path: 'eur-usd',
  },
  {
    contractAddress: '0xb405aa39dca6ecf18691cd638a98c090c477dcf9',
    name: 'CHF / USD aggregation',
    valuePrefix: '$',
    answerName: 'CHF',
    counter: 3600,
    path: 'chf-usd',
  },
  {
    contractAddress: '0xf277717f64815b17aa687cce701fa154185db4d9',
    name: 'GBP / USD aggregation',
    valuePrefix: '$',
    answerName: 'GBP',
    counter: 3600,
    path: 'gbp-usd',
  },
  {
    contractAddress: '0x8C4a5746d70C1aCeD6c3feB0915f59e3faBb18c3',
    name: 'JPY / USD aggregation',
    valuePrefix: '$',
    answerName: 'JPY',
    counter: 3600,
    path: 'jpy-usd',
  },
  {
    contractAddress: '0xe1ec292422194cfa6bb6089566af3019803b1af2',
    name: 'XAG / USD aggregation',
    valuePrefix: '$',
    answerName: 'XAG',
    counter: 3600,
    path: 'xag-usd',
  },
  {
    contractAddress: '0x121df6cfe2426e6c9ae1b273f2e38d0c362c03fa',
    name: 'XAU / USD aggregation',
    valuePrefix: '$',
    answerName: 'XAU',
    counter: 3600,
    path: 'xau-usd',
  },
  {
    contractAddress: '0xF5fff180082d6017036B771bA883025c654BC935',
    name: 'BTC / USD aggregation',
    valuePrefix: '$',
    answerName: 'BTC',
    counter: 600,
    path: 'btc-usd',
  },
  {
    contractVersion: 1,
    contractAddress: '0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9',
    name: 'ETH / USD aggregation',
    valuePrefix: '$',
    answerName: 'ETH',
    counter: 600,
    path: 'eth-usd',
  },
]

/**
 * withMainnet enhancer
 */
const withMainnet = BaseComponent => {
  const Mainnet = props => {
    const { params } = props.match
    const hasContract = MAINNET_CONTRACTS.filter(
      contract => contract.path.toLowerCase() === params.pair.toLowerCase(),
    )

    if (!hasContract.length) {
      return <Redirect to={'/'} />
    }

    const options = { ...DEFAULT_OPTIONS, ...hasContract[0] }

    return <BaseComponent {...props} options={options} />
  }

  return Mainnet
}

export default BaseComponent => {
  return withMainnet(withRouter(BaseComponent))
}
