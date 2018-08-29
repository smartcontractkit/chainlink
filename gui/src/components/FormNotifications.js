import React, { Fragment } from 'react'
import Flash from './Flash'
import { Link } from 'react-static'

export default class FormNotifications extends React.Component {
  render () {
    const { success, classes, jobOrBridge } = this.props
    return (
      <Fragment>
        {JSON.stringify(success) !== '{}' && (
          <Flash success className={classes.flash}>
            {jobOrBridge === 'Bridge' && (
              <Fragment>
                Bridge <Link to={`/bridges/${success.name}`}>{success.name}</Link> was successfully created
              </Fragment>
            )}
            {jobOrBridge === 'Job' && (
              <Fragment>
                Job
                <Link to={`/job_specs/${success.id}`}>{success.id}</Link> was successfully created
              </Fragment>
            )}
          </Flash>
        )}
      </Fragment>
    )
  }
}
