import { AggregatorVis } from 'components/aggregatorVis'
import { AnswerHistory } from 'components/answerHistory'
import { DeviationHistory } from 'components/deviationHistory'
import { OracleTable } from 'components/oracleTable'
import { FeedConfig } from 'config'
import React, { useEffect } from 'react'
import { connect } from 'react-redux'
import { aggregatorOperations } from 'state/ducks/aggregator'

interface OwnProps {
  config: FeedConfig
}

interface DispatchProps {
  initContract: any
  clearContract: any
}

interface Props extends OwnProps, DispatchProps {}

const Aggregator: React.FC<Props> = ({
  initContract,
  config,
  clearContract,
}) => {
  useEffect(() => {
    try {
      initContract(config)
    } catch (error) {
      console.error('Could not initiate contract:', error)
    }
    return clearContract
  }, [initContract, clearContract, config])

  const history = config.history && [
    <AnswerHistory key="answerHistory" config={config} />,
    <DeviationHistory key="deviationHistory" config={config} />,
  ]

  return (
    <>
      <AggregatorVis config={config} />
      {history}
      <OracleTable />
    </>
  )
}

const mapDispatchToProps = {
  initContract: aggregatorOperations.initContract,
  clearContract: aggregatorOperations.clearContract,
}

export default connect(null, mapDispatchToProps)(Aggregator)
