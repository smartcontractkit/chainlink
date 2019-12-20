import React from 'react'
import PropTypes from 'prop-types'
import { FIRST_PAGE, GenericList } from '../GenericList'
import { useHooks, useState, useEffect } from 'use-react-hooks'

const buildItems = txs => {
  return txs.map(tx => {
    return [
      { type: 'link', text: tx.hash, to: `/transactions/${tx.hash}` },
      { type: 'text', text: tx.from },
      { type: 'text', text: tx.to },
      { type: 'text', text: tx.nonce },
    ]
  })
}
export const List = useHooks(props => {
  const { transactions, count, fetchTransactions, pageSize } = props
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    const queryPage =
      (props.match && parseInt(props.match.params.transactionsPage, 10)) ||
      FIRST_PAGE
    setPage(queryPage - 1)
    fetchTransactions(queryPage, pageSize)
  }, [])
  const handleChangePage = (e, page) => {
    if (e) {
      setPage(page)
      fetchTransactions(page + 1, pageSize)
      if (props.history) props.history.push(`/transactions/page/${page + 1}`)
    }
  }
  return (
    <GenericList
      emptyMsg="You haven't created any transactions yet"
      headers={['Hash', 'From', 'To', 'Nonce']}
      items={transactions && buildItems(transactions)}
      onChangePage={handleChangePage}
      count={count}
      currentPage={page}
    />
  )
})

List.propTypes = {
  count: PropTypes.number.isRequired,
  pageSize: PropTypes.number.isRequired,
  fetchTransactions: PropTypes.func.isRequired,
  transactions: PropTypes.array,
  error: PropTypes.string,
}

export default List
