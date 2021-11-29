import React from 'react'

import CircularProgress from '@material-ui/core/CircularProgress'
import Grid from '@material-ui/core/Grid'
import { Theme, withStyles, WithStyles } from '@material-ui/core/styles'

const styles = (theme: Theme) => ({
  root: {
    margin: theme.spacing.unit * 2,
  },
  gridItem: {
    display: 'flex',
    justifyContent: 'center',
  },
})

interface Props extends WithStyles<typeof styles> {}

export const Loading = withStyles(styles)(({ classes }: Props) => (
  <Grid container className={classes.root}>
    <Grid item xs={12} className={classes.gridItem}>
      <CircularProgress data-testid="loading-spinner" />
    </Grid>
  </Grid>
))
