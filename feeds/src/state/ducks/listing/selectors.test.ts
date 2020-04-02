import { groups, answer } from './selectors'
import { partialAsFull } from '@chainlink/ts-helpers/src'
import { FeedConfig } from 'feeds'
import { AppState } from 'state/reducers'
import { ListingAnswer } from 'state/ducks/listing/operations'
import { HealthCheck } from 'state/ducks/listing/reducers'

const feed1 = partialAsFull<FeedConfig>({
  contractAddress: 'A',
  listing: true,
  pair: ['BTC', 'USD'],
})
const feed2 = partialAsFull<FeedConfig>({
  contractAddress: 'B',
  listing: true,
  pair: ['BTC', 'ETH'],
})
const feed3 = partialAsFull<FeedConfig>({
  contractAddress: 'C',
  listing: false,
  pair: ['XBT', 'USD'],
})
const feeds = [feed1, feed2, feed3]

describe('state/ducks/listing/selectors', () => {
  describe('groups', () => {
    it('returns a USD & ETH listing group', () => {
      const selected = groups.resultFunc([])
      expect(selected).toHaveLength(2)
      expect(selected[0].name).toMatch('USD')
      expect(selected[1].name).toMatch('ETH')
    })

    it('returns listed feed configs grouped by quote asset', () => {
      const selected = groups.resultFunc(feeds)

      const group1 = selected[0]
      expect(group1.feeds.length).toEqual(1)
      expect(group1.feeds[0].contractAddress).toEqual(feed1.contractAddress)

      const group2 = selected[1]
      expect(group2.feeds.length).toEqual(1)
      expect(group2.feeds[0].contractAddress).toEqual(feed2.contractAddress)
    })
  })

  describe('answer', () => {
    const feedA = partialAsFull<FeedConfig>({ contractAddress: 'A' })
    const answerA: ListingAnswer = { answer: '10.1', config: feedA }
    const answers: ListingAnswer[] = [answerA]
    const healthChecks: Record<string, HealthCheck> = {}
    const listingState = partialAsFull<AppState['listing']>({
      answers,
      healthChecks,
    })
    const state = partialAsFull<AppState>({
      listing: listingState,
    })

    it('returns the answer for the contract', () => {
      expect(answer(state, 'A')).toEqual(answerA)
    })

    it('returns undefined when there is no answer for the contract', () => {
      expect(answer(state, 'B')).toBeUndefined()
    })
  })
})
