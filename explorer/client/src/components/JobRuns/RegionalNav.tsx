import React from 'react'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import { ThemeStyle } from '@material-ui/core/styles/createTypography'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Typography from '@material-ui/core/Typography'
import Hidden from '@material-ui/core/Hidden'
import moment from 'moment'
import TimeAgo from '../../components/TimeAgo'

type Variant = ThemeStyle | 'srOnly' | 'inherit'

interface IJobRunProps {
  jobRun?: IJobRun
  variant: Variant
}

const JobRunId = ({ jobRun, variant }: IJobRunProps) => {
  return (
    <Typography variant={variant} color="secondary" gutterBottom>
      {jobRun ? jobRun.runId : '...'}
    </Typography>
  )
}

interface ICreatedProps {
  jobRun?: IJobRun
  showTimeAgo?: boolean
}

const Created = ({ jobRun, showTimeAgo }: ICreatedProps) => {
  return (
    <>
      {jobRun && (
        <Typography variant="subtitle2" color="textSecondary">
          Created <TimeAgo tooltip={false}>{jobRun.createdAt}</TimeAgo>
          {showTimeAgo && ` (${moment(jobRun.createdAt).format()})`}
        </Typography>
      )}
    </>
  )
}

const regionalNavStyles = ({ spacing }: Theme) =>
  createStyles({
    container: {
      padding: spacing.unit * 5
    }
  })

interface IRegionalNavProps extends WithStyles<typeof regionalNavStyles> {
  jobRun?: IJobRun
}

const RegionalNav = withStyles(regionalNavStyles)(
  ({ jobRun, classes }: IRegionalNavProps) => {
    return (
      <Paper square className={classes.container}>
        <Grid container spacing={0}>
          <Hidden xsDown>
            <Grid item xs={12}>
              <JobRunId jobRun={jobRun} variant="h3" />
            </Grid>
            <Grid item xs={12}>
              <Created jobRun={jobRun} showTimeAgo />
            </Grid>
          </Hidden>

          <Hidden smUp>
            <Grid item xs={12}>
              <JobRunId jobRun={jobRun} variant="subtitle1" />
            </Grid>
            <Grid item xs={12}>
              <Created jobRun={jobRun} />
            </Grid>
          </Hidden>
        </Grid>
      </Paper>
    )
  }
)

export default RegionalNav
