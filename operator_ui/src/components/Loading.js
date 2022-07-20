import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme) => ({
  wrapper: {
    marginTop: theme.spacing.unit * 5,
  },
  text: {
    textAlign: 'center',
  },
})

const Loading = ({ classes }) => (
  <Grid container alignItems="center">
    <Grid item xs={12} className={classes.wrapper}>
      <Typography variant="h4" className={classes.text}>
        Loading...
      </Typography>
    </Grid>
  </Grid>
)

export default withStyles(styles)(Loading)
