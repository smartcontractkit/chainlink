import React, { Fragment } from 'react'
import Button from '@material-ui/core/Button'
import { connect } from 'react-redux'
import Flash from './Flash'
import { Link } from 'react-static'
import { receiveSignoutSuccess } from 'actions'
import { bindActionCreators } from 'redux'

export class FormNotifications extends React.Component {
  signOutLocally = () => {
    this.props.receiveSignoutSuccess()
  }

  render () {
    const { errors, success, networkError, authenticated, classes, jobOrBridge } = this.props
    return (
      <Fragment>
        {errors.length > 0 &&
          authenticated && (
            <Flash error className={classes.flash}>
              {errors.map((err, i) => {
                if (err.status === 401) {
                  return (<span key={i}>
                    <span>{err.detail}</span>
                    <Button onClick={this.signOutLocally}>
                      Log back in
                    </Button>
                  </span>)
                }
                return <span key={i}>{err.detail}</span>
              })}
            </Flash>
          )}
        {!authenticated && (
          <Flash warning className={classes.flash}>
            Session expired. <Link to='/signin'>Please sign back in.</Link>
          </Flash>
        )}
        {errors.length === 0 &&
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
}

const mapDispatchToProps = dispatch => bindActionCreators({receiveSignoutSuccess}, dispatch)
export default connect(state => state, mapDispatchToProps)(FormNotifications)
