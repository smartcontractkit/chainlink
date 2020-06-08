import React from 'react'
import { KeyValueList } from '@chainlink/styleguide'
import { OperatorShowData } from '../../../reducers/adminOperatorsShow'

const LOADING_MESSAGE = 'Loading operator...'

interface Props {
  operatorData?: OperatorShowData
}

const entries = (operatorData: OperatorShowData): [string, string][] => {
  return [
    ['url', operatorData.url ?? ''],
    ['uptime', operatorData.uptime.toString()],
    ['job runs completed', operatorData.jobCounts.completed.toString()],
    ['job runs errored', operatorData.jobCounts.errored.toString()],
    ['job runs in progress', operatorData.jobCounts.inProgress.toString()],
    ['total job runs', operatorData.jobCounts.total.toString()],
    ['core version', operatorData.coreVersion ?? 'N/A'],
    ['core sha', operatorData.coreSha ?? 'N/A'],
  ]
}

const Operator: React.FC<Props> = ({ operatorData }) => {
  const title = operatorData ? operatorData.name : LOADING_MESSAGE
  const _entries = operatorData ? entries(operatorData) : []
  return (
    <KeyValueList title={title} entries={_entries} showHead={false} titleize />
  )
}

export default Operator
