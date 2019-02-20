import React from 'react'
import { bindActionCreators } from 'redux'
import { connect } from 'react-redux'
import { receiveSignoutSuccess } from 'actions'
import Flash from 'components/Flash'
import Unhandled from 'components/Notifications/UnhandledError'

export const Notifications = props => {
  const { errors, successes } = props

  return (
    <div>
      {errors.length > 0 && (
        <Flash error>
          {errors.map(({ component, props }, i) => {
            if (component) return <p key={i}>{component(props)}</p>
            return (
              <p key={i}>
                <Unhandled />
              </p>
            )
          })}
        </Flash>
      )}
      {successes.length > 0 && (
        <Flash success>
          {successes.map(({ component, props }, i) => (
            <p key={i}>{component(props)}</p>
          ))}
        </Flash>
      )}
    </div>
  )
}

const mapStateToProps = state => ({
  errors: state.notifications.errors,
  successes: state.notifications.successes
})

const mapDispatchToProps = dispatch =>
  bindActionCreators({ receiveSignoutSuccess }, dispatch)

export const ConnectedNotifications = connect(
  mapStateToProps,
  mapDispatchToProps
)(Notifications)

export default ConnectedNotifications
