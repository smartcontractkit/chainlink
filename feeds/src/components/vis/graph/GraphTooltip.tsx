import { FeedConfig } from 'config'
import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import { AppState } from 'state'
import { humanizeUnixTimestamp } from 'utils'

export interface Props {
  config: FeedConfig
  tooltip: any
}

interface PositionTypesItem {
  x: number
  y: number
}

interface PositionTypes {
  [key: string]: PositionTypesItem
}

interface PositionStyles {
  left: number
  top: number
}

const positionTypes: PositionTypes = {
  oracle: { x: 10, y: 10 },
  contract: { x: 20, y: 20 },
}

const Tooltip: React.FC<Props> = ({ config, tooltip }) => {
  const [position, setPosition] = useState<Partial<PositionStyles>>()

  useEffect(() => {
    if (tooltip && tooltip.x) {
      setPosition({
        left: tooltip.x + positionTypes[tooltip.type].x,
        top: tooltip.y + positionTypes[tooltip.type].y,
      })
    } else {
      setPosition(undefined)
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
            {tooltip.answerFormatted && (
              <div className="price">
                {config.valuePrefix} {tooltip.answerFormatted}
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
              {config.valuePrefix} {tooltip.latestAnswer || '...'}
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

const mapStateToProps = (state: AppState) => ({
  tooltip: state.networkGraph.tooltip,
})

export default connect(mapStateToProps)(Tooltip)
