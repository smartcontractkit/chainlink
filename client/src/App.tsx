import React, { useState } from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import { Router } from '@reach/router'
import Header from './containers/Header'
import JobRunsIndex from './containers/JobRuns/Index'
import JobRunsShow from './containers/JobRuns/Show'

interface IProps {
  children: any
  path: string
}

const DEFAULT_HEIGHT = 82

const Main = ({ children }: IProps) => {
  const [height, setHeight] = useState<number>(DEFAULT_HEIGHT)
  const onHeaderResize = (width: number, height: number) => {
    setHeight(height)
  }

  return (
    <>
      <Header onResize={onHeaderResize} />
      <main style={{ paddingTop: height }}>{children}</main>
    </>
  )
}

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
