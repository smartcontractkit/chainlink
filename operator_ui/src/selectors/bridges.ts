import build from 'redux-object'
import { IState } from '../connectors/redux/reducers/index'

export default ({ bridges }: IState) => {
  return (
    bridges.currentPage &&
    bridges.currentPage.map(id => build(bridges, 'items', id)).filter(b => b)
  )
}
