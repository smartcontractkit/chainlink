import jobsSelector from 'selectors/jobs'

describe('selectors - jobs', () => {
  it('returns the jobs in the current page', () => {
    const state = {
      jobs: {
        currentPage: ['jobA', 'jobB'],
        items: {
          jobA: { id: 'jobA' },
          jobB: { id: 'jobB' },
        },
      },
    }
    const jobs = jobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }, { id: 'jobB' }])
  })

  it('excludes job items that are not present', () => {
    const state = {
      jobs: {
        currentPage: ['jobA', 'jobB'],
        items: {
          jobA: { id: 'jobA' },
        },
      },
    }
    const jobs = jobsSelector(state)

    expect(jobs).toEqual([{ id: 'jobA' }])
  })
})
