import reducer from 'connectors/redux/reducers'
import {
  REQUEST_CONFIGURATION,
  RECEIVE_CONFIGURATION_SUCCESS,
  RECEIVE_CONFIGURATION_ERROR
} from 'actions'

describe('configuration reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.configuration).toEqual({
      config: {},
      fetching: false,
      networkError: false
    })
  })

  it('REQUEST_CONFIGURATION starts fetching and disables the network error', () => {
    const action = {type: REQUEST_CONFIGURATION}
    const state = reducer(undefined, action)

    expect(state.configuration.fetching).toEqual(true)
    expect(state.configuration.networkError).toEqual(false)
  })

  it('RECEIVE_CONFIGURATION_SUCCESS stops fetching and assigns the config', () => {
    const previousState = {
      configuration: {
        fetching: true,
        networkError: true
      }
    }
    let configMap = {singer: 'bob'}
    const action = {
      type: RECEIVE_CONFIGURATION_SUCCESS,
      config: configMap
    }
    const state = reducer(previousState, action)

    expect(state.configuration.config).toEqual(configMap)
    expect(state.configuration.fetching).toEqual(false)
    expect(state.configuration.networkError).toEqual(false)
  })

  it('RECEIVE_CONFIGURATION_ERROR stops fetching and assigns a network error', () => {
    const previousState = {
      configuration: {
        fetching: true,
        networkError: false
      }
    }
    const action = {
      type: RECEIVE_CONFIGURATION_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.configuration.fetching).toEqual(false)
    expect(state.configuration.networkError).toEqual(true)
  })
})
