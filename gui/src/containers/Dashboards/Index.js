import React, { Component } from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Title from 'components/Title'
import Button from '@material-ui/core/Button'
import { fetchJobs, fetchAccountBalance } from 'actions'
import jobsSelector from 'selectors/jobs'
import ReactStaticLinkComponent from 'components/ReactStaticLinkComponent'
import JobList from 'components/JobList'
import TokenBalance from 'components/TokenBalance'
import MetaInfo from 'components/MetaInfo'
import Footer from 'components/Footer'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'

const renderJobsList = props => {
  const { jobs, jobCount, pageSize, jobsError, fetchJobs, history, match } = props
  return (
    <JobList
      jobs={jobs}
      jobCount={jobCount}
      pageSize={pageSize}
      error={jobsError}
      fetchJobs={fetchJobs}
      history={history}
      match={match}
    />
  )
}

const renderSidebar = ({ ethBalance, linkBalance, jobCount, accountBalanceError }) => (
  <Grid container spacing={24}>
    <Grid item xs={12}>
      <TokenBalance title='Link Balance' value={linkBalance} error={accountBalanceError} />
    </Grid>
    <Grid item xs={12}>
      <TokenBalance title='Ether Balance' value={ethBalance} error={accountBalanceError} />
    </Grid>
    <Grid item xs={12}>
      <MetaInfo title='Jobs' value={jobCount} />
    </Grid>
  </Grid>
)

export class Index extends Component {
  componentDidMount () {
    this.props.fetchAccountBalance()
  }

  render () {
    return (
      <div>
        <Grid container alignItems='center' >
          <Grid item xs={9}>
            <Title>Jobs</Title>
          </Grid>
          <Grid item xs={3}>
            <Grid container justify='flex-end' >
              <Grid item>
                <Button variant='outlined' color='primary' component={ReactStaticLinkComponent} to={'/jobs/new'}>
                  New Job
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
        <Grid container spacing={40}>
          <Grid item xs={9}>
            {renderJobsList(this.props)}
          </Grid>
          <Grid item xs={3}>
            {renderSidebar(this.props)}
          </Grid>
        </Grid>
        <Footer />
      </div>
    )
  }
}

Index.propTypes = {
  ethBalance: PropTypes.string,
  linkBalance: PropTypes.string,
  accountBalanceError: PropTypes.string,
  jobCount: PropTypes.number.isRequired,
  jobs: PropTypes.array.isRequired,
  jobsError: PropTypes.string,
  pageSize: PropTypes.number
}

Index.defaultProps = {
  pageSize: 10
}

const mapStateToProps = state => {
  let accountBalanceError
  if (state.accountBalance.networkError) {
    accountBalanceError = 'error fetching balance'
  }
  let jobsError
  if (state.jobs.networkError) {
    jobsError = 'There was an error fetching the jobs. Please reload the page.'
  }

  return {
    ethBalance: state.accountBalance.eth,
    linkBalance: state.accountBalance.link,
    accountBalanceError: accountBalanceError,
    jobCount: state.jobs.count,
    jobs: jobsSelector(state),
    jobsError: jobsError
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchAccountBalance, fetchJobs })
)(Index)

export default ConnectedIndex
