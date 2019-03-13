import React from 'react'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import Header from './components/Header'
import Home from './containers/Home'
import { withStyles } from '@material-ui/core/styles'

const styles = () => {
  return {
    main: {
      marginTop: 90
    }
  }
}

const App = (props: any) => {
  return (
    <Grid container>
      <Grid item xs={12}>
        <Header />

        <main className={props.classes.main}>
          <Home />
        </main>
      </Grid>
    </Grid>
  )
}

export default withStyles(styles)(App)
