import React from 'react'
import { Redirect, useLocation, useHistory } from 'react-router-dom'

import { v2 } from 'api'
import { EditFeedsManagerView } from './EditFeedsManagerView'
import { FormValues } from 'components/Forms/FeedsManagerForm'

import { useFetchFeedsManagers } from 'src/hooks/useFetchFeedsManager'

export const EditFeedsManagerScreen: React.FC = () => {
  const history = useHistory()
  const location = useLocation()
  const { data, loading, error, refetch } = useFetchFeedsManagers()

  if (loading) {
    return null
  }

  if (error) {
    return <div>error</div>
  }

  // We currently only support a single feeds manager, but plan to support more
  // in the future.
  const manager =
    data != undefined && data.feedsManagers.results[0]
      ? data.feedsManagers.results[0]
      : undefined

  if (!manager) {
    return (
      <Redirect
        to={{
          pathname: '/feeds_manager/new',
          state: { from: location },
        }}
      />
    )
  }

  const handleSubmit = async (values: FormValues) => {
    try {
      await v2.feedsManagers.updateFeedsManager(manager.id, values)

      refetch()
      history.push('/feeds_manager')
    } catch (e) {
      console.log(e)
    }
  }

  return <EditFeedsManagerView data={manager} onSubmit={handleSubmit} />
}
