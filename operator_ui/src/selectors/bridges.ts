import { AppState } from 'reducers'
import build from 'redux-object'

export default ({ bridges }: Pick<AppState, 'bridges'>) => {
  return (
    bridges.currentPage &&
    bridges.currentPage
      .map((id) => build(bridges, 'items', id))
      .filter((b) => b)
  )
}
