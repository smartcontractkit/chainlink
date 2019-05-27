import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import PaddedCard from 'components/PaddedCard'
import PrettyJson from 'components/PrettyJson'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRun } from 'actions'
import jobRunSelector from 'selectors/jobRun'
import Content from 'components/Content'
import RegionalNav from 'components/JobRuns/RegionalNav'
import StatusCard from 'components/JobRuns/StatusCard'

const renderDetails = ({ fetching, jobRun }) => {
  if (fetching || !jobRun) {
    return <div>Fetching job run...</div>
  }

  return (
    <Grid container spacing={40}>
      <Grid item xs={8}>
        <PaddedCard>
          <PrettyJson object={jobRun} />
        </PaddedCard>
      </Grid>
      <Grid item xs={4}>
        <StatusCard>{jobRun.status}</StatusCard>
      </Grid>
    </Grid>
  )
}

const Show = props => {
  useEffect(() => {
    props.fetchJobRun(props.jobRunId)
  }, [])

  return (
    <div>
      <RegionalNav
        jobSpecId={props.jobSpecId}
        jobRunId={props.jobRunId}
        jobRun={props.jobRun}
      />

      <Content>{renderDetails(props)}</Content>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  const { jobSpecId, jobRunId } = ownProps.match.params
  const jobRun = jobRunSelector(state, jobRunId)
  const fetching = state.jobRuns.fetching

  return {
    jobSpecId,
    jobRunId,
    jobRun,
    fetching
  }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRun })
)(Show)

export default ConnectedShow
