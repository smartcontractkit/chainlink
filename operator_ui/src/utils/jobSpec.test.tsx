import { JobSpecFormats, getJobSpecFormat, stringifyJobSpec } from './jobSpec'

describe('utils/jobSpec', () => {
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
    it('stringify and indent job spec', async () => {
      expect(
        stringifyJobSpec({
          value: '{"foo":"bar"',
        }),
      ).toEqual('{"foo":"bar"')

      expect(
        stringifyJobSpec({
          value: '{"foo":"bar"}',
        }),
      ).toEqual(
        `{
    "foo": "bar"
}`,
      )
    })
  })
})
