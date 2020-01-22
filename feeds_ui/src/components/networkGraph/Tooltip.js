import React, { useEffect, useState } from 'react'
import { humanizeUnixTimestamp } from 'utils'
import { connect } from 'react-redux'

const positionTypes = {
  oracle: { x: 10, y: 10 },
  contract: { x: 20, y: 20 },
}

const Tooltip = ({ options, tooltip }) => {
  const [position, setPosition] = useState({})

  useEffect(() => {
    if (tooltip && tooltip.x) {
      setPosition({
        left: tooltip.x + positionTypes[tooltip.type].x,
        top: tooltip.y + positionTypes[tooltip.type].y,
      })
    } else {
      setPosition(null)
    }
  }, [tooltip])

  if (!position || !tooltip) {
    return null
  }

  return (
    <div className="vis__tooltip" style={position}>
      {tooltip.type === 'oracle' && (
        <>
          <div className="type">Oracle</div>
          <div className="data">
            <div className="name">{tooltip.name}</div>
            {tooltip.responseFormatted && (
              <div className="price">
                {options.valuePrefix} {tooltip.responseFormatted}
              </div>
            )}
          </div>
        </>
      )}
      {tooltip.type === 'contract' && (
        <>
          <div className="type">Smart Contract</div>
          <div className="data">
            <div className="price">
              {options.valuePrefix} {tooltip.currentAnswer || '...'}
            </div>
          </div>
        </>
      )}
      {tooltip.meta && (
        <div className="meta">
          <div className="date">
            {humanizeUnixTimestamp(tooltip.meta.timestamp)}
          </div>
          <div className="block">Block: {tooltip.meta.blockNumber}</div>
          <div className="gas">Gas Price (Gwei): {tooltip.meta.gasPrice}</div>
        </div>
      )}
    </div>
  )
}

const mapStateToProps = state => ({
  tooltip: state.networkGraph.tooltip,
})

export default connect(mapStateToProps)(Tooltip)
