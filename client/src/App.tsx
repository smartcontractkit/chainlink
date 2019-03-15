import React from 'react'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Header from './components/Header'
import Home from './containers/Home'
import { createStyles, Theme, withStyles, WithStyles } from '@material-ui/core/styles'

const styles = ({spacing}: Theme) => createStyles({
  main: {
    marginTop: 90,
    paddingLeft: spacing.unit * 5,
    paddingRight: spacing.unit * 5
  }
})

interface IProps extends WithStyles<typeof styles> {
}

const App = (props: IProps) => {
  return (
    <Grid container spacing={24}>
      <Grid item xs={12}>
        <Header />

        <Paper className={props.classes.main} elevation={0}>
          <Home />
        </Paper>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(App)
