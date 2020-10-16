import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRun } from 'actionCreators'
import jobRunSelector from 'selectors/jobRun'
import Content from 'components/Content'
import StatusCard from 'components/JobRuns/StatusCard'
import TaskExpansionPanel from 'components/JobRuns/TaskExpansionPanel'
import RegionalNav from './RegionalNav'

const styles = (theme) => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5,
  },
})

const renderDetails = ({ fetching, jobRun }) => {
  if (fetching || !jobRun) {
    return <div>Fetching job run...</div>
  }

  return (
    <Grid container spacing={40}>
      <Grid item xs={8}>
        <Card>
          <TaskExpansionPanel jobRun={jobRun} />
        </Card>
      </Grid>
      <Grid item xs={4}>
        <StatusCard title={jobRun.status} jobRun={jobRun} />
      </Grid>
    </Grid>
  )
}

export const Show = (props) => {
  const { fetchJobRun, jobRunId } = props
  useEffect(() => {
    document.title = 'Show Job Run'
    fetchJobRun(jobRunId)
  }, [fetchJobRun, jobRunId])

  return (
    <div>
      <RegionalNav {...props} />
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
    fetching,
  }
}

export const ConnectedShow = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobRun }),
)(Show)

export default withStyles(styles)(ConnectedShow)
