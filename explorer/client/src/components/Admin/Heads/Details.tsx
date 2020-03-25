import React from 'react'
import { KeyValueList } from '@chainlink/styleguide'
import { Head } from 'explorer/models'

const LOADING_MESSAGE = 'Loading head...'

const buildEntries = (head: Head): [string, string][] => {
  return [
    ['id', head.id.toString()],
    ['parentHash', head.parentHash],
    ['txHash', head.txHash],
    ['number', head.number],
  ]
}

interface Props {
  head: Head | null
}

const Details: React.FC<Props> = ({ head }) => {
  const title: string = head ? head.txHash : LOADING_MESSAGE
  const entries = head ? buildEntries(head) : []

  return (
    <KeyValueList title={title} entries={entries} showHead={false} titleize />
  )
}

export default Details
