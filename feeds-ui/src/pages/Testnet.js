import React, { useEffect, useState } from 'react'
import { compose } from 'recompose'
import { connect } from 'react-redux'

import { aggregationOperations } from 'state/ducks/aggregation'

import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'

const OPTIONS = {
  contractAddress: '0x1c44616CdB7FAe1ba69004ce6010248147CE019e',
  name: 'BTC / USD',
  valuePrefix: '$',
  network: 'ropsten',
}

const NetworkPage = ({ initContract, clearState }) => {
  const [init, setInit] = useState()

  useEffect(() => {
    async function init() {
      await initContract(OPTIONS)
      setInit(true)
    }
    init()

    return () => {
      clearState()
    }
  }, [initContract, setInit, clearState])

  return (
    <div className="page-wrapper network-page">
      {init && <NetworkGraph options={OPTIONS} />}
      <NetworkGraphInfo options={OPTIONS} />
    </div>
  )
}

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState,
}

export default compose(connect(null, mapDispatchToProps))(NetworkPage)
