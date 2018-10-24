import jobSelector from './job'

export default (state, jobSpecId) => {
  const spec = jobSelector(state, jobSpecId)
  return spec ? spec.runsCount : 0
}
