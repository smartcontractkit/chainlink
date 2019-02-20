import reducer from 'connectors/redux/reducers'
import {
  REQUEST_CREATE,
  RECEIVE_CREATE_SUCCESS,
  RECEIVE_CREATE_ERROR
} from 'actions'

describe('connectors/reducers/create', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.create).toEqual({
      networkError: false
    })
  })

  it('REQUEST_CREATE disables the network error', () => {
    const action = { type: REQUEST_CREATE }
    const state = reducer(undefined, action)

    expect(state.create.networkError).toEqual(false)
  })

  describe('RECEIVE_CREATE_SUCCESS', () => {
    it('assigns correct object and sets networkError to false', () => {
      const previousState = { create: { networkError: true } }
      const action = {
        type: RECEIVE_CREATE_SUCCESS,
        response: { successful: 'success message' }
      }
      const state = reducer(previousState, action)

      expect(state.create.networkError).toEqual(false)
    })
  })

  describe('RECEIVE_CREATE_ERROR', () => {
    it("does nothing because that's handled by global errors", () => {
      const previousState = { create: {} }
      const error = { errors: [{ detail: 'errored' }] }
      const action = { type: RECEIVE_CREATE_ERROR, error: error }
      const state = reducer(previousState, action)
      expect(state.create.errors).toEqual(undefined)
    })
  })
})
