import transformJob from 'actions/transforms/jobs'

describe('actions/transforms/jobs', () => {
  it('returns an action for the type with serialized job items', () => {
    const json = {data: [{}, {}]}
    const action = transformJob('MY_ACTION', json, j => j)

    expect(action).toEqual({
      type: 'MY_ACTION',
      items: [{}, {}]
    })
  })

  it('includes the count when provided in meta', () => {
    const json = {
      data: [{}, {}],
      meta: {count: 10}
    }
    const action = transformJob('MY_ACTION', json, j => j)

    expect(action).toEqual({
      type: 'MY_ACTION',
      items: [{}, {}],
      count: 10
    })
  })
})
