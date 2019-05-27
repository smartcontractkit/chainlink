import React, { useEffect } from 'react'
import Grid from '@material-ui/core/Grid'
import Card from '@material-ui/core/Card'
import CardContent from '@material-ui/core/CardContent'
import Typography from '@material-ui/core/Typography'
import Divider from '@material-ui/core/Divider'
import Content from 'components/Content'
import RegionalNav from 'components/Jobs/RegionalNav'
import PrettyJson from 'components/PrettyJson'
import { connect } from 'react-redux'
import { fetchJob, createJobRun } from 'actions'
import jobSelector from 'selectors/job'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import jobSpecDefinition from 'utils/jobSpecDefinition'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  definitionTitle: {
    marginTop: theme.spacing(2),
    marginBottom: theme.spacing(2)
  },
  divider: {
    marginTop: theme.spacing(1),
    marginBottom: theme.spacing(3)
  }
}))

const renderDetails = ({ job }) => {
  const definition = job && jobSpecDefinition(job)
  const classes = useStyles()
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
  useEffect(() => {
    document.title = 'Job Definition'
    props.fetchJob(props.jobSpecId)
  }, [])
  const { jobSpecId, job } = props

  return (
    <div>
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
    job
  }
}

export const ConnectedDefinition = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJob, createJobRun })
)(Definition)

export default ConnectedDefinition
