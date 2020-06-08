import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { withStyles } from '@material-ui/core/styles'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import Divider from '@material-ui/core/Divider'
import Content from 'components/Content'
import PrettyJson from 'components/PrettyJson'
import RegionalNav from './RegionalNav'
import { fetchJob, createJobRun } from 'actions'
import jobSelector from 'selectors/job'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import jobSpecDefinition from 'utils/jobSpecDefinition'

const styles = theme => ({
  definitionTitle: {
    marginTop: theme.spacing.unit * 2,
    marginBottom: theme.spacing.unit * 2,
  },
  divider: {
    marginTop: theme.spacing.unit,
    marginBottom: theme.spacing.unit * 3,
  },
})

const renderDetails = ({ job, classes }) => {
  const definition = job && jobSpecDefinition(job)

  if (definition) {
    return (
      <Grid container spacing={0}>
        <Grid item xs={12}>
          <Typography variant="h5" className={classes.definitionTitle}>
            Definition
          </Typography>
        </Grid>
        <Grid item xs={12}>
          <Divider light className={classes.divider} />
        </Grid>
        <Grid item xs={12}>
          <PrettyJson object={definition} />
        </Grid>
      </Grid>
    )
  }

  return <React.Fragment>Fetching ...</React.Fragment>
}

const Definition = props => {
  const { fetchJob, job, jobSpecId } = props

  useEffect(() => {
    document.title = 'Job Definition'
    fetchJob(jobSpecId)
  }, [fetchJob, jobSpecId])

  return (
    <div>
      {/* TODO - RYAN - get rid of jobSpecId as argument */}
      <RegionalNav jobSpecId={jobSpecId} job={job} />
      <Content>
        <Card>
          <CardContent>{renderDetails(props)}</CardContent>
        </Card>
      </Content>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  const jobSpecId = ownProps.match.params.jobSpecId
  const job = jobSelector(state, jobSpecId)

  return {
    jobSpecId,
    job,
  }
}

export const ConnectedDefinition = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob, createJobRun }),
)(Definition)

export default withStyles(styles)(ConnectedDefinition)
