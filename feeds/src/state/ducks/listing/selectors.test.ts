import { partialAsFull } from '@chainlink/ts-helpers/src'
import { feedGroups, answer } from './selectors'
import { FeedConfig } from 'config'
import { AppState } from 'state/reducers'
import { HealthCheck } from 'state/ducks/listing/reducers'

const feedA = partialAsFull<FeedConfig>({
  contractAddress: 'A',
  pair: ['BTC', 'USD'],
})
const feedB = partialAsFull<FeedConfig>({
  contractAddress: 'B',
  pair: ['BTC', 'ETH'],
})
const feedC = partialAsFull<FeedConfig>({
  contractAddress: 'C',
  pair: ['FTSE', 'GBP'],
})
const listedFeeds = [feedA, feedB, feedC]

describe('state/ducks/listing/selectors', () => {
  describe('feedGroups', () => {
    it('returns a Fiat & ETH listing group', () => {
      const selected = feedGroups.resultFunc([])
      expect(selected).toHaveLength(2)
      expect(selected[0].name).toMatch('Fiat')
      expect(selected[1].name).toMatch('ETH')
    })

    it('returns a list of feed configs grouped by quote asset', () => {
      const selected = feedGroups.resultFunc(listedFeeds)

      const group1 = selected[0]
      expect(group1.feeds.length).toEqual(2)
      expect(group1.feeds[0].contractAddress).toEqual(feedA.contractAddress)
      expect(group1.feeds[1].contractAddress).toEqual(feedC.contractAddress)

      const group2 = selected[1]
      expect(group2.feeds.length).toEqual(1)
      expect(group2.feeds[0].contractAddress).toEqual(feedB.contractAddress)
    })
  })

  describe('answer', () => {
    const feedA = partialAsFull<FeedConfig>({ contractAddress: 'A' })
    const answers: Record<FeedConfig['contractAddress'], string> = {
      [feedA.contractAddress]: '10.1',
    }
    const healthChecks: Record<string, HealthCheck> = {}
    const listingState = partialAsFull<AppState['listing']>({
      answers,
      healthChecks,
    })
    const state = partialAsFull<AppState>({
      listing: listingState,
    })

    it('returns the answer for the contract', () => {
      expect(answer(state, 'A')).toEqual('10.1')
    })

    it('returns undefined when there is no answer for the contract', () => {
      expect(answer(state, 'B')).toBeUndefined()
    })
  })
})
