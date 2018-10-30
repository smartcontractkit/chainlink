import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'
import { Tooltip } from '@material-ui/core';


const MetaInfo = ({title, value, unformattedValue}) => (
  <PaddedCard>
    <Typography gutterBottom variant='headline' component='h2'>
      {title}
    </Typography>
    <Typography variant='body1' color='textSecondary'>
      <Tooltip title={unformattedValue} placement='left'>
        <div>
          {value}
        </div>
      </Tooltip>
    </Typography>
  </PaddedCard>
)

MetaInfo.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ]),
  unformattedValue: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ])
}

export default MetaInfo
