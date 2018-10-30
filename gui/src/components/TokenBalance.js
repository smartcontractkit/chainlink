import React from 'react'
import PropTypes from 'prop-types'
import MetaInfo from 'components/MetaInfo'
import numeral from 'numeral'
import { BigNumber } from 'bignumber.js'

const WEI_PER_TOKEN = new BigNumber(10 ** 18)

const formatBalance = (val) => {
  const b = new BigNumber(val)
  const tokenBalance = b.dividedBy(WEI_PER_TOKEN).toNumber()
  return {formatted: numeral(tokenBalance).format('0.200000a'), unformatted: tokenBalance}
}

const TokenBalance = ({title, value, className, error}) => {
  let val
  let unformattedVal
  if (error) {
    val = error
  } else if (value == null) {
    val = '...'
  } else {
    val = formatBalance(value).formatted
    unformattedVal = formatBalance(value).unformatted
  }
  return (
    <MetaInfo
      title={title}
      value={val}
      unformattedValue={unformattedVal}
      className={className}
    />
  )
}

TokenBalance.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.string,
  className: PropTypes.string,
  error: PropTypes.string
}

export default TokenBalance
