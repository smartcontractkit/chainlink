import React from 'react'
import Paper from '@material-ui/core/Paper'
import Hidden from '@material-ui/core/Hidden'
import Table, { Props as TableProps } from '../Table'
import { LinkColumn, TextColumn, TimeAgoColumn } from '../Table/TableCell'
import { JobRun } from 'explorer/models'

export interface Props {
  loading: boolean
  error: boolean
  currentPage: number
  onChangePage: TableProps['onChangePage']
  jobRuns?: JobRun[]
  rowsPerPage: number
  count?: number
  loadingMsg: string
  emptyMsg: string
  className?: string
}

const DEFAULT_HEADERS = ['Node', 'Run ID', 'Job ID', 'Created At']
const MOBILE_HEADERS = ['Run ID', 'Job ID', 'Created At', 'Node']

const List: React.FC<Props> = props => {
  const tableProps = {
    currentPage: props.currentPage,
    rowsPerPage: props.rowsPerPage,
    count: props.count,
    onChangePage: props.onChangePage,
    loading: props.loading,
    error: props.error,
    loadingMsg: props.loadingMsg,
    emptyMsg: props.emptyMsg,
  }

  return (
    <Paper className={props.className}>
      <Hidden xsDown>
        <Table
          headers={DEFAULT_HEADERS}
          rows={defaultRows(props.jobRuns)}
          {...tableProps}
        />
      </Hidden>
      <Hidden smUp>
        <Table
          headers={MOBILE_HEADERS}
          rows={mobileRows(props.jobRuns)}
          {...tableProps}
        />
      </Hidden>
    </Paper>
  )
}

function buildNodeCol(jobRun: JobRun): TextColumn {
  return { type: 'text', text: jobRun.chainlinkNode.name }
}

function buildIdCol(jobRun: JobRun): LinkColumn {
  return {
    type: 'link',
    text: jobRun.runId,
    to: `/job-runs/${jobRun.id}`,
  }
}

function buildJobIdCol(jobRun: JobRun): TextColumn {
  return { type: 'text', text: jobRun.jobId }
}

function buildCreatedAtCol(jobRun: JobRun): TimeAgoColumn {
  return {
    type: 'time_ago',
    text: jobRun.createdAt,
  }
}

function mobileRows(
  jobRuns?: JobRun[],
): Array<[LinkColumn, TextColumn, TimeAgoColumn, TextColumn]> | undefined {
  return jobRuns?.map((r: JobRun) => {
    return [
      buildIdCol(r),
      buildJobIdCol(r),
      buildCreatedAtCol(r),
      buildNodeCol(r),
    ]
  })
}

function defaultRows(
  jobRuns?: JobRun[],
): Array<[TextColumn, LinkColumn, TextColumn, TimeAgoColumn]> | undefined {
  return jobRuns?.map((r: JobRun) => {
    return [
      buildNodeCol(r),
      buildIdCol(r),
      buildJobIdCol(r),
      buildCreatedAtCol(r),
    ]
  })
}

export default List
