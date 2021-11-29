import React from 'react'
import { Route, useRouteMatch } from 'react-router-dom'

import { JobProposalScreen } from 'src/screens/JobProposal/JobProposalScreen'

export const JobProposalsPage = function () {
  const { path } = useRouteMatch()

  return (
    <>
      <Route exact path={`${path}/:id`}>
        <JobProposalScreen />
      </Route>
    </>
  )
}
