import React from 'react'
import Grid from '@material-ui/core/Grid'
import { KeyValueList } from '../../KeyValueList'
import { OperatorShowData } from '../../../reducers/adminOperatorsShow'

interface Props {
  operatorData: OperatorShowData
  loadingMessage: string
}

const entries = (operatorData: OperatorShowData): [string, string][] => {
  return [
    ['url', operatorData.url || ''],
    ['uptime', operatorData.uptime.toString()],
    ['job runs completed', operatorData.jobCounts.completed.toString()],
    ['job runs errored', operatorData.jobCounts.errored.toString()],
    ['job runs in progress', operatorData.jobCounts.inProgress.toString()],
    ['total job runs', operatorData.jobCounts.total.toString()],
  ]
}

const Operator: React.FC<Props> = ({ operatorData, loadingMessage }) => {
  const title = operatorData ? operatorData.name : loadingMessage
  const _entries = operatorData ? entries(operatorData) : []
  // TODO refactor ^
  return (
    <KeyValueList title={title} entries={_entries} showHead={false} titleize />
  )
}

export default Operator
