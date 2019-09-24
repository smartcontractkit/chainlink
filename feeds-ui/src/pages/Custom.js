import React, { useEffect, useState } from 'react'
import { compose } from 'recompose'
import { connect } from 'react-redux'
import { withRouter } from 'react-router'
import queryString from 'query-string'
import { message } from 'antd'

import { aggregationOperations } from 'state/ducks/aggregation'

import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'

const CustomPage = ({ initContract, clearState, history }) => {
  const [urlOptions] = useState(queryString.parse(history.location.search))

  useEffect(() => {
    async function init() {
      try {
        await initContract(urlOptions)
      } catch (error) {
        message.error('Error! Something went wrong', 10000)
      }
    }
    init()

    return () => {
      clearState()
      message.destroy()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <div className="page-wrapper network-page">
      <NetworkGraph options={urlOptions} />
      <NetworkGraphInfo options={urlOptions} />
    </div>
  )
}

const mapStateToProps = state => ({})

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState
}

export default compose(
  connect(
    mapStateToProps,
    mapDispatchToProps
  ),
  withRouter
)(CustomPage)
