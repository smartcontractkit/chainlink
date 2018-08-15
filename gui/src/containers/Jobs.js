import React, { Component } from 'react'
import PropTypes from 'prop-types'
import Grid from '@material-ui/core/Grid'
import Typography from '@material-ui/core/Typography'
import JobList from 'components/JobList'
import TokenBalance from 'components/TokenBalance'
import MetaInfo from 'components/MetaInfo'
import Footer from 'components/Footer'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
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

const renderJobsList = (props) => {
  const {jobs, jobCount, pageSize, jobsError, fetchJobs, history, match} = props
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

const renderSidebar = ({
  ethBalance,
  linkBalance,
  jobCount,
  accountBalanceError
}) => (
  <Grid container spacing={24}>
    <Grid item xs={12}>
      <TokenBalance
        title='Ethereum'
        value={ethBalance}
        error={accountBalanceError}
      />
    </Grid>
    <Grid item xs={12}>
      <TokenBalance
        title='Link'
        value={linkBalance}
        error={accountBalanceError}
      />
    </Grid>
    <Grid item xs={12}>
      <MetaInfo title='Jobs' value={jobCount} />
    </Grid>
  </Grid>
)

const CreateButton = ({construct, buttonText}) => (
  <Button variant='outlined' color='primary' component={Link} to={`/create/${construct}`}>
    {buttonText}
  </Button>
)

export class Jobs extends Component {
  componentDidMount () {
    this.props.fetchAccountBalance()
  }

  render () {
    const { classes } = this.props
    return (
      <div>
        <Grid container spacing={8} xs={9} alignItems='center' >
          <Grid item xs>
            <Typography variant='display2' color='inherit' className={classes.title}>
              Jobs
            </Typography>
          </Grid>
          <Grid item>
            <CreateButton construct='bridge' buttonText='Create Bridge' />
          </Grid>
          <Grid item>
            <CreateButton construct='job' buttonText='Create Job' />
          </Grid>
        </Grid>
        <Grid container spacing={40}>
          <Grid item xs={9}>{renderJobsList(this.props)}</Grid>
          <Grid item xs={3}>{renderSidebar(this.props)}</Grid>
        </Grid>
        <Footer />
      </div>
    )
  }
}

Jobs.propTypes = {
  classes: PropTypes.object.isRequired,
  ethBalance: PropTypes.string,
  linkBalance: PropTypes.string,
  accountBalanceError: PropTypes.string,
  jobCount: PropTypes.number.isRequired,
  jobs: PropTypes.array.isRequired,
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
    accountBalanceError: accountBalanceError,
    jobCount: state.jobs.count,
    jobs: jobsSelector(state),
    jobsError: jobsError
  }
}

export const ConnectedJobs = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({fetchAccountBalance, fetchJobs})
)(Jobs)

export default withStyles(styles)(ConnectedJobs)
