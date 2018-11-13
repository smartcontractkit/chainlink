import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'

export default job => {
  let unshaped = jsonApiJobSpecsFactory([job])
  return { data: unshaped.data[0] }
}
