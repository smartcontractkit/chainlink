import {
  jobsSelector,
  jobSpecSelector,
  jobRunsSelector,
  jobRunsCountSelector,
  latestJobRunsSelector
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

  describe('jobRunsSelectors', () => {
    it('returns the job runs for the given job spec id', () => {
      const state = {
        jobs: {
          items: {
            jobA: {id: 'jobA', runs: ['runA', 'runB']}
          }
        },
        jobRuns: {
          items: {
            'runA': {id: 'runA'},
            'runB': {id: 'runB'}
          }
        }
      }
      const runs = jobRunsSelector(state, 'jobA')

      expect(runs).toEqual([
        {id: 'runA'},
        {id: 'runB'}
      ])
    })

    it('returns an empty array when the job does not exist', () => {
      const state = {
        jobs: {
          items: {}
        }
      }
      const runs = jobRunsSelector(state, 'jobA')

      expect(runs).toEqual([])
    })

    it('returns an empty array when the job does not have the runs attribute', () => {
      const state = {
        jobs: {
          items: {
            'jobA': {id: 'jobA'}
          }
        }
      }
      const runs = jobRunsSelector(state, 'jobA')

      expect(runs).toEqual([])
    })

    it('excludes job runs that do not have items', () => {
      const state = {
        jobs: {
          items: {
            jobA: {id: 'jobA', runs: ['runA', 'runB']}
          }
        },
        jobRuns: {
          items: {
            'runA': {id: 'runA'}
          }
        }
      }
      const runs = jobRunsSelector(state, 'jobA')

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
            jobA: {id: 'jobA', runs: ['runA', 'runB', 'runC', 'runD', 'runE', 'runF']}
          }
        },
        jobRuns: {
          items: {
            'runA': {id: 'runA'},
            'runB': {id: 'runB'},
            'runC': {id: 'runC'},
            'runD': {id: 'runD'},
            'runE': {id: 'runE'},
            'runF': {id: 'runF'}
          }
        }
      }

      expect(jobRunsCountSelector(state, 'jobA')).toEqual(6)
    })
  })

  describe('latestJobRunsSelector', () => {
    it('returns the 5 latest runs by creation date', () => {
      const state = {
        jobs: {
          items: {
            jobA: {id: 'jobA', runs: ['runA', 'runB', 'runC']}
          }
        },
        jobRuns: {
          items: {
            'runA': {id: 'runA', createdAt: '2018-05-01T16:54:16.255900955-07:00'},
            'runB': {id: 'runB', createdAt: '2018-05-02T16:54:16.255900955-07:00'},
            'runC': {id: 'runC', createdAt: '2018-05-03T16:54:16.255900955-07:00'}
          }
        }
      }
      const take = 2
      const runs = latestJobRunsSelector(state, 'jobA', take)

      expect(runs).toEqual([
        {id: 'runC', createdAt: '2018-05-03T16:54:16.255900955-07:00'},
        {id: 'runB', createdAt: '2018-05-02T16:54:16.255900955-07:00'}
      ])
    })
  })
})
