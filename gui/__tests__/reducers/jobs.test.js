import reducer from 'connectors/redux/reducers'
import {
  REQUEST_JOBS,
  RECEIVE_JOBS_SUCCESS,
  RECEIVE_JOBS_ERROR
} from 'actions'

describe('jobs reducer', () => {
  it('should return the initial state', () => {
    const state = reducer(undefined, {})

    expect(state.jobs).toEqual({
      items: {},
      currentPage: [],
      count: 0,
      fetching: false,
      networkError: false
    })
  })

  it('REQUEST_JOBS starts fetching and disables the network error', () => {
    const action = {type: REQUEST_JOBS}
    const state = reducer(undefined, action)

    expect(state.jobs.fetching).toEqual(true)
    expect(state.jobs.networkError).toEqual(false)
  })

  it('RECEIVE_JOBS_SUCCESS stores the job items and the current page', () => {
    const action = {
      type: RECEIVE_JOBS_SUCCESS,
      items: [{id: 'a'}, {id: 'b'}]
    }
    const state = reducer(undefined, action)

    expect(state.jobs.items).toEqual({
      'a': {id: 'a'},
      'b': {id: 'b'}
    })
    expect(state.jobs.currentPage).toEqual(['a', 'b'])
    expect(state.jobs.fetching).toEqual(false)
    expect(state.jobs.networkError).toEqual(false)
  })

  it('RECEIVE_JOBS_ERROR stops fetching and updates the network error', () => {
    const previousState = {
      jobs: {networkError: false, fetching: true}
    }
    const action = {
      type: RECEIVE_JOBS_ERROR,
      networkError: true
    }
    const state = reducer(previousState, action)

    expect(state.jobs.fetching).toEqual(false)
    expect(state.jobs.networkError).toEqual(true)
  })
})
