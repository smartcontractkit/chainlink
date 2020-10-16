import React from 'react'
import { connect, MapStateToProps } from 'react-redux'
import { AppState } from 'reducers'
import { Notification } from 'reducers/notifications'
import Flash from 'components/Flash'
import Unhandled from 'components/Notifications/UnhandledError'

export function renderNotification(notification: Notification) {
  if (typeof notification === 'string') {
    return notification
  } else if (notification.component) {
    return notification.component(notification.props)
  }

  return <Unhandled />
}

export function render(notification: Notification, idx: number) {
  let n

  if (typeof notification === 'string') {
    n = notification
  } else if (notification.component) {
    n = notification.component(notification.props)
  } else {
    n = <Unhandled />
  }

  return <p key={idx}>{n}</p>
}

interface NotificationProps {
  notifications: Notification[]
}

const Error: React.FC<NotificationProps> = ({ notifications }) => {
  return <Flash error>{notifications.map(render)}</Flash>
}

const Success: React.FC<NotificationProps> = ({ notifications }) => {
  return <Flash success>{notifications.map(render)}</Flash>
}

/* eslint-disable-next-line @typescript-eslint/no-empty-interface */
interface OwnProps {}

interface StateProps {
  errors: Notification[]
  successes: Notification[]
}

interface Props extends OwnProps, StateProps {}

export const Notifications: React.FC<Props> = ({ errors, successes }) => {
  return (
    <div>
      {errors.length > 0 && <Error notifications={errors} />}
      {successes.length > 0 && <Success notifications={successes} />}
    </div>
  )
}

const mapStateToProps: MapStateToProps<StateProps, OwnProps, AppState> = (
  state,
) => {
  return {
    errors: state.notifications.errors,
    successes: state.notifications.successes,
  }
}

export const ConnectedNotifications = connect(mapStateToProps)(Notifications)

export default ConnectedNotifications
