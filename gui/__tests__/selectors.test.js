import {
  jobsSelector,
  jobSpecSelector,
  jobRunsSelector,
  jobRunsCountSelector,
  bridgeSelector,
  bridgesSelector,
  configsSelector,
  isFetchingSelector
} from 'selectors'

describe('selectors', () => {
  describe('jobsSelector', () => {
    it('returns the jobs in the current page', () => {
      const state = {
        jobs: {
          currentPage: ['jobA', 'jobB'],
          items: {
            jobA: {id: 'jobA'},
            jobB: {id: 'jobB'}
          }
        }
      }
      const jobs = jobsSelector(state)

      expect(jobs).toEqual([
        {id: 'jobA'},
        {id: 'jobB'}
      ])
    })

    it('excludes job items that are not present', () => {
      const state = {
        jobs: {
          currentPage: ['jobA', 'jobB'],
          items: {
            jobA: {id: 'jobA'}
          }
        }
      }
      const jobs = jobsSelector(state)

      expect(jobs).toEqual([
        {id: 'jobA'}
      ])
    })
  })

  describe('jobSpecSelector', () => {
    it('returns the job item for the given id and undefined otherwise', () => {
      const state = {
        jobs: {
          items: {
            jobA: {id: 'jobA'}
          }
        }
      }

      expect(jobSpecSelector(state, 'jobA')).toEqual({id: 'jobA'})
      expect(jobSpecSelector(state, 'joba')).toBeUndefined()
    })
  })

  describe('jobRunsSelector', () => {
    it('returns the job runs for the given job spec id', () => {
      const state = {
        jobRuns: {
          currentPage: ['runA', 'runB'],
          items: {
            'runA': {id: 'runA'},
            'runB': {id: 'runB'},
            'runC': {id: 'runC'}
          }
        }
      }
      const runs = jobRunsSelector(state)

      expect(runs).toEqual([
        {id: 'runA'},
        {id: 'runB'}
      ])
    })

    it('returns an empty array when the currentPage is empty', () => {
      const state = {
        jobRuns: {
          currentPage: [],
          items: {
            'runA': {id: 'runA'}
          }
        }
      }
      const runs = jobRunsSelector(state)

      expect(runs).toEqual([])
    })

    it('excludes job runs that do not have items', () => {
      const state = {
        jobRuns: {
          currentPage: ['runA', 'runB'],
          items: {
            'runA': {id: 'runA'}
          }
        }
      }
      const runs = jobRunsSelector(state)

      expect(runs).toEqual([
        {id: 'runA'}
      ])
    })
  })

  describe('jobRunsCountSelector', () => {
    it('returns the number of runs for the job', () => {
      const state = {
        jobs: {
          items: {
            jobA: {id: 'jobA', runsCount: 6}
          }
        }
      }

      expect(jobRunsCountSelector(state, 'jobA')).toEqual(6)
    })

    it('returns the number 0 when the job doesn\'t exist', () => {
      const state = {
        jobs: {
          items: {}
        }
      }

      expect(jobRunsCountSelector(state, 'jobA')).toEqual(0)
    })
  })

  describe('bridgeSelector', () => {
    it('returns the bridge with the given id', () => {
      const state = {
        bridges: {
          items: {
            a: {name: 'A'},
            b: {name: 'B'}
          }
        }
      }

      expect(bridgeSelector(state, 'a')).toEqual({name: 'A'})
    })
  })

  describe('bridgesSelector', () => {
    it('returns the current page of bridges', () => {
      const state = {
        bridges: {
          items: {
            a: {name: 'A'},
            b: {name: 'B'},
            c: {name: 'C'}
          },
          currentPage: ['c', 'a']
        }
      }

      expect(bridgesSelector(state, 'a')).toEqual([
        {name: 'C'},
        {name: 'A'}
      ])
    })

    it('does not return items that cannot be found', () => {
      const state = {
        bridges: {
          items: {
            a: {name: 'A'},
            b: {name: 'B'},
            c: {name: 'C'}
          },
          currentPage: ['C', 'A']
        }
      }

      expect(bridgesSelector(state, 'a')).toEqual([])
    })
  })

  describe('configsSelector', () => {
    it('returns a tuple per key/value pair', () => {
      const state = {
        configuration: {
          config: {
            camelCased: 'value',
            key: 'value'
          }
        }
      }

      let expectation = [
        ['CAMEL_CASED', 'value'],
        ['KEY', 'value']
      ]
      expect(configsSelector(state)).toEqual(expectation)
    })
  })

  describe('isFetchingSelector', () => {
    it('is true when count > 0', () => {
      const state = {fetching: {count: 1}}

      expect(isFetchingSelector(state)).toEqual(true)
    })

    it('is false when count = 0', () => {
      const state = {fetching: {count: 0}}

      expect(isFetchingSelector(state)).toEqual(false)
    })
  })
})
