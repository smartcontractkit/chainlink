import jobSpecDefinition from 'utils/jobSpecDefinition'

const input = `{
  "initiators": [
    {
      "type": "web",
      "params": {
      }
    }
  ],
  "tasks": [
    {
      "ID": 1,
      "CreatedAt": "2019-05-27T11:31:27.81655-04:00",
      "UpdatedAt": "2019-05-27T11:31:27.81655-04:00",
      "DeletedAt": null,
      "type": "httpget",
      "confirmations": 0,
      "params": {
        "ID": "KeepMeBecauseIAmUserDefined",
        "get": "https://localhost:9000/file.tar.gz"
      }
    }
  ],
  "startAt": null,
  "endAt": null
}`

describe('utils/jobSpecDefinition', () => {
  it('scrubs unwanted keys', () => {
    const output = jobSpecDefinition(JSON.parse(input))
    expect(output.tasks).toHaveLength(1)
    expect(output.tasks[0].ID).toBeUndefined()
    expect(output.tasks[0].CreatedAt).toBeUndefined()
    expect(output.tasks[0].DeletedAt).toBeUndefined()
    expect(output.tasks[0].UpdatedAt).toBeUndefined()
    expect(output.tasks[0].params.ID).toEqual('KeepMeBecauseIAmUserDefined')
    expect(output.tasks[0].params.get).toEqual(
      'https://localhost:9000/file.tar.gz',
    )
  })
})
