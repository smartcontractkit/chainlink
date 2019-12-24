import React from 'react'
import { Redirect } from 'react-router-dom'
import { withRouter } from 'react-router'

const DEFAULT_OPTIONS = {
  contractVersion: 2,
  network: 'mainnet',
  history: true,
  bollinger: false,
  decimalPlaces: 3,
  multiply: '100000000',
}

const MAINNET_CONTRACTS = [
  {
    contractAddress: '0x05cf62c4ba0ccea3da680f9a8744ac51116d6231',
    name: 'AUD / USD aggregation',
    valuePrefix: '$',
    answerName: 'AUD',
    counter: 3600,
    path: 'aud-usd',
  },
  {
    contractAddress: '0x25Fa978ea1a7dc9bDc33a2959B9053EaE57169B5',
    name: 'EUR / USD aggregation',
    valuePrefix: '$',
    answerName: 'EUR',
    counter: 3600,
    path: 'eur-usd',
  },
  {
    contractAddress: '0x02d5c618dbc591544b19d0bf13543c0728a3c4ec',
    name: 'CHF / USD aggregation',
    valuePrefix: '$',
    answerName: 'CHF',
    counter: 3600,
    path: 'chf-usd',
  },
  {
    contractAddress: '0x151445852b0cfdf6a4cc81440f2af99176e8ad08',
    name: 'GBP / USD aggregation',
    valuePrefix: '$',
    answerName: 'GBP',
    counter: 3600,
    path: 'gbp-usd',
  },
  {
    contractAddress: '0xe1407BfAa6B5965BAd1C9f38316A3b655A09d8A6',
    name: 'JPY / USD aggregation',
    valuePrefix: '$',
    answerName: 'JPY',
    counter: 3600,
    path: 'jpy-usd',
  },
  {
    contractAddress: '0x8946a183bfafa95becf57c5e08fe5b7654d2807b',
    name: 'XAG / USD aggregation',
    valuePrefix: '$',
    answerName: 'XAG',
    counter: 3600,
    path: 'xag-usd',
  },
  {
    contractAddress: '0xafce0c7b7fe3425adb3871eae5c0ec6d93e01935',
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
  {
    contractAddress: '0x8770Afe90c52Fd117f29192866DE705F63e59407',
    name: 'LRC / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'LRC',
    path: 'lrc-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
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
