import jsonApiJobSpecsFactory from 'factories/jsonApiJobSpecs'

export default (job) => {
  const unshaped = jsonApiJobSpecsFactory([job])
  return { data: unshaped.data[0] }
}
