import React from 'react'
import { bindActionCreators, Dispatch } from 'redux'
import { connect } from 'react-redux'
import { receiveSignoutSuccess } from '../actions'
import Flash from '../components/Flash'
import Unhandled from '../components/Notifications/UnhandledError'
import { AppState } from '../connectors/redux/reducers'

interface Notification {
  component?: any
  props?: any
}

interface NotificationProps {
  notifications: Notification[]
}

const Error = ({ notifications }: NotificationProps) => {
  return (
    <Flash error>
      {notifications.map(({ component, props }, i) => {
        if (component) {
          return <p key={i}>{component(props)}</p>
        } else if (props && props.msg) {
          return <p key={i}>{props.msg}</p>
        }

        return (
          <p key={i}>
            <Unhandled />
          </p>
        )
      })}
    </Flash>
  )
}

const Success = ({ notifications }: NotificationProps) => {
  return (
    <Flash success>
      {notifications.map(({ component, props }, i) => (
        <p key={i}>{component(props)}</p>
      ))}
    </Flash>
  )
}

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface OwnProps {}

interface StateProps {
  errors: any[]
  successes: any[]
}

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface DispatchProps {}

interface Props extends OwnProps, StateProps, DispatchProps {}

export const Notifications = ({ errors, successes }: Props) => {
  return (
    <div>
      {errors.length > 0 && <Error notifications={errors} />}
      {successes.length > 0 && <Success notifications={successes} />}
    </div>
  )
}

function mapStateToProps(state: AppState): StateProps {
  return {
    errors: state.notifications.errors,
    successes: state.notifications.successes,
  }
}

function mapDispatchToProps(dispatch: Dispatch): DispatchProps {
  return bindActionCreators({ receiveSignoutSuccess }, dispatch)
}

export const ConnectedNotifications = connect(
  mapStateToProps,
  mapDispatchToProps,
)(Notifications)

export default ConnectedNotifications
