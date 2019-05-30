import React, { useState } from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import Grid from '@material-ui/core/Grid'
import { Router } from '@reach/router'
import Header from './containers/Header'
import SearchCard from './components/Cards/Search'
import TermsOfUse from './components/TermsOfUse'
import JobRunsIndex from './containers/JobRuns/Index'
import JobRunsShow from './containers/JobRuns/Show'

interface IMainProps {
  children: any
  path: string
}

const DEFAULT_HEIGHT = 82

const Main = ({ children }: IMainProps) => {
  const [height, setHeight] = useState<number>(DEFAULT_HEIGHT)
  const onHeaderResize = (width: number, height: number) => {
    setHeight(height)
  }

  return (
    <Grid container spacing={24}>
      <Grid item xs={12}>
        <Header onResize={onHeaderResize} />
        <main style={{ paddingTop: height }}>{children}</main>
        <TermsOfUse />
      </Grid>
    </Grid>
  )
}

const App = () => {
  return (
    <>
      <CssBaseline />

      <Router style={{ display: 'flex', height: '100%' }}>
        <SearchCard path="/" />

        <Main path="/job-runs">
          <JobRunsIndex path="/" />
          <JobRunsShow path="/:jobRunId" />
        </Main>
      </Router>
    </>
  )
}

export default hot(App)
