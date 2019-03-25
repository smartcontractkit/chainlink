import React, { useEffect } from 'react'
import { denormalize } from 'normalizr'
import { connect } from 'react-redux'
import { bindActionCreators, Dispatch } from 'redux'
import { getJobRun } from '../../actions/jobRuns'
import { IState } from '../../reducers'
import { JobRun } from '../../entities'

type IProps = {
  jobRunId?: string
  jobRun?: IJobRun
  getJobRun: Function
  path: string
}

interface IOwnProps {
  jobRunId?: string
}

const Show = ({ jobRunId, jobRun, getJobRun }: IProps) => {
  useEffect(() => {
    getJobRun(jobRunId)
  }, [])

  if (jobRun) {
    return (
      <div>
        <div>Run ID: {jobRun.id}</div>
        <div>Job ID: {jobRun.jobId}</div>
      </div>
    )
  }

  return <div>Loading run {jobRunId}...</div>
}

const jobRunSelector = (
  { jobRuns }: IState,
  jobRunId?: string
): IJobRun | undefined => {
  if (jobRuns.items) {
    return denormalize(jobRunId, JobRun, { jobRuns: jobRuns.items })
  }
}

const mapStateToProps = (state: IState, { jobRunId }: IOwnProps) => {
  return {
    jobRun: jobRunSelector(state, jobRunId)
  }
}

const mapDispatchToProps = (dispatch: Dispatch<any>) =>
  bindActionCreators({ getJobRun }, dispatch)

const ConnectedShow = connect(
  mapStateToProps,
  mapDispatchToProps
)(Show)

export default ConnectedShow
