import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import { ChainlinkNode } from 'explorer/models'
import { KeyValueList } from '@chainlink/styleguide'

interface Props {
  operator: ChainlinkNode
}

const Operator: React.FC<Props> = ({ operator }) => {
  return (
    // <Paper className={className}>
    <Paper>
      <KeyValueList
        title={operator.name}
        entries={[['foo', 'bar']]}
        showHead
        titleize
      />
    </Paper>
  )
}

export default Operator
