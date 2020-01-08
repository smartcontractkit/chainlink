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
    decimalPlaces: 9,
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

  //Aave

  {
    contractAddress: '0x1EeaF25f2ECbcAf204ECADc8Db7B0db9DA845327',
    name: 'LEND / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'LEND',
    path: 'lend-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xF17BE6Ea8506980D0c6832EdE67e337F1684Ae02',
    name: 'AMPL / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'AMPL',
    path: 'ampl-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0x0133Aa47B6197D0BA090Bf2CD96626Eb71fFd13c',
    name: 'BTC / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'BTC',
    path: 'btc-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xda3d675d50ff6c555973c4f0424964e1f6a4e7d3',
    name: 'MKR / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'MKR',
    path: 'mkr-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xc89c4ed8f52Bb17314022f6c0dCB26210C905C97',
    name: 'MANA / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'MANA',
    path: 'mana-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xd0e785973390fF8E77a83961efDb4F271E6B8152',
    name: 'KNC / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'KNC',
    path: 'knc-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xeCfA53A8bdA4F0c4dd39c55CC8deF3757aCFDD07',
    name: 'LINK / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'LINK',
    path: 'link-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xdE54467873c3BCAA76421061036053e371721708',
    name: 'USDC / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'USDC',
    path: 'usdc-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xb8b513d9cf440C1b6f5C7142120d611C94fC220c',
    name: 'REP / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'REP',
    path: 'rep-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xA0F9D94f060836756FFC84Db4C78d097cA8C23E8',
    name: 'ZRX / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'ZRX',
    path: 'zrx-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0x9b4e2579895efa2b4765063310Dc4109a7641129',
    name: 'BAT / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'BAT',
    path: 'bat-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0x037E8F2125bF532F3e228991e051c8A7253B642c',
    name: 'DAI / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'DAI',
    path: 'dai-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0x73ead35fd6A572EF763B13Be65a9db96f7643577',
    name: 'TUSD / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'TUSD',
    path: 'tusd-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xa874fe207DF445ff19E7482C746C4D3fD0CB9AcE',
    name: 'USDT / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'USDT',
    path: 'usdt-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0x6d626Ff97f0E89F6f983dE425dc5B24A18DE26Ea',
    name: 'SUSD / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'SUSD',
    path: 'susd-eth',
    multiply: '1000000000000000000',
    decimalPlaces: 9,
    history: false,
  },
  {
    contractAddress: '0xE23d1142dE4E83C08bb048bcab54d50907390828',
    name: 'SNX / ETH aggregation',
    valuePrefix: 'Ξ',
    answerName: 'SNX',
    path: 'snx-eth',
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
