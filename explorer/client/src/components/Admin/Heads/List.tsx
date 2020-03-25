import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import { join } from 'path'
import Table, { ChangePageEvent } from '../../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../../Table/TableCell'
import { Head } from 'explorer/models'

const HEADERS = ['Block Height', 'Hash', 'Created At']
const LOADING_MSG = 'Loading heads...'
const EMPTY_MSG = 'The Explorer has not yet observed any heads.'

function buildBlockHeightCol(head: Head): TextColumn {
  return {
    type: 'text',
    text: head.number,
  }
}

function buildNameCol(head: Head): UrlColumn {
  return {
    type: 'link',
    text: head.txHash,
    to: join('/', 'admin', 'heads', head.id.toString()),
  }
}

type UrlColumn = LinkColumn | TextColumn

function buildCreatedAtCol(head: Head): TimeAgoColumn {
  return {
    type: 'time_ago',
    text: head.createdAt,
  }
}

function rows(
  heads?: Head[],
): [TextColumn, UrlColumn, TimeAgoColumn][] | undefined {
  return heads?.map(o => {
    return [buildBlockHeightCol(o), buildNameCol(o), buildCreatedAtCol(o)]
  })
}

interface Props {
  loaded: boolean
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  heads?: Head[]
  count?: number
  className?: string
}

const List: React.FC<Props> = ({
  loaded,
  heads,
  count,
  currentPage,
  className,
  onChangePage,
}) => {
  return (
    <Paper className={className}>
      <Hidden xsDown>
        <Table
          loaded={loaded}
          headers={HEADERS}
          currentPage={currentPage}
          rows={rows(heads)}
          count={count}
          onChangePage={onChangePage}
          loadingMsg={LOADING_MSG}
          emptyMsg={EMPTY_MSG}
        />
      </Hidden>
    </Paper>
  )
}

export default List
