import React from 'react'
import { RouteComponentProps } from 'react-router'
import { useEffect, useHooks, useState } from 'use-react-hooks'
import { FIRST_PAGE, GenericList } from '../GenericList'
import { ChangePageEvent, Column } from '@chainlink/styleguide'

const buildItems = (bridges: any[]): Column[][] =>
  bridges.map(b => [
    { type: 'link', text: b.name, to: `/bridges/${b.name}` },
    { type: 'text', text: b.url },
    { type: 'text', text: b.confirmations },
    { type: 'text', text: b.minimumContractPayment },
  ])

// CHECKME
interface OwnProps {
  bridges: any[]
  bridgeCount: number
  pageSize: number
  fetching: boolean
  error: string
  fetchBridges: (...args: any[]) => any
}

// CHECKME
type RouteProps = RouteComponentProps<{
  bridgePage: string
}>

type Props = OwnProps & RouteProps

// FIXME - remove unused export?
export const BridgeList = useHooks<Props>(props => {
  const { bridges, bridgeCount, fetchBridges, pageSize } = props
  const [page, setPage] = useState(FIRST_PAGE)
  useEffect(() => {
    const queryPage =
      (props.match && parseInt(props.match.params.bridgePage, 10)) || FIRST_PAGE
    setPage(queryPage - 1)
    fetchBridges(queryPage, pageSize)
  }, [])
  const handleChangePage = (event: ChangePageEvent, page: number) => {
    if (event) {
      setPage(page)
      fetchBridges(page + 1, pageSize)
      if (props.history) props.history.push(`/bridges/page/${page + 1}`)
    }
  }
  return (
    <GenericList
      emptyMsg="You haven’t created any bridges yet. Create a new bridge"
      headers={[
        'Name',
        'URL',
        'Default Confirmations',
        'Minimum Contract Payment',
      ]}
      items={bridges && buildItems(bridges)}
      onChangePage={handleChangePage}
      count={bridgeCount}
      currentPage={page}
      rowsPerPage={pageSize}
    />
  )
})

export default BridgeList
