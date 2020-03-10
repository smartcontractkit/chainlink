import React from 'react'
import { Icon, Tooltip } from 'antd'
import { TooltipProps } from 'antd/lib/tooltip'

const TooltipQuestion: React.FC<TooltipProps> = props => {
  return (
    <Tooltip {...props}>
      <Icon type="question-circle" />
    </Tooltip>
  )
}

export default TooltipQuestion
