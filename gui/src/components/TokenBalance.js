import React from 'react'
import PropType from 'prop-types'
import Typography from '@material-ui/core/Typography'
import Card from '@material-ui/core/Card'
import { BigNumber } from 'bignumber.js'
import numeral from 'numeral'

const WEI_PER_TOKEN = new BigNumber(10 ** 18)

const formatBalance = (val) => {
  const b = new BigNumber(val)
  const tokenBalance = b.dividedBy(WEI_PER_TOKEN).toNumber()

  return numeral(tokenBalance).format('0.20a')
}

const TokenBalance = ({title, value, className}) => (
  <Card className={className}>
    <Typography gutterBottom variant='headline' component='h2'>
      {title}
    </Typography>
    <Typography variant='display2' color='inherit'>
      {formatBalance(value)}
    </Typography>
  </Card>
)

TokenBalance.propTypes = {
  title: PropType.string,
  value: PropType.string,
  className: PropType.string
}

export default TokenBalance
