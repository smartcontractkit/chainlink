import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRun } from 'actions'
import jobRunSelector from 'selectors/jobRun'
import Content from 'components/Content'
import RegionalNav from 'components/JobRuns/RegionalNav'
import StatusCard from 'components/JobRuns/StatusCard'
import TaskExpansionPanel from 'components/JobRuns/TaskExpansionPanel'

const renderDetails = ({ fetching, jobRun }) => {
  if (fetching || !jobRun) {
    return <div>Fetching job run...</div>
  }

  return (
    <Grid container spacing={5}>
      <Grid item xs={8}>
        <Card>
          <TaskExpansionPanel>{jobRun}</TaskExpansionPanel>
        </Card>
      </Grid>
      <Grid item xs={4}>
        <StatusCard>{jobRun.status}</StatusCard>
      </Grid>
    </Grid>
  )
}

export const Show = props => {
  useEffect(() => {
    document.title = 'Show Job Run'
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
