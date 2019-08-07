import { constantCase } from 'change-case'
import { IState } from 'connectors/redux/reducers/index'

export default ({ configuration }: Pick<IState, 'configuration'>) => {
  const { data } = configuration

  return Object.keys(data)
    .sort()
    .map(key => [constantCase(key), data[key]])
}
