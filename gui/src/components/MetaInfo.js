import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'

const MetaInfo = ({title, value}) => (
  <PaddedCard>
    <Typography gutterBottom variant='headline' component='h2'>
      {title}
    </Typography>
    <Typography variant='body1' color='textSecondary'>
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
