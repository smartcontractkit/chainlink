import React from 'react'

import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles,
} from '@material-ui/core/styles'
import Typography from '@material-ui/core/Typography'

import { JobData } from './sharedTypes'
import { TimeAgo } from 'src/components/TimeAgo'

const styles = (theme: Theme) =>
  createStyles({
    paper: {
      margin: `${theme.spacing.unit * 2.5}px 0`,
      padding: theme.spacing.unit * 3,
    },
  })

interface Props extends WithStyles<typeof styles> {
  job: JobData['job']
}

export const DetailsBox = withStyles(styles)(({ classes, job }: Props) => {
  if (!job) {
    return null
  }

  return (
    <Paper className={classes.paper}>
      <Grid container>
        <Grid item xs={12} sm={6} md={1}>
          <Typography variant="subtitle2" gutterBottom>
            ID
          </Typography>
          <Typography variant="body1">{job.id}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={2}>
          <Typography variant="subtitle2" gutterBottom>
            Type
          </Typography>
          <Typography variant="body1">{job.specType}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={5}>
          <Typography variant="subtitle2" gutterBottom>
            External Job ID
          </Typography>
          <Typography variant="body1">{job.externalJobID}</Typography>
        </Grid>
        <Grid item xs={12} sm={6} md={4}>
          <Typography variant="subtitle2" gutterBottom>
            Created At
          </Typography>
          <Typography variant="body1">
            <TimeAgo tooltip>{job.createdAt}</TimeAgo>
          </Typography>
        </Grid>
      </Grid>
    </Paper>
  )
})
