import React, { useEffect } from 'react'
import { connect } from 'react-redux'
// import Grid from '@material-ui/core/Grid'
import Content from 'components/Content'
import { KeyValueList } from '@chainlink/styleguide'
import { JobSpec } from 'operator_ui'
import { AppState } from 'src/reducers'
import jobSelector from 'selectors/job'
import { fetchJob } from 'actions'
import RegionalNav from './RegionalNav'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

interface Props {
  jobSpecId: string
  job?: JobSpec
  fetchJob: (id: string) => Promise<any>
}

export const JobSpecErrors: React.FC<Props> = ({
  jobSpecId,
  job,
  fetchJob,
}) => {
  useEffect(() => {
    document.title = 'Job Errors'
    fetchJob(jobSpecId)
  }, [fetchJob, jobSpecId])
  return (
    <>
      {/* TODO: Regional nav should handle job = undefined */}
      {job && <RegionalNav jobSpecId={jobSpecId} job={job} />}
      <Content>
        <KeyValueList title="Errors" entries={[]} showHead titleize />
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
  matchRouteAndMapDispatchToProps({ fetchJob }),
)(JobSpecErrors)

export default ConnectedJobSpecErrors
