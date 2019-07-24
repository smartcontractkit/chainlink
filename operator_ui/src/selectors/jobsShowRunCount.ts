import { IState } from '../connectors/redux/reducers/index'

export default ({ jobRuns }: IState) => jobRuns.currentJobRunsCount
