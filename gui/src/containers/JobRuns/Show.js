import React, { Component } from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import Breadcrumb from 'components/Breadcrumb'
import BreadcrumbItem from 'components/BreadcrumbItem'
import PaddedCard from 'components/PaddedCard'
import PrettyJson from 'components/PrettyJson'
import Title from 'components/Title'
import TimeAgo from 'components/TimeAgo'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { fetchJobRun } from 'actions'
import jobRunSelector from 'selectors/jobRun'

const styles = theme => ({
  breadcrumb: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  }
})

const renderDetails = ({fetching, jobRun}) => {
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
        <PaddedCard>
          <Grid container spacing={16}>
            <Grid item xs={12}>
              <Typography variant='subheading' color='textSecondary'>ID</Typography>
              <Typography variant='body1' color='inherit'>
                {jobRun.id}
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography variant='subheading' color='textSecondary'>Status</Typography>
              <Typography variant='body1' color='inherit'>
                {jobRun.status}
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography variant='subheading' color='textSecondary'>Created</Typography>
              <Typography variant='body1' color='inherit'>
                <TimeAgo>{jobRun.createdAt}</TimeAgo>
              </Typography>
            </Grid>
            <Grid item xs={12}>
              <Typography variant='subheading' color='textSecondary'>Result</Typography>
              <Typography variant='body1' color='inherit'>
                {jobRun.result && JSON.stringify(jobRun.result.data)}
              </Typography>
            </Grid>
          </Grid>
        </PaddedCard>
      </Grid>
    </Grid>
  )
}

export class Show extends Component {
  componentDidMount () {
    this.props.fetchJobRun(this.props.jobRunId)
  }

  render () {
    const {props} = this

    return (
      <div>
        <Breadcrumb className={props.classes.breadcrumb}>
          <BreadcrumbItem href='/'>Dashboard</BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem href={`/jobs/${props.jobSpecId}`}>
            Job ID: {props.jobSpecId}
          </BreadcrumbItem>
          <BreadcrumbItem>></BreadcrumbItem>
          <BreadcrumbItem>Job Run ID: {props.jobRunId}</BreadcrumbItem>
        </Breadcrumb>
        <Title>Job Run Detail</Title>

        {renderDetails(props)}
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => {
  const {jobSpecId, jobRunId} = ownProps.match.params
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
  matchRouteAndMapDispatchToProps({fetchJobRun})
)(Show)

export default withStyles(styles)(ConnectedShow)
