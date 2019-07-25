import React from 'react'
import MuiTableCell from '@material-ui/core/TableCell'
import Link from '../Link'
import TimeAgo from '../TimeAgo'

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
export type Column = TextColumn | LinkColumn | TimeAgoColumn

interface IProps {
  column: Column
}

const renderCol = (col: Column) => {
  switch (col.type) {
    case 'link':
      return col.to.startsWith('http') ? (
        <a href={col.to}>{col.text}</a>
      ) : (
        <Link to={col.to}>{col.text}</Link>
      )
    case 'time_ago':
      return <TimeAgo tooltip>{col.text}</TimeAgo>
    case 'text':
      return col.text
  }
}

const TableCell = ({ column }: IProps) => {
  return <MuiTableCell>{renderCol(column)}</MuiTableCell>
}

export default TableCell
