import build from 'redux-object'

export default ({ jobRuns }, id) => build(jobRuns, 'items', id)
