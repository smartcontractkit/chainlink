import React from 'react'
import { Route, useRouteMatch } from 'react-router-dom'

import Content from 'components/Content'
import { EditFeedsManagerScreen } from '../../screens/EditFeedsManager/EditFeedsManagerScreen'
import { FeedsManagerScreen } from '../../screens/FeedsManager/FeedsManagerScreen'
import { NewFeedsManagerScreen } from '../../screens/NewFeedsManager/NewFeedsManagerScreen'

export const FeedsManagerPage = function () {
  const { path } = useRouteMatch()

  return (
    <Content>
      <Route exact path={`${path}/new`}>
        <NewFeedsManagerScreen />
      </Route>

      <Route exact path={path}>
        <FeedsManagerScreen />
      </Route>

      <Route exact path={`${path}/edit`}>
        <EditFeedsManagerScreen />
      </Route>
    </Content>
  )
}
