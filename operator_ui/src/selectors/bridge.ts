import build from 'redux-object'
import { IState } from '../connectors/redux/reducers/index'

export default ({ bridges }: Pick<IState, 'bridges'>, id: string) => {
  return build(bridges, 'items', id)
}
