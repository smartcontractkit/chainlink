import React, { useEffect, useState } from 'react'
import Tooltip from './Tooltip'
import Oracle from './Oracle'
import Contract from './Contract'
import Line from './Line'
import SideDrawer from './SideDrawer'

const getPositions = (width, height, oracles) => {
  return oracles.map((_, i) => {
    const angle = (i / (oracles.length / 2)) * Math.PI
    const x = (20 - height / 2) * Math.cos(angle) + width / 2
    const y = (20 - height / 2) * Math.sin(angle) + height / 2
    return { x, y }
  })
}

const AggregationGraph = ({
  oraclesData,
  currentAnswer,
  updateHeight,
  options,
  fetchJobId,
}) => {
  const [oracles, setOracles] = useState([])
  const [positions, setPositions] = useState([])
  const [svgSize] = useState({ width: 1200, height: 600 })

  useEffect(() => {
    setPositions(getPositions(svgSize.width, svgSize.height, oraclesData))
  }, [oraclesData, svgSize.width, svgSize.height])

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
        {oracles.map((o, i) => (
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

        {oracles.map((o, i) => (
          <Oracle
            key={o.id}
            index={i}
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

export default AggregationGraph
