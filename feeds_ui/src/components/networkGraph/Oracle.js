import React from 'react'
import { connect } from 'react-redux'
import { networkGraphOperations } from 'state/ducks/networkGraph'

const Oracle = ({ data, position, options, setTooltip, setDrawer }) => {
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
        {data.responseFormatted && (
          <text className="vis__oracle-label--price" x="20" y="10">
            {options.valuePrefix} {data.responseFormatted}
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
