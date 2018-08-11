import React, { Fragment } from 'react'
import Flash from './Flash'
import { Link } from 'react-static'

export const BridgeAndJobNotifications = props => {
  const { error, success, networkError, authenticated, classes, jobOrBridge } = props
  return (
    <Fragment>
      {error.length > 0 &&
        authenticated && (
          <Flash error className={classes.flash}>
            {(Array.isArray(error) && error.map((msg, i) => <span key={i}>{msg}</span>)) || error}
          </Flash>
        )}
      {!authenticated && (
        <Flash warning className={classes.flash}>
          Session expired. <Link to='/signin'>Please sign back in.</Link>
        </Flash>
      )}
      {error.length === 0 &&
        networkError && (
          <Flash error className={classes.flash}>
            Received a Network Error.
          </Flash>
        )}
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
