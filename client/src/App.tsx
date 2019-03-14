import React from 'react'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import Paper from '@material-ui/core/Paper'
import Header from './components/Header'
import Home from './containers/Home'
import { withStyles } from '@material-ui/core/styles'

const styles = (theme: any) => {
  return {
    main: {
      marginTop: 90,
      paddingLeft: theme.spacing.unit * 5,
      paddingRight: theme.spacing.unit * 5
    }
  }
}

const App = (props: any) => {
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
