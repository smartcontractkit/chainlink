import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'

const MetaInfo = ({title, value, classes}) => (
  <PaddedCard>
    <Typography gutterBottom variant='headline' component='h2'>
      {title}
    </Typography>
    <Typography variant='display2' color='inherit'>
      {value}
    </Typography>
  </PaddedCard>
)

MetaInfo.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ])
}

export default MetaInfo
