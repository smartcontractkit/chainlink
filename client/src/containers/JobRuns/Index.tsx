import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch } from 'redux'
import { denormalize } from 'normalizr'
import List from '../../components/JobRuns/List'
import { getJobRuns } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { JobRun } from '../../entities'

type IProps = {
  query?: string
  jobRuns?: IJobRun[]
  getJobRuns: Function
  path: string
}

const Index = ({ query, jobRuns, getJobRuns }: IProps) => {
  useEffect(() => {
    getJobRuns(query)
  }, [])

  return <List jobRuns={jobRuns} />
}

const jobRunsSelector = ({
  jobRunsIndex,
  jobRuns
}: IState): IJobRun[] | undefined => {
  if (jobRunsIndex.items) {
    return denormalize(jobRunsIndex.items, [JobRun], { jobRuns: jobRuns.items })
  }
}

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
