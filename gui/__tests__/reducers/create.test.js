import reducer from 'connectors/redux/reducers'
import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR
} from 'actions'

describe('create reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.create).toEqual({
      fetching: false,
      errors: [],
      successMessage: {},
      networkError: false
    })
  })

  it('REQUEST_CREATE starts fetching and disables the network error', () => {
    const action = {type: REQUEST_CREATE}
    const state = reducer(undefined, action)

    expect(state.create.fetching).toEqual(true)
    expect(state.create.networkError).toEqual(false)
  })

  describe('RECEIVE_CREATE_SUCCESS', () => {
    it('stops fetching', () => {
      const previousState = {
        create: {
          fetching: true,
          networkError: true
        }
      }
      const action = {type: RECEIVE_CREATE_SUCCESS, response: {successful: 'success message'}}
      const state = reducer(previousState, action)

      expect(state.create.fetching).toEqual(false)
      expect(state.create.successMessage).toEqual({successful: 'success message'})
      expect(state.create.networkError).toEqual(false)
    })
  })

  describe('RECEIVE_CREATE_ERROR', () => {
    it('stops fetching and assigns a network error', () => {
      const previousState = {
        create: {
          fetching: true,
          networkError: false
        }
      }
      const action = {type: RECEIVE_CREATE_ERROR, networkError: true}
      const state = reducer(previousState, action)

      expect(state.create.fetching).toEqual(false)
      expect(state.create.networkError).toEqual(true)
    })
  })
})
