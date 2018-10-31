import React from 'react'
import PropTypes from 'prop-types'
import Typography from '@material-ui/core/Typography'
import PaddedCard from 'components/PaddedCard'
import { Tooltip } from '@material-ui/core'

const withTooltip = ({tooltip, value}) => (
  <Tooltip title={tooltip} placement='left'>
    <span>{value}</span>
  </Tooltip>
)

const MetaInfo = props => (
  <PaddedCard>
    <Typography gutterBottom variant='headline' component='h2'>
      {props.title}
    </Typography>
    <Typography variant='body1' color='textSecondary'>
      {props.tooltip !== undefined ? withTooltip(props) : props.value}
    </Typography>
  </PaddedCard>
)

MetaInfo.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ]),
  tooltip: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.number
  ])
}

export default MetaInfo
