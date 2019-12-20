import React from 'react'
import { Table, ChangePageEvent, Column } from '@chainlink/styleguide'

interface Props {
  currentPage: number
  onChangePage: (event: ChangePageEvent, page: number) => void
  headers: string[]
  items?: Column[][]
  count?: number
  rowsPerPage?: number
  emptyMsg?: string
  className?: string
}

export const GenericList = (props: Props) => {
  const items = props.items || []
  const headers = props.headers
  return (
    <Table
      headers={headers}
      currentPage={props.currentPage}
      rowsPerPage={props.rowsPerPage}
      rows={items}
      count={props.count}
      onChangePage={props.onChangePage}
      emptyMsg={props.emptyMsg}
    />
  )
}

export const FIRST_PAGE = 1
