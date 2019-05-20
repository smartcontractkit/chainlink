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
          <Grid item xs={12}>
            <Hidden xsDown>
              <JobRunId jobRun={jobRun} variant="h3" />
            </Hidden>
            <Hidden smUp>
              <JobRunId jobRun={jobRun} variant="h5" />
            </Hidden>
          </Grid>
          <Grid item xs={12}>
            <Typography variant="subtitle2" color="textSecondary">
              {jobRun && (
                <>
                  Created <TimeAgo tooltip={false}>{jobRun.createdAt}</TimeAgo>{' '}
                  ({moment(jobRun.createdAt).format()})
                </>
              )}
            </Typography>
          </Grid>
        </Grid>
      </Paper>
    )
  }
)

export default RegionalNav
