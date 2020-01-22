import React from 'react'
import { Redirect } from 'react-router-dom'
import { withRouter } from 'react-router'

const DEFAULT_OPTIONS = {
  network: 'ropsten',
  history: true,
  bollinger: false,
}

const ROPSTEN_CONTRACTS = [
  {
    contractAddress: '0xEa7302472d58761582408F8512C24fc5bb8DBB70',
    name: 'BNB / USD aggregation',
    valuePrefix: '$',
    answerName: 'BNB',
    counter: 1800,
    path: 'bnb-usd',
  },
  {
    contractAddress: '0x1c44616CdB7FAe1ba69004ce6010248147CE019e',
    name: 'BTC / USD aggregation',
    valuePrefix: '$',
    answerName: 'BTC',
    counter: 3600,
    path: 'btc-usd',
  },
  {
    contractAddress: '0x46aD082e62D86089b7365320081685115F50d8B3',
    name: 'SNX / USD aggregation',
    valuePrefix: '$',
    answerName: 'SNX',
    counter: 1800,
    path: 'snx-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0x60289a8aa8ae8e2e2d03d5c1ad851a991cf29100',
    name: 'XTZ / USD aggregation',
    valuePrefix: '$',
    answerName: 'XTZ',
    counter: 1800,
    path: 'xtz-usd',
  },
  {
    contractAddress: '0x0Be00A19538Fac4BE07AC360C69378B870c412BF',
    name: 'ETH / USD aggregation',
    valuePrefix: '$',
    answerName: 'ETH',
    path: 'eth-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0xf230921e0bd70d5d3054aad25f8203e9853d5080',
    name: 'MKR / USD aggregation',
    valuePrefix: '$',
    answerName: 'MKR',
    counter: 1800,
    path: 'mkr-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0x4121c27b5bed6628b46493445108c34e130c996a',
    name: 'TRX / USD aggregation',
    valuePrefix: '$',
    answerName: 'TRX',
    counter: 1800,
    path: 'trx-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0x1b7ce2481149328c5e00efa6daa82de8e24f078b',
    name: 'AUD / USD aggregation',
    valuePrefix: '$',
    answerName: 'AUD',
    counter: 3600,
    path: 'aud-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0xc8bc999deab18feca1c5fbd6bffe9975ac396402',
    name: 'XAG / USD aggregation',
    valuePrefix: '$',
    answerName: 'XAG',
    counter: 3600,
    path: 'xag-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0xcd80a8f6915c78b3e65d30f94468547e021ccf9b',
    name: 'CHF / USD aggregation',
    valuePrefix: '$',
    answerName: 'CHF',
    counter: 3600,
    path: 'chf-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0x174754491a4ca333bf777387b0926bc8ecaf7f6e',
    name: 'GBP / USD aggregation',
    valuePrefix: '$',
    answerName: 'GBP',
    counter: 3600,
    path: 'gbp-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0xf45a5bb73124907e8c391c6a1001896f62f8f290',
    name: 'XAU / USD aggregation',
    valuePrefix: '$',
    answerName: 'XAU',
    counter: 3600,
    path: 'xau-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0x152cfa5d0e11ab0355179cd812035c2c64d750bd',
    name: 'EUR / USD aggregation',
    valuePrefix: '$',
    answerName: 'EUR',
    counter: 3600,
    path: 'eur-usd',
  },
  {
    contractVersion: 2,
    contractAddress: '0xe3153d946c958e334285f4aa93c6a3d8f5dfbff7',
    name: 'JPY / USD aggregation',
    valuePrefix: '$',
    answerName: 'JPY',
    counter: 3600,
    path: 'jpy-usd',
  },
]

/**
 * withRopsten enhancer
 */
const withRopsten = BaseComponent => {
  const Ropsten = props => {
    const { params } = props.match
    const hasContract = ROPSTEN_CONTRACTS.filter(
      contract => contract.path.toLowerCase() === params.pair.toLowerCase(),
    )

    if (!hasContract.length) {
      return <Redirect to={'/'} />
    }

    const options = { ...DEFAULT_OPTIONS, ...hasContract[0] }

    return <BaseComponent {...props} options={options} />
  }

  return Ropsten
}

export default BaseComponent => {
  return withRopsten(withRouter(BaseComponent))
}
