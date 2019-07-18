import { constantCase } from 'change-case'
import { IState } from '../connectors/redux/reducers/index'

export default ({ configuration }: IState) => {
  const { data } = configuration

  return Object.keys(data)
    .sort()
    .map(key => [constantCase(key), data[key]])
}
