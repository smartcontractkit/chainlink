import React from 'react'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import { makeStyles } from '@material-ui/styles'

const useStyles = makeStyles(theme => ({
  wrapper: {
    marginTop: theme.spacing(5)
  },
  text: {
    textAlign: 'center'
  }
}))

const Loading = () => {
  const classes = useStyles()
  return (
    <Grid container alignItems="center">
      <Grid item xs={12} className={classes.wrapper}>
        <Typography variant="h4" className={classes.text}>
          Loading...
        </Typography>
      </Grid>
    </Grid>
  )
}

export default Loading
