import { PaddedCard } from '@chainlink/styleguide'
import Grid from '@material-ui/core/Grid'
import { withStyles } from '@material-ui/core/styles'
import { fetchJobRun } from 'actions'
import Content from 'components/Content'
import StatusCard from 'components/JobRuns/StatusCard'
import PrettyJson from 'components/PrettyJson'
import React from 'react'
import { connect } from 'react-redux'
import jobRunSelector from 'selectors/jobRun'
import { useEffect, useHooks } from 'use-react-hooks'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import RegionalNav from './RegionalNav'

const styles = theme => ({})

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
        <StatusCard title={jobRun.status} jobRun={jobRun} />
      </Grid>
    </Grid>
  )
}

const Show = useHooks(props => {
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
})

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

export default withStyles(styles)(ConnectedShow)
