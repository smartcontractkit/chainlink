import React, { useEffect } from 'react'
import { createChart, updateData, updateState } from './networkGraph'
function NetworkGraph({
  networkGraphNodes,
  networkGraphLinks,
  networkGraphData
}) {
  useEffect(() => {
    createChart()
  }, [])

  useEffect(() => {
    updateData(networkGraphNodes, networkGraphLinks)
  }, [networkGraphNodes, networkGraphLinks])

  useEffect(() => {
    updateState(networkGraphData)
  }, [networkGraphData])

  return (
    <>
      <div className="network-graph">
        <div className="network-graph__tooltip oracle">
          <div className="network-graph__tooltip--type"></div>
          <div className="network-graph__tooltip--oracle">
            <div className="network-graph__tooltip--name"></div>
            <div className="network-graph__tooltip--price"></div>
          </div>
          <div className="network-graph__tooltip--oracle-details">
            <div className="network-graph__tooltip--date"></div>
            <div className="network-graph__tooltip--block"></div>
          </div>
        </div>
      </div>
    </>
  )
}

export default NetworkGraph
