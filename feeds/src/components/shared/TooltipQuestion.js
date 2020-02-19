import React from 'react'
import { Icon, Tooltip } from 'antd'

const TooltipQuestion = props => {
  return (
    <Tooltip {...props}>
      <Icon type="question-circle" />
    </Tooltip>
  )
}

export default TooltipQuestion
