import React from 'react'
import { KeyValueList } from '@chainlink/styleguide'

const LOADING_MESSAGE = 'Loading head...'

type AttrName = string

const buildEntries = (head: any): [AttrName, string][] => {
  return [
    ['id', head.id.toString()],
    ['parentHash', head.parentHash.toString('hex')],
    ['txHash', head.txHash.toString('hex')],
    ['number', head.number.toString()],
  ]
}

interface Props {
  head: any | null
}

const Head: React.FC<Props> = ({ head }) => {
  const title: string = head ? head.txHash.toString('hex') : LOADING_MESSAGE
  const entries = head ? buildEntries(head) : []

  return (
    <KeyValueList title={title} entries={entries} showHead={false} titleize />
  )
}

export default Head
