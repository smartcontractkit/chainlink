import { AppState } from 'reducers'
import build from 'redux-object'

export default (
  { transactions }: Pick<AppState, 'transactions'>,
  id: string,
) => {
  return build(transactions, 'items', id, { eager: true })
}
