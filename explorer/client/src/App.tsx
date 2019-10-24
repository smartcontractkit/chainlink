import React from 'react'
import { hot } from 'react-hot-loader/root'
import CssBaseline from '@material-ui/core/CssBaseline'
import { Router, Redirect } from '@reach/router'
import PublicLayout from './layouts/Public'
import AdminMinimalLayout from './layouts/AdminMinimal'
import AdminLayout from './layouts/Admin'
import NotFound from './components/NotFound'
import SearchCard from './components/Cards/Search'
import JobRunsIndex from './containers/JobRuns/Index'
import JobRunsShow from './containers/JobRuns/Show'
import AdminSignIn from './containers/Admin/SignIn'
import AdminSignOut from './containers/Admin/SignOut'
import AdminOperatorIndex from './containers/Admin/Operator/Index'

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

        <AdminMinimalLayout path="/admin/signin">
          <AdminSignIn default />
        </AdminMinimalLayout>
        <AdminMinimalLayout path="/admin/signout">
          <AdminSignOut default />
        </AdminMinimalLayout>

        <AdminLayout path="/admin">
          <AdminOperatorIndex path="/operators" />
          <Redirect path="/" from="/" to="/admin/operators" noThrow />
          <NotFound default />
        </AdminLayout>
      </Router>
    </>
  )
}

export default hot(App)
