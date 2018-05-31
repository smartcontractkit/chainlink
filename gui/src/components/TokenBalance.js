import React from 'react'
import PropTypes from 'prop-types'
import MetaInfo from 'components/MetaInfo'
import numeral from 'numeral'
import { BigNumber } from 'bignumber.js'

const WEI_PER_TOKEN = new BigNumber(10 ** 18)

const formatBalance = (val) => {
  const b = new BigNumber(val)
  const tokenBalance = b.dividedBy(WEI_PER_TOKEN).toNumber()

  return numeral(tokenBalance).format('0.20a')
}

const TokenBalance = ({title, value, className, fetching, error}) => {
  let val
  if (fetching) {
    val = '...'
  } else if (error) {
    val = error
  } else {
    val = formatBalance(value)
  }

  return (
    <MetaInfo
      title={title}
      value={val}
      className={className}
    />
  )
}

TokenBalance.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.string,
  className: PropTypes.string,
  fetching: PropTypes.bool,
  error: PropTypes.string
}

export default TokenBalance
