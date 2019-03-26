import React from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import Header from './containers/Header'
import JobRunsIndex from './containers/JobRuns/Index'
import {
  createStyles,
  Theme,
  withStyles,
  WithStyles
} from '@material-ui/core/styles'

const styles = ({ spacing }: Theme) =>
  createStyles({
    main: {
      marginTop: 90,
      paddingLeft: spacing.unit * 5,
      paddingRight: spacing.unit * 5
    }
  })

interface IProps extends WithStyles<typeof styles> {}

const App = (props: IProps) => {
  return (
    <>
      <CssBaseline />

      <Grid container spacing={24}>
        <Grid item xs={12}>
          <Header />

          <main className={props.classes.main}>
            <JobRunsIndex />
          </main>
        </Grid>
      </Grid>
    </>
  )
}

export default hot(withStyles(styles)(App))
