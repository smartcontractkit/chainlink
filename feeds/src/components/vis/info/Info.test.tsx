import { partialAsFull } from '@chainlink/ts-helpers'
import React from 'react'
import '@testing-library/jest-dom/extend-expect'
import { render } from '@testing-library/react'
import { FeedConfig } from 'config'
import Info from './Info'

const feed = partialAsFull<FeedConfig>({
  proxyAddress: 'ProxyAddress1',
  contractAddress: 'contractAddress1',
})

const props = {
  latestAnswer: '1',
  latestRequestTimestamp: 1,
  minimumAnswers: 1,
  oracleAnswers: undefined,
  oracleList: [],
  latestAnswerTimestamp: undefined,
  pendingRoundId: undefined,
}

describe('Info', () => {
  describe('Proxy address', () => {
    it('renders proxy and aggregator contract addresses when there is a proxy address', () => {
      const { container } = render(<Info config={feed} {...props} />)

      expect(container).toHaveTextContent(`Feed: ${feed.proxyAddress}`)
      expect(container).toHaveTextContent(`Aggregator: ${feed.contractAddress}`)
    })

    it('renders only aggregator contract address by default', () => {
      const feedWithoutProxy = { ...feed }
      delete feedWithoutProxy.proxyAddress

      const { container } = render(
        <Info config={feedWithoutProxy} {...props} />,
      )

      expect(container).not.toHaveTextContent(`Feed: ${feed.proxyAddress}`)
      expect(container).toHaveTextContent(`Aggregator: ${feed.contractAddress}`)
    })

    it('renders only aggregator contract address when there is no proxy address (blank)', () => {
      const feedWithoutProxy = partialAsFull<FeedConfig>({
        proxyAddress: '',
        contractAddress: 'contractAddress1',
      })

      const { container } = render(
        <Info config={feedWithoutProxy} {...props} />,
      )

      expect(container).not.toHaveTextContent(`Feed: ${feed.proxyAddress}`)
      expect(container).toHaveTextContent(`Aggregator: ${feed.contractAddress}`)
    })
  })
})
