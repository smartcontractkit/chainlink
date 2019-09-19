import React from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import { Router } from '@reach/router'
import PublicLayout from './layouts/Public'
import AdminLayout from './layouts/Admin'
import SearchCard from './components/Cards/Search'
import JobRunsIndex from './containers/JobRuns/Index'
import JobRunsShow from './containers/JobRuns/Show'
import AdminSignIn from './containers/Admin/SignIn'

const App = () => {
  return (
    <>
      <CssBaseline />

      <Router style={{ display: 'flex', height: '100%', overflowX: 'hidden' }}>
        <SearchCard path="/" />

        <PublicLayout path="/job-runs">
          <JobRunsIndex path="/" />
          <JobRunsShow path="/:jobRunId" />
        </PublicLayout>

        <AdminLayout path="/admin">
          <AdminSignIn path="/signin" />
        </AdminLayout>
      </Router>
    </>
  )
}

export default hot(App)
