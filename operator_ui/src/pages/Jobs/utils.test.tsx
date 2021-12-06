import { getTaskList } from './utils'

describe('pages/jobs/utils', () => {
  describe('getTaskList', () => {
    it('parse string to Toml Stratify list', () => {
      expect(
        getTaskList({
          value:
            'observationSource = """ ds [type=ds]; ds_parse [type=ds_parse];  """',
        }),
      ).toEqual({
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
        error: '',
      })
    })

    it('return false on bad toml format', () => {
      expect(
        getTaskList({
          value:
            'observationSource = "" ds [type=ds]; ds_parse [type=ds_parse];  """',
        }),
      ).toEqual({
        list: false,
        error: '',
      })
    })

    it('return false on circular dependency', () => {
      expect(
        getTaskList({
          value: 'observationSource = """ ds -> ds_parse -> ds  """',
        }),
      ).toEqual({
        list: false,
        error: '',
      })
    })

    it('returns an error on duplicate parents', () => {
      expect(
        getTaskList({
          value:
            'observationSource = """ ds -> ds_parse; ds1 -> ds1_parse; ds -> ds_parse;  """',
        }),
      ).toEqual({
        list: false,
        error: 'ds has duplicate ds_parse children',
      })
    })
  })
})
