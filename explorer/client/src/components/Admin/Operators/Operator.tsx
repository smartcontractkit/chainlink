import React from 'react'
import Paper from '@material-ui/core/Paper'
import { KeyValueList } from '../../KeyValueList'
import { OperatorShowData } from '../../../reducers/adminOperatorsShow'

interface Props {
  operatorData: OperatorShowData
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

const Operator: React.FC<Props> = ({ operatorData }) => {
  const title = operatorData ? operatorData.name : 'Loading...'
  const _entries = operatorData ? entries(operatorData) : []
  // TODO refactor ^
  return (
    // <Paper className={className}>
    <Paper>
      <KeyValueList
        title={title}
        entries={_entries}
        showHead={false}
        titleize
      />
    </Paper>
  )
}

export default Operator
