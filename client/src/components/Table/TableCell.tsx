import React from 'react'
import MuiTableCell from '@material-ui/core/TableCell'
import Link from '../Link'

export interface TextColumn {
  type: 'text'
  text: string | number
}
export interface LinkColumn {
  type: 'link'
  text: string | number
  to: string
}
export type Column = TextColumn | LinkColumn

interface IProps {
  column: Column
}

const renderCol = (col: Column) => {
  switch (col.type) {
    case 'link':
      return <Link to={col.to}>{col.text}</Link>
    case 'text':
      return col.text
  }
}

const TableCell = ({ column }: IProps) => {
  return <MuiTableCell>{renderCol(column)}</MuiTableCell>
}

export default TableCell
