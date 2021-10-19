import React from 'react'
import { Redirect, Route, useRouteMatch } from 'react-router-dom'

import { v2 } from 'api'
import Content from 'components/Content'
import { Resource, FeedsManager } from 'core/store/models'
import { useErrorHandler } from 'hooks/useErrorHandler'

import { EditFeedsManagerView } from './EditFeedsManagerView'
import { RegisterFeedsManagerView } from './RegisterFeedsManagerView'
import { FeedsManagerView } from './FeedsManagerView'

export const FeedsManagerScreen: React.FC = () => {
  const { path } = useRouteMatch()
  const { error, ErrorComponent, setError } = useErrorHandler()
  const [manager, setManager] = React.useState<Resource<FeedsManager>>()
  const [isLoading, setIsLoading] = React.useState(true)

  // Fetch the feeds managers.
  //
  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  React.useEffect(() => {
    v2.feedsManagers
      .getFeedsManagers()
      .then((managers) => {
        if (managers.data.length > 0) {
          setManager(managers.data[0])
        }
      })
      .catch(setError)
      .finally(() => {
        setIsLoading(false)
      })
  }, [setError])

  if (isLoading) {
    return null
  }

  if (error) {
    return <ErrorComponent />
  }

  return (
    <Content>
      <Route
        exact
        path={`${path}/new`}
        render={({ location }) =>
          manager ? (
            <Redirect
              to={{
                pathname: '/feeds_manager',
                state: { from: location },
              }}
            />
          ) : (
            <RegisterFeedsManagerView
              onSuccess={(manager) => setManager(manager)}
            />
          )
        }
      />

      <Route
        exact
        path={path}
        render={({ location }) =>
          manager ? (
            <FeedsManagerView manager={manager.attributes} />
          ) : (
            <Redirect
              to={{
                pathname: '/feeds_manager/new',
                state: { from: location },
              }}
            />
          )
        }
      />

      <Route
        exact
        path={`${path}/edit`}
        render={({ location }) =>
          manager ? (
            <EditFeedsManagerView
              manager={manager}
              onSuccess={(manager) => setManager(manager)}
            />
          ) : (
            <Redirect
              to={{
                pathname: '/feeds_manager',
                state: { from: location },
              }}
            />
          )
        }
      />
    </Content>
  )
}
