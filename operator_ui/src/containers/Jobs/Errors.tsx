import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import Content from 'components/Content'
import { JobSpec } from 'operator_ui'
import { AppState } from 'src/reducers'
import jobSelector from 'selectors/job'
import { fetchJob } from 'actionCreators'
import RegionalNav from './RegionalNav'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import List from 'components/JobErrors/List'
import { deleteJobSpecError } from 'actionCreators'

interface Props {
  jobSpecId: string
  job?: JobSpec
  fetchJob: (id: string) => Promise<any>
  deleteJobSpecError: (jobSpecErrorId: number) => Promise<any>
}

export const JobSpecErrors: React.FC<Props> = ({
  jobSpecId,
  job,
  fetchJob,
  deleteJobSpecError,
}) => {
  useEffect(() => {
    document.title = 'Job Errors'
    fetchJob(jobSpecId)
  }, [fetchJob, jobSpecId])

  const handleDismiss = (jobSpecErrorId: number) => {
    deleteJobSpecError(jobSpecErrorId)
    // fetchJob(jobSpecId)
  }

  return (
    <>
      <RegionalNav jobSpecId={jobSpecId} job={job} />
      <Content>
        <List errors={job?.errors} dismiss={handleDismiss} />
      </Content>
    </>
  )
}

interface Match {
  params: {
    jobSpecId: string
  }
}

const mapStateToProps = (state: AppState, ownProps: { match: Match }) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const job = jobSelector(state, jobSpecId)

  return {
    jobSpecId,
    job,
  }
}

export const ConnectedJobSpecErrors = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob, deleteJobSpecError }),
)(JobSpecErrors)

export default ConnectedJobSpecErrors
