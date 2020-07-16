import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import { join } from 'path'
import Table, { Props as TableProps } from '../../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../../Table/TableCell'
import { ChainlinkNode } from 'explorer/models'

const HEADERS = ['Name', 'URL', 'Created At'] as const
const LOADING_MSG = 'Loading operators...'
const EMPTY_MSG = 'There are no operators added to the Explorer yet.'
const ERROR_MSG = 'Error loading operators.'

interface Props {
  currentPage: number
  onChangePage: TableProps['onChangePage']
  loading: boolean
  error: boolean
  operators?: ChainlinkNode[]
  count?: number
  className?: string
}

const List: React.FC<Props> = ({
  loading,
  error,
  operators,
  count,
  currentPage,
  className,
  onChangePage,
}) => {
  return (
    <Paper className={className}>
      <Hidden xsDown>
        <Table
          headers={HEADERS}
          currentPage={currentPage}
          rows={rows(operators)}
          count={count}
          onChangePage={onChangePage}
          loading={loading}
          error={error}
          loadingMsg={LOADING_MSG}
          emptyMsg={EMPTY_MSG}
          errorMsg={ERROR_MSG}
        />
      </Hidden>
    </Paper>
  )
}

function buildNameCol(operator: ChainlinkNode): UrlColumn {
  return {
    type: 'link',
    text: operator.name,
    to: join('/', 'admin', 'operators', operator.id.toString()),
  }
}

type UrlColumn = LinkColumn | TextColumn

function buildUrlCol(operator: ChainlinkNode): UrlColumn {
  if (operator.url) {
    return {
      type: 'link',
      text: operator.url,
      to: operator.url,
    }
  }

  return { type: 'text', text: '-' }
}

function buildCreatedAtCol(operator: ChainlinkNode): TimeAgoColumn {
  return {
    type: 'time_ago',
    text: operator.createdAt,
  }
}

function rows(
  operators?: ChainlinkNode[],
): [UrlColumn, UrlColumn, TimeAgoColumn][] | undefined {
  return operators?.map(o => {
    return [buildNameCol(o), buildUrlCol(o), buildCreatedAtCol(o)]
  })
}

export default List
