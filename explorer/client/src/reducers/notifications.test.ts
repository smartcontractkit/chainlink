import reducer, { INITIAL_STATE } from '../reducers'
import { FetchAdminSigninErrorAction } from '../reducers/actions'

describe('reducers/notifications', () => {
  describe('FETCH_ADMIN_SIGNIN_ERROR', () => {
    it('adds a notification for invalid credentials', () => {
      const action: FetchAdminSigninErrorAction = {
        type: 'FETCH_ADMIN_SIGNIN_ERROR',
        errors: [{ status: 401, detail: 'Unauthorized' }],
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual([
        'Invalid username and password',
      ])
    })

    it("adds a notification when the server can't process the request", () => {
      const action: FetchAdminSigninErrorAction = {
        type: 'FETCH_ADMIN_SIGNIN_ERROR',
        errors: [{ status: 500, detail: 'Internal Server Error' }],
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual([
        'Error processing your request. Please ensure your connection is active and try again',
      ])
    })

    it('adds a notification with a pass through message by default', () => {
      const action: FetchAdminSigninErrorAction = {
        type: 'FETCH_ADMIN_SIGNIN_ERROR',
        errors: [{ status: 404, detail: 'Not Found' }],
      }
      const state = reducer(INITIAL_STATE, action)

      expect(state.notifications.errors).toEqual(['Not Found'])
    })
  })
})
