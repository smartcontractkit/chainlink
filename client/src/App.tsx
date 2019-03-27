import React from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'
import { Router } from '@reach/router'
import Header from './containers/Header'
import JobRunsIndex from './containers/JobRuns/Index'
import JobRunsShow from './containers/JobRuns/Show'

const styles = ({ spacing }: Theme) =>
  createStyles({
    main: {
      marginTop: 90
    }
  })

interface IProps extends WithStyles<typeof styles> {
  children: any
  path: string
}

const Main = withStyles(styles)(({ children, classes }: IProps) => (
  <>
    <Header />
    <main className={classes.main}>{children}</main>
  </>
))

const App = () => {
  return (
    <>
      <CssBaseline />

      <Grid container spacing={24}>
        <Grid item xs={12}>
          <Router>
            <Main path="/">
              <JobRunsIndex path="/" />
              <JobRunsShow path="/job-runs/:jobRunId" />
            </Main>
          </Router>
        </Grid>
      </Grid>
    </>
  )
}

export default hot(App)
