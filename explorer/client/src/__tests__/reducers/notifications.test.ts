import * as jsonapi from '@chainlink/json-api-client'
import { partialAsFull } from '@chainlink/ts-test-helpers'
import reducer, { INITIAL_STATE } from '../../reducers'

describe('reducers/jobRuns', () => {
  describe('FETCH_ADMIN_SIGNIN_ERROR', () => {
    it('adds a notification for AuthenticationError', () => {
      const response = partialAsFull<Response>({})
      const action = {
        type: 'FETCH_ADMIN_SIGNIN_ERROR',
        error: new jsonapi.AuthenticationError(response),
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual([
        'Invalid username and password.',
      ])
    })
  })
})
