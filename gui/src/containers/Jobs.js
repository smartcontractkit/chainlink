import React, { Component } from 'react'
import PropType from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import JobList from 'components/JobList'
import TokenBalance from 'components/TokenBalance'
import MetaInfo from 'components/MetaInfo'
import { withSiteData } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import { connect } from 'react-redux'
import { bindActionCreators } from 'redux'
import { fetchJobs, fetchAccountBalance } from 'actions'

const styles = theme => ({
  title: {
    marginTop: theme.spacing.unit * 5,
    marginBottom: theme.spacing.unit * 5
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
        <Typography variant='display2' color='inherit' className={classes.title}>
          Jobs
        </Typography>

        <Grid container spacing={40}>
          <Grid item xs={9}>
            {renderJobsList(this.props)}
          </Grid>
          <Grid item xs={3}>
            {renderSidebar(this.props)}
          </Grid>
        </Grid>
      </div>
    )
  }
}

Jobs.propTypes = {
  classes: PropType.object.isRequired,
  ethBalance: PropType.string.isRequired,
  linkBalance: PropType.string.isRequired,
  accountBalanceFetching: PropType.bool.isRequired,
  accountBalanceError: PropType.string,
  jobCount: PropType.number.isRequired,
  jobs: PropType.array.isRequired,
  jobsFetching: PropType.bool.isRequired,
  jobsError: PropType.string,
  pageSize: PropType.number
}

Jobs.defaultProps = {
  pageSize: 10
}

const jobsSelector = (state) => state.jobs.currentPage.map(id => state.jobs.items[id])

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

const mapDispatchToProps = (dispatch) => {
  return bindActionCreators({
    fetchAccountBalance,
    fetchJobs
  }, dispatch)
}

export const ConnectedJobs = connect(mapStateToProps, mapDispatchToProps)(Jobs)

export default withSiteData(
  withStyles(styles)(ConnectedJobs)
)
