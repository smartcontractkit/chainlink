import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import JobList from 'components/JobList'
import TokenBalance from 'components/TokenBalance'
import MetaInfo from 'components/MetaInfo'
import Footer from 'components/Footer'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import { withSiteData } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { fetchJobs, fetchAccountBalance } from 'actions'
import { jobsSelector } from 'selectors'
import Button from '@material-ui/core/Button'
import Link from 'components/Link'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
  },
  createButtons: {
    padding: theme.spacing.unit * 2
  }
})

const renderJobsList = ({jobs, jobCount, pageSize, jobsFetching, jobsError, fetchJobs}) => (
  <JobList
    jobs={jobs}
    jobCount={jobCount}
    pageSize={pageSize}
    fetching={jobsFetching}
    error={jobsError}
    fetchJobs={fetchJobs}
  />
)

const renderSidebar = ({
  ethBalance,
  linkBalance,
  jobCount,
  accountBalanceFetching,
  accountBalanceError
}) => (
  <Grid container spacing={24}>
    <Grid item xs={12}>
      <TokenBalance
        title='Ethereum'
        fetching={accountBalanceFetching}
        value={ethBalance}
        error={accountBalanceError}
      />
    </Grid>
    <Grid item xs={12}>
      <TokenBalance
        title='Link'
        fetching={accountBalanceFetching}
        value={linkBalance}
        error={accountBalanceError}
      />
    </Grid>
    <Grid item xs={12}>
      <MetaInfo title='Jobs' value={jobCount} />
    </Grid>
  </Grid>
)

export class Jobs extends Component {
  componentDidMount () {
    this.props.fetchJobs(1, this.props.pageSize)
    this.props.fetchAccountBalance()
  }

  render () {
    const { classes } = this.props

    return (
      <div>
        <Grid container alignItems='center' >
          <Grid item xs={1}>
            <Typography variant='display2' color='inherit' className={classes.title}>
              Jobs
            </Typography>
          </Grid>
          <Grid item xs={1}>
            <Button variant='outlined' color='primary' component={Link} to={{ pathname: '/create', state: { tab: 0 } }} >
              Create Bridge
            </Button>
          </Grid>
          <Grid item xs={1}>
            <Button variant='outlined' color='primary' component={Link} to={{ pathname: '/create', state: { tab: 1 } }} >
              Create Job
            </Button>
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

Jobs.propTypes = {
  classes: PropTypes.object.isRequired,
  ethBalance: PropTypes.string.isRequired,
  linkBalance: PropTypes.string.isRequired,
  accountBalanceFetching: PropTypes.bool.isRequired,
  accountBalanceError: PropTypes.string,
  jobCount: PropTypes.number.isRequired,
  jobs: PropTypes.array.isRequired,
  jobsFetching: PropTypes.bool.isRequired,
  jobsError: PropTypes.string,
  pageSize: PropTypes.number
}

Jobs.defaultProps = {
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
    accountBalanceFetching: state.accountBalance.fetching,
    accountBalanceError: accountBalanceError,
    jobCount: state.jobs.count,
    jobs: jobsSelector(state),
    jobsFetching: state.jobs.fetching,
    jobsError: jobsError
  }
}

export const ConnectedJobs = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchAccountBalance, fetchJobs})
)(Jobs)

export default withSiteData(
  withStyles(styles)(ConnectedJobs)
)
