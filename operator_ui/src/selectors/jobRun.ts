import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ jobRuns }: Pick<AppState, 'jobRuns'>, id: string) => {
  return build(jobRuns, 'items', id)
}
