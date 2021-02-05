import {
  JobSpecFormats,
  getJobSpecFormat,
  stringifyJobSpec,
  getTaskList,
} from './utils'

describe('pages/jobs/utils', () => {
  describe('getJobSpecFormat', () => {
    it('return job spec format', async () => {
      expect(
        getJobSpecFormat({
          value: '"foo"="bar"',
        }),
      ).toEqual(JobSpecFormats.TOML)

      expect(
        getJobSpecFormat({
          value: '"foo""bar"',
        }),
      ).toEqual(false)

      expect(
        getJobSpecFormat({
          value: '{"foo":"bar"}',
        }),
      ).toEqual(JobSpecFormats.JSON)

      expect(
        getJobSpecFormat({
          value: '{"foo":"bar"',
        }),
      ).toEqual(false)

      expect(
        getJobSpecFormat({
          value: '',
        }),
      ).toEqual(false)
    })
  })

  describe('stringifyJobSpec', () => {
    it('stringify and indent JSON job spec', async () => {
      expect(
        stringifyJobSpec({
          value: { foo: 'bar' },
          format: JobSpecFormats.JSON,
        }),
      ).toEqual(
        `{
    "foo": "bar"
}`,
      )
    })

    it('stringify TOML spec', async () => {
      expect(
        stringifyJobSpec({
          value: { foo: 'bar' },
          format: JobSpecFormats.TOML,
        }),
      ).toEqual(`foo = "bar"
`)
    })
  })

  describe('getTaskList', () => {
    it('parse string to Json TaskSpec list', async () => {
      expect(
        getTaskList({
          value: '{"tasks":[{ "type": "HTTPGet"}, { "type": "JSONParse"}]}',
        }),
      ).toEqual({
        format: 'json',
        list: [{ type: 'HTTPGet' }, { type: 'JSONParse' }],
      })
    })

    it('return false on bad json format', async () => {
      expect(
        getTaskList({
          value: '{"tasks":[{ "type": HTTPGet}, { "type": JSONParse}]}',
        }),
      ).toEqual({
        format: false,
        list: false,
      })
    })

    it('parse string to Toml Stratify list', async () => {
      expect(
        getTaskList({
          value:
            'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse];  """',
        }),
      ).toEqual({
        format: 'toml',
        list: [
          {
            attributes: {
              type: 'ds',
            },
            id: 'ds',
            parentIds: [],
          },
          {
            attributes: {
              type: 'ds_parse',
            },
            id: 'ds_parse',
            parentIds: [],
          },
        ],
      })
    })

    it('return false on bad toml format', async () => {
      expect(
        getTaskList({
          value:
            'observationSource = "" ds [type=ds]; ds_parse [type=ds_parse];  """',
        }),
      ).toEqual({
        format: false,
        list: false,
      })
    })
  })
})
