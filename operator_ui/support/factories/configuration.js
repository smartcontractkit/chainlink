export default (configOptions) => {
  return {
    data: {
      id: 'someConfigId',
      type: 'configPrinters',
      attributes: configOptions,
    },
  }
}
