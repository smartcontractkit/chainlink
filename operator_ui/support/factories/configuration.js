import { decamelizeKeys } from 'humps'

export default (configOptions) => {
  return decamelizeKeys({
    data: {
      id: 'someConfigId',
      type: 'configPrinters',
      attributes: configOptions,
    },
  })
}
