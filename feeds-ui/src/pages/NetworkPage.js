import React, { useEffect } from 'react'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import { aggregationOperations } from 'state/ducks/aggregation'

import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'

const NetworkPage = ({ fetchInitData }) => {
  useEffect(() => {
    fetchInitData()
  })
  return (
    <div className="page-wrapper">
      <NetworkGraphInfo />
      <NetworkGraph />
      <AnswerHistory />
    </div>
  )
}

const mapStateToProps = state => ({})

const mapDispatchToProps = {
  fetchInitData: aggregationOperations.fetchInitData
}

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps
  )
)(NetworkPage)
