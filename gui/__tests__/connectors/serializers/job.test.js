import jobSerializer from 'connectors/redux/serializers/job'

describe('connectors/redux/serializers/job', () => {
  it('transforms a JSON api response to an object that will be stored', () => {
    const json = {
      id: 'idA',
      attributes: {
        createdAt: '2018-10-22T21:54:04.84278Z'
      }
    }

    const job = jobSerializer(json)

    expect(job.id).toEqual('idA')
    expect(job.createdAt).toEqual(1540245244842)
  })
})
