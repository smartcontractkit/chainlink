import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
/* import { DispatchBinding } from '@chainlink/ts-helpers' */
import { aggregationOperations } from 'state/ducks/aggregation'
import { NetworkGraph } from 'components/networkGraph'
import { NetworkGraphInfo } from 'components/networkGraphInfo'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { Header } from 'components/header'
import { parseQuery, uIntFrom } from 'utils'

interface OwnProps {
  history: any
}

interface StateProps {}

interface DispatchProps {
  initContract: any
  clearState: any
  /* initContract: DispatchBinding<typeof aggregationOperations.initContract> */
  /* clearState: DispatchBinding<typeof aggregationOperations.clearState> */
}

interface Props extends OwnProps, StateProps, DispatchProps {}

const Page: React.FC<Props> = ({ initContract, clearState, history }) => {
  const [options] = useState(formatOptions(parseQuery(history.location.search)))

  useEffect(() => {
    initContract(options).catch(() => {
      console.error('Could not initiate contract')
    })
    return () => {
      clearState()
    }
  }, [initContract, clearState, options])

  return (
    <>
      <div className="page-container-full-width">
        <Header />
      </div>
      <div className="page-wrapper network-page">
        <NetworkGraph options={options} />
        <NetworkGraphInfo options={options} />
        {options && options.history && <AnswerHistory options={options} />}
        {options && options.history && <DeviationHistory options={options} />}
        <OracleTable />
      </div>
    </>
  )
}

const mapDispatchToProps = {
  initContract: aggregationOperations.initContract,
  clearState: aggregationOperations.clearState,
}

function formatOptions(options: any) {
  return {
    ...options,
    networkId: uIntFrom(options.networkId),
    contractVersion: 2,
    decimalPlaces: uIntFrom(options.decimalPlaces),
    counter: uIntFrom(options.counter) || false,
  }
}

export default connect(null, mapDispatchToProps)(Page)
