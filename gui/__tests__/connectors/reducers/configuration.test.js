import reducer from 'connectors/redux/reducers'
import {
  REQUEST_CONFIGURATION,
  RECEIVE_CONFIGURATION_SUCCESS,
  RECEIVE_CONFIGURATION_ERROR
} from 'actions'

describe('connectors/reducers/configuration', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.configuration).toEqual({
      config: {},
      networkError: false
    })
  })

  it('REQUEST_CONFIGURATION disables the network error', () => {
    const action = { type: REQUEST_CONFIGURATION }
    const state = reducer(undefined, action)

    expect(state.configuration.networkError).toEqual(false)
  })

  it('RECEIVE_CONFIGURATION_SUCCESS assigns the config', () => {
    const previousState = {
      configuration: {
        networkError: true
      }
    }
    let configMap = { singer: 'bob' }
    const action = {
      type: RECEIVE_CONFIGURATION_SUCCESS,
      config: configMap
    }
    const state = reducer(previousState, action)

    expect(state.configuration.config).toEqual(configMap)
    expect(state.configuration.networkError).toEqual(false)
  })

  it('RECEIVE_CONFIGURATION_ERROR assigns a network error', () => {
    const previousState = {
      configuration: {
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_CONFIGURATION_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.configuration.networkError).toEqual(true)
  })
})
