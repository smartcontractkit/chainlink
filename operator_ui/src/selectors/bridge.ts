import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ bridges }: Pick<AppState, 'bridges'>, id: string) => {
  return build(bridges, 'items', id)
}
