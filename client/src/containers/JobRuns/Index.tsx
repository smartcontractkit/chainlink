import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch } from 'redux'
import JobRunsList from '../../components/JobRunsList'
import { getJobRuns } from '../../actions/jobRuns'
import { IState } from '../../reducers'

type IProps = {
  query?: string
  jobRuns?: IJobRun[]
  getJobRuns: Function
}

const Index = ({ query, jobRuns, getJobRuns }: IProps) => {
  useEffect(() => {
    getJobRuns(query)
  }, [])

  return <JobRunsList jobRuns={jobRuns} />
}

const jobRunsSelector = (state: IState): IJobRun[] | undefined =>
  state.jobRuns.items

const mapStateToProps = (state: IState) => ({
  query: state.search.query,
  jobRuns: jobRunsSelector(state)
})

const mapDispatchToProps = (dispatch: Dispatch<any>) =>
  bindActionCreators({ getJobRuns }, dispatch)

const ConnectedIndex = connect(
  mapStateToProps,
  mapDispatchToProps
)(Index)

export default ConnectedIndex
