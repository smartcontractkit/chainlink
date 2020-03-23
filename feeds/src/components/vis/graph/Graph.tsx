import React, { useEffect, useState } from 'react'
import Tooltip from './GraphTooltip'
import Oracle from './GraphOracle'
import Contract from './GraphContract'
import Line from './GraphLine'
import SideDrawer from './GraphSideDrawer'

interface Props {
  latestOraclesState: any[]
  latestAnswer: any
  latestAnswerTimestamp: any
  config: any
  fetchJobId: any
}

export interface Position {
  x: number
  y: number
}

const Graph = ({
  latestOraclesState,
  latestAnswer,
  latestAnswerTimestamp,
  config,
  fetchJobId,
}: Props) => {
  const [oracles, setOracles] = useState<any[]>([])
  const [positions, setPositions] = useState<Position[]>([])
  const [svgSize] = useState({ width: 1200, height: 600 })

  useEffect(() => {
    setPositions(
      getPositions(svgSize.width, svgSize.height, latestOraclesState),
    )
  }, [latestOraclesState, svgSize.width, svgSize.height])

  useEffect(() => {
    setOracles(latestOraclesState)
  }, [latestOraclesState])

  return (
    <div className="vis__wrapper">
      <Tooltip config={config} />
      <svg
        viewBox={`0 0 ${svgSize.width} ${svgSize.height}`}
        width={svgSize.width}
        height={svgSize.height}
        preserveAspectRatio="xMidYMid meet"
      >
        {oracles.map((o: any, i: number) => (
          <Line
            key={o.id}
            data={latestOraclesState[i]}
            position={{
              x1: svgSize.width / 2,
              y1: svgSize.height / 2,
              x2: positions[i].x,
              y2: positions[i].y,
            }}
          />
        ))}

        {oracles.map((o: any, i: number) => (
          <Oracle
            key={o.id}
            position={positions[i]}
            data={latestOraclesState[i]}
            config={config}
          />
        ))}

        <Contract
          latestAnswer={latestAnswer}
          latestAnswerTimestamp={latestAnswerTimestamp}
          position={{ x: svgSize.width / 2, y: svgSize.height / 2 }}
          config={config}
        />
      </svg>

      <SideDrawer config={config} fetchJobId={fetchJobId} />
    </div>
  )
}

function getPositions(
  width: number,
  height: number,
  oracles: any[],
): Position[] {
  return oracles.map((_, i: number) => {
    const angle = (i / (oracles.length / 2)) * Math.PI
    const x = (20 - height / 2) * Math.cos(angle) + width / 2
    const y = (20 - height / 2) * Math.sin(angle) + height / 2
    return { x, y }
  })
}

export default Graph
