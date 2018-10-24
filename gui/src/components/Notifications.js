import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { Link } from 'react-static'
import { withStyles } from '@material-ui/core/styles'
import { receiveSignoutSuccess } from 'actions'
import Flash from './Flash'
import Unhandled from 'components/Errors/Unhandled'

const styles = theme => ({
  flash: {
    textAlign: 'center'
  }
})

export class Notifications extends React.Component {
  errorPresenter = (err, i) => {
    return <p key={i}>{err.detail}</p>
  }

  successPresenter = (success, i) => {
    const isJobRun = success => success.data && success.data.type === 'runs'
    const attributes = success.data.attributes

    if (isJobRun(success)) {
      return <p key={i}>
        Job <Link to={`/jobs/${attributes.jobId}/runs/id/${attributes.id}`}>{attributes.id}</Link> was successfully run
      </p>
    }
  }

  render () {
    const { errors, successes, classes } = this.props
    return (
      <React.Fragment>
        {errors.length > 0 &&
          <Flash error className={classes.flash}>
            {errors.map((err, i) => {
              if (err.type === 'component' && err.component) {
                return <p key={i}>{err.component(err.props)}</p>
              } else if (err.type === 'component') {
                return <p key={i}><Unhandled /></p>
              }
              return this.errorPresenter(err, i)
            })}
          </Flash>
        }
        {successes.length > 0 && (
          <Flash success className={classes.flash}>
            {successes.map((succ, i) => {
              if (succ.type === 'component') {
                return <p key={i}>{succ.component(succ.props)}</p>
              }
              return this.successPresenter(succ, i)
            })}
          </Flash>
        )}
      </React.Fragment>
    )
  }
}

const mapStateToProps = state => ({
  errors: state.notifications.errors,
  successes: state.notifications.successes
})

const mapDispatchToProps = dispatch => bindActionCreators(
  {receiveSignoutSuccess},
  dispatch
)

export const ConnectedNotifications = connect(
  mapStateToProps,
  mapDispatchToProps
)(Notifications)

export default withStyles(styles)(ConnectedNotifications)
