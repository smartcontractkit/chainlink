import { partialAsFull } from '@chainlink/ts-helpers'
import { FeedConfig } from 'config'
import reducer, { INITIAL_STATE } from './reducers'
import {
  fetchFeedsBegin,
  fetchFeedsSuccess,
  fetchFeedsError,
  fetchAnswerSuccess,
  fetchAnswerTimestampSuccess,
} from './actions'

describe('state/ducks/listing/reducers', () => {
  describe('FETCH_FEEDS_*', () => {
    it('toggles loadingFeeds when the request starts & finishes', () => {
      let state

      const beginAction = fetchFeedsBegin()
      state = reducer(INITIAL_STATE, beginAction)
      expect(state.loadingFeeds).toEqual(true)

      const payload: FeedConfig[] = []
      const successAction = fetchFeedsSuccess(payload)
      state = reducer(state, successAction)
      expect(state.loadingFeeds).toEqual(false)

      state = reducer(state, beginAction)
      expect(state.loadingFeeds).toEqual(true)

      const errorAction = fetchFeedsError(new Error('Not Found'))
      state = reducer(state, errorAction)
      expect(state.loadingFeeds).toEqual(false)
    })
  })

  describe('FETCH_FEEDS_SUCCESS', () => {
    it('indexes listed feeds by contract address', () => {
      const feedA = partialAsFull<FeedConfig>({
        contractAddress: 'A',
        listing: true,
      })
      const feedB = partialAsFull<FeedConfig>({
        contractAddress: 'B',
        listing: false,
      })
      const payload: FeedConfig[] = [feedA, feedB]
      const successAction = fetchFeedsSuccess(payload)

      const state = reducer(INITIAL_STATE, successAction)
      expect(state.feedItems).toEqual({ A: feedA })
    })

    it('stores the order of listed contract addresses', () => {
      const feedA = partialAsFull<FeedConfig>({
        contractAddress: 'A',
        listing: true,
      })
      const feedB = partialAsFull<FeedConfig>({
        contractAddress: 'B',
        listing: true,
      })
      const payload: FeedConfig[] = [feedA, feedB]
      const successAction = fetchFeedsSuccess(payload)

      const state = reducer(INITIAL_STATE, successAction)
      expect(state.feedOrder).toEqual(['A', 'B'])
    })
  })

  describe('FETCH_ANSWER_SUCCESS', () => {
    it('should replace answers', () => {
      const config = partialAsFull<FeedConfig>({ contractAddress: 'A' })
      const payload = { answer: 'answer', config }
      const action = fetchAnswerSuccess(payload)
      const state = reducer(INITIAL_STATE, action)

      expect(state.answers).toEqual({ A: 'answer' })
    })
  })

  describe('FETCH_ANSWER_TIMESTAMP_SUCCESS', () => {
    it('should replace answers timestamp', () => {
      const config = partialAsFull<FeedConfig>({ contractAddress: 'A' })
      const payload = { timestamp: 123, config }
      const action = fetchAnswerTimestampSuccess(payload)
      const state = reducer(INITIAL_STATE, action)

      expect(state.answersTimestamp).toEqual({ A: 123 })
    })
  })
})
