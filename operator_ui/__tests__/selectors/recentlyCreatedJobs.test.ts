import recentlyCreatedJobsSelector from 'selectors/recentlyCreatedJobs'

describe('selectors - jobs', () => {
  it('returns null when not loaded', () => {
    const state = {
      jobs: {
        recentlyCreated: null,
      },
    }
    const jobs = recentlyCreatedJobsSelector(state)

    expect(jobs).toEqual(null)
  })

  it('returns the job objects in items and excludes those not present', () => {
    const state = {
      jobs: {
        recentlyCreated: ['jobA', 'jobB', 'jobC'],
        items: {
          jobA: { id: 'jobA' },
          jobB: { id: 'jobB' },
        },
      },
    }
    const jobs = recentlyCreatedJobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }, { id: 'jobB' }])
  })
})
