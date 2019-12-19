import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import { ChainlinkNode } from 'explorer/models'
import { KeyValueList } from '../../KeyValueList'

interface Props {
  operator: ChainlinkNode
}

enum operatorProperties {
  url = 'url',
  createdAt = 'createdAt',
}

const OPERATOR_PROPERTIES: operatorProperties[] = [
  operatorProperties.url,
  operatorProperties.createdAt,
]

const entries = (operator: ChainlinkNode): [string, string][] => {
  return OPERATOR_PROPERTIES.map(property => [
    property.toString(),
    operator[property] || '', // TODO defult value?
  ])
}

const Operator: React.FC<Props> = ({ operator }) => {
  const title = operator ? operator.name : 'Loading...'
  const _entries = operator ? entries(operator) : []
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
