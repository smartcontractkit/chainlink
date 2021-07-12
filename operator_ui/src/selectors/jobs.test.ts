import { partialAsFull } from 'support/test-helpers/partialAsFull'
import { INITIAL_STATE, AppState } from '../../src/reducers'
import jobsSelector from '../../src/selectors/jobs'

describe('selectors - jobs', () => {
  type JobsState = typeof INITIAL_STATE.jobs

  it('returns the jobs in the current page', () => {
    const jobsState = partialAsFull<JobsState>({
      currentPage: ['jobA', 'jobB'],
      items: {
        jobA: { id: 'jobA' },
        jobB: { id: 'jobB' },
      },
    })
    const state: Pick<AppState, 'jobs'> = {
      jobs: jobsState,
    }
    const jobs = jobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }, { id: 'jobB' }])
  })

  it('excludes job items that are not present', () => {
    const jobsState = partialAsFull<JobsState>({
      currentPage: ['jobA', 'jobB'],
      items: {
        jobA: { id: 'jobA' },
      },
    })
    const state: Pick<AppState, 'jobs'> = {
      jobs: jobsState,
    }
    const jobs = jobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }])
  })
})
