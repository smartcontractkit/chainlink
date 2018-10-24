import jobSelector from 'selectors/job'

describe('selectors - job', () => {
  it('returns the job item for the given id and undefined otherwise', () => {
    const state = {
      jobs: {
        items: {
          jobA: {id: 'jobA'}
        }
      }
    }

    expect(jobSelector(state, 'jobA')).toEqual({id: 'jobA'})
    expect(jobSelector(state, 'joba')).toBeUndefined()
  })
})
