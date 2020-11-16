import { JobSpecFormats, getJobSpecFormat, stringifyJobSpec } from './utils'

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
          value: '{"foo":"bar"',
          format: JobSpecFormats.JSON,
        }),
      ).toEqual('{"foo":"bar"')

      expect(
        stringifyJobSpec({
          value: '{"foo":"bar"}',
          format: JobSpecFormats.JSON,
        }),
      ).toEqual(
        `{
    "foo": "bar"
}`,
      )
    })

    it('returns TOML format value', async () => {
      expect(
        stringifyJobSpec({
          value: 'foo="bar"',
          format: JobSpecFormats.TOML,
        }),
      ).toEqual('foo="bar"')
    })
  })
})
