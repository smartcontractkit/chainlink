import React from 'react'
import { KeyValueList } from '@chainlink/styleguide'
import { HeadShowData } from '../../../reducers/adminHeadsShow'

const LOADING_MESSAGE = 'Loading head...'

interface Props {
  headData: HeadShowData | null
}

const entries = (head: HeadShowData): [string, string][] => {
  return [
    ['id', head.id.toString()],
    ['parentHash', head.parentHash.toString('hex')],
    ['txHash', head.txHash.toString('hex')],
    ['number', head.number.toString()],
  ]
}

const Head: React.FC<Props> = ({ headData }) => {
  const title = headData ? headData.txHash.toString('hex') : LOADING_MESSAGE
  const _entries = headData ? entries(headData) : []
  return (
    <KeyValueList title={title} entries={_entries} showHead={false} titleize />
  )
}

export default Head
