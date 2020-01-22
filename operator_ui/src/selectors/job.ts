import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ jobs }: Pick<AppState, 'jobs'>, id: string) => {
  return build(jobs, 'items', id, { eager: true })
}
