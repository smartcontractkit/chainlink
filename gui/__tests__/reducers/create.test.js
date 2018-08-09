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
      errors: [],
      successMessage: {},
      networkError: false
    })
  })

  it('REQUEST_CREATE disables the network error', () => {
    const action = {type: REQUEST_CREATE}
    const state = reducer(undefined, action)

    expect(state.create.networkError).toEqual(false)
  })

  describe('RECEIVE_CREATE_SUCCESS', () => {
    it('assigns correct object and sets networkError to false', () => {
      const previousState = { create: { networkError: true } }
      const action = {type: RECEIVE_CREATE_SUCCESS, response: {successful: 'success message'}}
      const state = reducer(previousState, action)

      expect(state.create.successMessage).toEqual({successful: 'success message'})
      expect(state.create.networkError).toEqual(false)
    })
  })

  describe('RECEIVE_CREATE_ERROR', () => {
    it('assigns a network error', () => {
      const previousState = { create: { networkError: false } }
      const action = {type: RECEIVE_CREATE_ERROR, networkError: true}
      const state = reducer(previousState, action)

      expect(state.create.networkError).toEqual(true)
    })
  })
})
