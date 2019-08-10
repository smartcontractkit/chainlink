import { AppState } from 'connectors/redux/reducers'

export default ({ jobRuns }: AppState) => jobRuns.currentJobRunsCount
