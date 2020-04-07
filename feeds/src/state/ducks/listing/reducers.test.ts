import { partialAsFull } from '../../../../../tools/ts-helpers/src'
import { setAnswers } from './actions'
import reducer, { INITIAL_STATE, ListingAnswer } from './reducers'

describe('state/ducks/listing/reducers', () => {
  describe('SET_ANSWERS', () => {
    it('should replace answers', () => {
      const data = [
        partialAsFull<ListingAnswer>({
          answer: 'answer',
        }),
      ]
      const action = setAnswers(data)
      const state = reducer(INITIAL_STATE, action)

      expect(state.answers).toEqual(data)
    })
  })
})
