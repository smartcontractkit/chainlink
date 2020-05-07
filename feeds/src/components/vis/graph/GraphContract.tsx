import { FeedConfig } from 'config'
import React from 'react'
import { connect } from 'react-redux'
import { networkGraphOperations } from 'state/ducks/networkGraph'
import { Position } from './Graph'

interface OwnProps {
  config: FeedConfig
  latestAnswer: any
  latestAnswerTimestamp: any
  position: Position
}

interface DispatchProps {
  setTooltip: any
  setDrawer: any
}

interface Props extends OwnProps, DispatchProps {}

const Contract: React.FC<Props> = ({
  latestAnswer,
  latestAnswerTimestamp,
  position,
  config,
  setTooltip,
  setDrawer,
}) => {
  return (
    <g
      className="vis__contract-group"
      transform={`translate(${position.x},${position.y})`}
    >
      <circle
        onClick={() =>
          setDrawer({
            type: 'contract',
            latestAnswerTimestamp,
            latestAnswer,
          })
        }
        onMouseEnter={() =>
          setTooltip({
            type: 'contract',
            latestAnswerTimestamp,
            latestAnswer,
            ...position,
          })
        }
        onMouseLeave={() => setTooltip(null)}
        className={`vis__contract ${
          latestAnswer ? 'vis__contract--fulfilled' : 'vis__contract--pending'
        }`}
        r="60"
      />
      <g transform="translate(-15,-35)">
        <path
          d="M866.9 169.9L527.1 54.1C523 52.7 517.5 52 512 52s-11 .7-15.1 2.1L157.1 169.9c-8.3 2.8-15.1 12.4-15.1 21.2v482.4c0 8.8 5.7 20.4 12.6 25.9L499.3 968c3.5 2.7 8 4.1 12.6 4.1s9.2-1.4 12.6-4.1l344.7-268.6c6.9-5.4 12.6-17 12.6-25.9V191.1c.2-8.8-6.6-18.3-14.9-21.2zM810 654.3L512 886.5 214 654.3V226.7l298-101.6 298 101.6v427.6zm-405.8-201c-3-4.1-7.8-6.6-13-6.6H336c-6.5 0-10.3 7.4-6.5 12.7l126.4 174a16.1 16.1 0 0 0 26 0l212.6-292.7c3.8-5.3 0-12.7-6.5-12.7h-55.2c-5.1 0-10 2.5-13 6.6L468.9 542.4l-64.7-89.1z"
          transform="scale(0.03)"
        ></path>
      </g>
      <g className="vis__contract-label">
        <text
          className="vis__contract-label--answer"
          y="15"
          textAnchor="middle"
        >
          {latestAnswer
            ? `${config.valuePrefix} ${latestAnswer}`
            : 'Loading...'}
        </text>
      </g>
    </g>
  )
}

const mapDispatchToProps = {
  setTooltip: networkGraphOperations.setTooltip,
  setDrawer: networkGraphOperations.setDrawer,
}

export default connect(null, mapDispatchToProps)(Contract)
