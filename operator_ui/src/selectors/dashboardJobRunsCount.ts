import { IState } from '../connectors/redux/reducers/index'

export default ({ dashboardIndex }: IState): number | undefined =>
  dashboardIndex.jobRunsCount
