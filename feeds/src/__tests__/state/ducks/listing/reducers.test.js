import reducer from 'state/ducks/listing'
import { initialState } from 'state/ducks/listing/reducers'
import { setAnswers } from 'state/ducks/listing/actions'

describe('state/ducks/listing/reducers', () => {
  describe('SET_ANSWERS', () => {
    it('should replace answers', () => {
      expect(initialState.answers).toEqual(null)
      const data = [
        {
          answer: 'answer',
        },
      ]
      const action = setAnswers(data)
      const state = reducer(initialState, action)
      expect(state.answers).toEqual(data)
    })
  })
})
