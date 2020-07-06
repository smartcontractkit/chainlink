import { AppState } from 'reducers'
import build from 'redux-object'
import { JobSpec } from 'operator_ui'

export default (
  { jobs }: Pick<AppState, 'jobs'>,
  id: string,
): JobSpec | undefined => {
  return build(jobs, 'items', id, { eager: true })
}
