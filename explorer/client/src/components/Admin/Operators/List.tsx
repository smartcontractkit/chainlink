import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import Table, { ChangePageEvent } from '../../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../../Table/TableCell'

const HEADERS = ['Name', 'URL', 'Created At']

const buildNameCol = (operator: ChainlinkNode): TextColumn => {
  return {
    type: 'text',
    text: operator.name,
  }
}

type UrlColumn = LinkColumn | TextColumn

const buildUrlCol = (operator: ChainlinkNode): UrlColumn => {
  if (operator.url) {
    return {
      type: 'link',
      text: operator.url,
      to: operator.url,
    }
  }

  return { type: 'text', text: '-' }
}

const buildCreatedAtCol = (operator: ChainlinkNode): TimeAgoColumn => {
  return {
    type: 'time_ago',
    text: operator.createdAt,
  }
}

const rows = (
  operators?: ChainlinkNode[],
): [TextColumn, UrlColumn, TimeAgoColumn][] | undefined => {
  if (operators) {
    return operators.map(o => {
      return [buildNameCol(o), buildUrlCol(o), buildCreatedAtCol(o)]
    })
  }
}

interface Props {
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  operators?: ChainlinkNode[]
  count?: number
  loadingMsg: string
  emptyMsg: string
  className?: string
}

const List = ({
  operators,
  count,
  currentPage,
  className,
  onChangePage,
  loadingMsg,
  emptyMsg,
}: Props) => {
  return (
    <Paper className={className}>
      <Hidden xsDown>
        <Table
          headers={HEADERS}
          currentPage={currentPage}
          rows={rows(operators)}
          count={count}
          onChangePage={onChangePage}
          loadingMsg={loadingMsg}
          emptyMsg={emptyMsg}
        />
      </Hidden>
    </Paper>
  )
}

export default List
