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

const valAndTooltip = ({value, title, error}) => {
  if (error) {
    return {
      val: error,
      tooltip: 'Error'
    }
  } else if (value == null) {
    return {
      val: '...',
      tooltip: 'Loading...'
    }
  }

  return {
    val: formatBalance(value).formatted,
    tooltip: formatBalance(value).unformatted
  }
}

const TokenBalance = props => {
  const {val, tooltip} = valAndTooltip(props)

  return (
    <MetaInfo
      className={props.className}
      title={props.title}
      value={val}
      tooltip={tooltip}
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
