import React, { useEffect, useState } from 'react'
import { connect } from 'react-redux'
import Tooltip from './Tooltip'
import Oracle from './Oracle'
import Contract from './Contract'
import Line from './Line'
import SideDrawer from './SideDrawer'
import {
  aggregationSelectors,
  aggregationOperations,
} from 'state/ducks/aggregation'

interface Props {
  oraclesData: any[]
  currentAnswer: any
  updateHeight: any
  options: any
  fetchJobId: any
}

interface Position {
  x: number
  y: number
}

function NetworkGraph({
  oraclesData,
  currentAnswer,
  updateHeight,
  options,
  fetchJobId,
}: Props) {
  const [oracles, setOracles] = useState<any[]>([])
  const [positions, setPositions] = useState<Position[]>([])
  const [svgSize] = useState({ width: 1200, height: 600 })

  useEffect(() => {
    setPositions(getPositions(svgSize.width, svgSize.height, oraclesData))
  }, [oraclesData, svgSize.width, svgSize.height, setPositions])

  useEffect(() => {
    setOracles(oraclesData)
  }, [oraclesData])

  return (
    <div className="vis__wrapper">
      <Tooltip options={options} />
      <svg
        viewBox={`0 0 ${svgSize.width} ${svgSize.height}`}
        width={svgSize.width}
        height={svgSize.height}
        preserveAspectRatio="xMidYMid meet"
      >
        {oracles.map((o: any, i: number) => (
          <Line
            key={o.id}
            data={oraclesData[i]}
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
            data={oraclesData[i]}
            options={options}
          />
        ))}

        <Contract
          currentAnswer={currentAnswer}
          updateHeight={updateHeight}
          position={{ x: svgSize.width / 2, y: svgSize.height / 2 }}
          options={options}
        />
      </svg>

      <SideDrawer options={options} fetchJobId={fetchJobId} />
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

const mapStateToProps = (state: any) => ({
  updateHeight: state.aggregation.updateHeight,
  oraclesData: aggregationSelectors.oraclesData(state),
  currentAnswer: state.aggregation.currentAnswer,
})

const mapDispatchToProps = {
  fetchJobId: aggregationOperations.fetchJobId,
}

export default connect(mapStateToProps, mapDispatchToProps)(NetworkGraph)
