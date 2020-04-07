import { FeedConfig } from 'config'
import React from 'react'
import { connect } from 'react-redux'
import { networkGraphOperations } from 'state/ducks/networkGraph'
import { Position } from './Graph'

interface OwnProps {
  config: FeedConfig
  data: any
  position: Position
}

interface DispatchProps {
  setTooltip: any
  setDrawer: any
}

interface Props extends OwnProps, DispatchProps {}

const Oracle: React.FC<Props> = ({
  data,
  position,
  config,
  setTooltip,
  setDrawer,
}) => {
  if (!data) return null

  return (
    <g
      className="vis__oracle-group"
      transform={`translate(${position.x},${position.y})`}
    >
      <circle
        onClick={() => setDrawer(data)}
        onMouseEnter={() => setTooltip({ ...data, ...position })}
        onMouseLeave={() => setTooltip(null)}
        className={`vis__oracle ${
          data.isFulfilled ? 'vis__oracle--fulfilled' : 'vis__oracle--pending'
        }`}
        r="10"
      />
      <g className="vis__oracle-label">
        <text className="vis__oracle-label--name" x="20" y="-5">
          {data.name}
        </text>
        {data.answerFormatted && (
          <text className="vis__oracle-label--price" x="20" y="10">
            {config.valuePrefix} {data.answerFormatted}
          </text>
        )}
      </g>
    </g>
  )
}

const mapDispatchToProps = {
  setTooltip: networkGraphOperations.setTooltip,
  setDrawer: networkGraphOperations.setDrawer,
}

export default connect(null, mapDispatchToProps)(Oracle)
