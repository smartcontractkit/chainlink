import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'react-redux'
import Grid from '@material-ui/core/Grid'
import Button from 'components/Button'
import Title from 'components/Title'
import List from 'components/Jobs/List'
import Content from 'components/Content'
import BaseLink from 'components/BaseLink'
import matchRouteAndMapDispatchToProps from 'utils/matchRouteAndMapDispatchToProps'
import jobsSelector from 'selectors/jobs'
import { fetchJobs } from 'actions'

export const Index = props => {
  document.title = 'Jobs'
  return (
    <Content>
      <Grid container>
        <Grid item xs={9}>
          <Title>Jobs</Title>
        </Grid>
        <Grid item xs={3}>
          <Grid container justify="flex-end">
            <Grid item>
              <Button variant="secondary" component={BaseLink} to={'/jobs/new'}>
                New Job
              </Button>
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <List
            jobs={props.jobs}
            jobCount={props.jobCount}
            pageSize={props.pageSize}
            fetchJobs={props.fetchJobs}
            history={props.history}
            match={props.match}
          />
        </Grid>
      </Grid>
    </Content>
  )
}
Index.propTypes = {
  jobCount: PropTypes.number.isRequired,
  jobs: PropTypes.array,
  pageSize: PropTypes.number
}

Index.defaultProps = {
  pageSize: 10
}

const mapStateToProps = state => {
  return {
    jobCount: state.jobs.count,
    jobs: jobsSelector(state)
  }
}

export const ConnectedIndex = connect(
  mapStateToProps,
  matchRouteAndMapDispatchToProps({ fetchJobs })
)(Index)

export default ConnectedIndex
