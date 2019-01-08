import React from 'react'
import PropTypes from 'prop-types'
import Card from '@material-ui/core/Card'
import Table from '@material-ui/core/Table'
import TableBody from '@material-ui/core/TableBody'
import TableCell from '@material-ui/core/TableCell'
import TableHead from '@material-ui/core/TableHead'
import TableRow from '@material-ui/core/TableRow'
import TablePagination from '@material-ui/core/TablePagination'
import Typography from '@material-ui/core/Typography'
import Link from 'components/Link'
import TableButtons, { FIRST_PAGE } from 'components/TableButtons'
import { useHooks, useState, useEffect } from 'use-react-hooks'

const renderBody = (transactions, error) => {
  if (error) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          {error}
        </TableCell>
      </TableRow>
    )
  } else if (transactions && transactions.length === 0) {
    return (
      <TableRow>
        <TableCell component='th' scope='row' colSpan={3}>
          You haven't created any transactions yet.
        </TableCell>
      </TableRow>
    )
  } else if (transactions) {
    return transactions.map(j => (
      <TableRow key={j.hash}>
        <TableCell component='th' scope='row'>
          <Link to={`/transactions/${j.hash}`}>{j.hash}</Link>
        </TableCell>
        <TableCell>
          <Typography variant='body1'>{j.txId}</Typography>
        </TableCell>
      </TableRow>
    ))
  }

  return (
    <TableRow>
      <TableCell component='th' scope='row' colSpan={3}>
        Loading...
      </TableCell>
    </TableRow>
  )
}

export const List = useHooks(props => {
  const [ page, setPage ] = useState(FIRST_PAGE)
  useEffect(() => {
    const queryPage = (props.match && parseInt(props.match.params.transactionsPage, 10)) || FIRST_PAGE
    setPage(queryPage)
    fetchTransactions(queryPage, pageSize)
  }, [])
  const { transactions, count, fetchTransactions, pageSize, error } = props
  const handleChangePage = (e, page) => {
    fetchTransactions(page, pageSize)
    setPage(page)
  }
  const TableButtonsWithProps = () => (
    <TableButtons
      {...props}
      count={count}
      onChangePage={handleChangePage}
      rowsPerPage={pageSize}
      page={page}
      replaceWith={`/transactions/page`}
    />
  )

  return (
    <Card>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>
              <Typography variant='body1' color='textSecondary'>Hash</Typography>
            </TableCell>
            <TableCell>
              <Typography variant='body1' color='textSecondary'>Nonce</Typography>
            </TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {renderBody(transactions, error)}
        </TableBody>
      </Table>
      <TablePagination
        component='div'
        count={count}
        rowsPerPage={pageSize}
        rowsPerPageOptions={[pageSize]}
        page={page - 1}
        onChangePage={() => { } /* handler required by component, so make it a no-op */}
        onChangeRowsPerPage={() => { } /* handler required by component, so make it a no-op */}
        ActionsComponent={TableButtonsWithProps}
      />
    </Card>
  )
})

List.propTypes = {
  count: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  fetchTransactions: PropTypes.func.isRequired,
  transactions: PropTypes.array,
  error: PropTypes.string
}

export default List
