import React from 'react'
import MuiTableCell from '@material-ui/core/TableCell'
import Link from '../Link'
import { TimeAgo } from '../TimeAgo'
import { StatusText } from './StatusText'

export interface TextColumn {
  type: 'text'
  text: string | number
}
export interface TimeAgoColumn {
  type: 'time_ago'
  text: string
}
export interface LinkColumn {
  type: 'link'
  text: string | number
  to: string
}

export interface StatusColumn {
  type: 'status'
  text: 'complete' | 'error' | 'pending'
}

export type Column = TextColumn | LinkColumn | TimeAgoColumn | StatusColumn

interface Props {
  column: Column
}

const renderCol = (col: Column) => {
  switch (col.type) {
    case 'link':
      return <Link href={col.to}>{col.text}</Link>
    case 'time_ago':
      return <TimeAgo tooltip>{col.text}</TimeAgo>
    case 'text':
      return col.text
    case 'status':
      return <StatusText status={col.text} />
  }
}

export const TableCell = ({ column }: Props) => {
  return <MuiTableCell>{renderCol(column)}</MuiTableCell>
}
