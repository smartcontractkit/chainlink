import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import Table, { ChangePageEvent } from '../../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../../Table/TableCell'
import { ChainlinkNode } from 'explorer/models'

const HEADERS = ['Name']

// function buildNameCol(operator: ChainlinkNode): UrlColumn {
//   return {
//     type: 'link',
//     text: operator.name,
//     to: `/admin/operators/${operator.id}`, // TODO
//   }
// }

// type UrlColumn = LinkColumn | TextColumn

// function buildUrlCol(operator: ChainlinkNode): UrlColumn {
//   if (operator.url) {
//     return {
//       type: 'link',
//       text: operator.url,
//       to: operator.url,
//     }
//   }

//   return { type: 'text', text: '-' }
// }

// function buildCreatedAtCol(operator: ChainlinkNode): TimeAgoColumn {
//   return {
//     type: 'time_ago',
//     text: operator.createdAt,
//   }
// }

// function rows(
//   operators?: ChainlinkNode[],
// ): [UrlColumn, UrlColumn, TimeAgoColumn][] | undefined {
//   if (operators) {
//     return operators.map(o => {
//       return [buildNameCol(o), buildUrlCol(o), buildCreatedAtCol(o)]
//     })
//   }
// }

function rows(operator: ChainlinkNode): [TextColumn][] | undefined {
  // if (operators) {
  //   return operators.map(o => {
  //     return [buildNameCol(o), buildUrlCol(o), buildCreatedAtCol(o)]
  //   })
  // }
  return [
    [
      {
        type: 'text',
        text: operator.name,
      },
    ],
  ]
}

interface Props {
  operator: ChainlinkNode
}

const Operator: React.FC<Props> = ({
  operator,
  // count,
  // currentPage,
  // className,
  // onChangePage,
  // loadingMsg,
  // emptyMsg,
}) => {
  return (
    // <Paper className={className}>
    <Paper>
      <Hidden xsDown>
        <Table
          headers={HEADERS}
          currentPage={1}
          rows={rows(operator)}
          // count={count}
          onChangePage={() => {}}
          // loadingMsg={loadingMsg}
          // emptyMsg={emptyMsg}
        />
      </Hidden>
    </Paper>
  )
}

export default Operator
