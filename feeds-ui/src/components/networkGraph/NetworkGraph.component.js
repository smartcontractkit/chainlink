import React, { useEffect, useRef, useState } from 'react'
import NetworkGraphD3 from './NetworkGraph.d3'

import NodeDetailsModal from './NodeDetailsModal'
import ContractDetailsModal from './ContractDetailsModal'

const NetworkGraph = ({
  networkGraphNodes,
  networkGraphState,
  options,
  pendingAnswerId,
  fetchJobId,
  updateHeight
}) => {
  const [nodeModalVisible, setNodeModalVisile] = useState(false)
  const [nodeModalData, setNodeModalData] = useState()
  const [nodeModalJobId, setNodeModalJobId] = useState()
  const [contractModalVisible, setContractModalVisile] = useState(false)
  const [contractModalData, setContractModalData] = useState()

  let graph = useRef()

  useEffect(() => {
    graph.current = new NetworkGraphD3(options)
    graph.current.onNodeClick = showNodeInfo
    graph.current.onContractClick = showContractInfo
    graph.current.build()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    if (networkGraphNodes) {
      graph.current.updateNodes(networkGraphNodes)
    }
  }, [networkGraphNodes])

  useEffect(() => {
    if (networkGraphState) {
      graph.current.updateState(networkGraphState, pendingAnswerId)
    }
  }, [networkGraphState, pendingAnswerId])

  const showNodeInfo = async node => {
    setNodeModalData(node)
    setNodeModalVisile(true)

    const jobId = await fetchJobId(node.address)
    setNodeModalJobId(jobId)
  }

  const closeNodeInfo = () => {
    setNodeModalVisile(false)
    setNodeModalData({})
    setNodeModalJobId(null)
  }

  const showContractInfo = contract => {
    setContractModalData(contract)
    setContractModalVisile(true)
  }

  const closeContractInfo = () => {
    setContractModalVisile(false)
    setContractModalData({})
  }

  return (
    <>
      <div className="network-graph">
        <svg className="network-graph__svg">
          <g className="network-graph__links"></g>
          <g className="network-graph__nodes"></g>
        </svg>

        <div className="network-graph__tooltip--oracle">
          <div className="type">Oracle</div>
          <div className="oracle">
            <div className="name"></div>
            <div className="price"></div>
          </div>
          <div className="details">
            <div className="date"></div>
            <div className="block"></div>
          </div>
        </div>

        <div className="network-graph__tooltip--contract">
          <div className="type">Contract</div>
          <div className="contract">
            <div className="name">Aggregation Contract</div>
            <div className="price"></div>
          </div>
        </div>
      </div>
      <NodeDetailsModal
        onClose={closeNodeInfo}
        visible={nodeModalVisible}
        data={nodeModalData}
        pendingAnswerId={pendingAnswerId}
        jobId={nodeModalJobId}
        options={options}
      />
      <ContractDetailsModal
        updateHeight={updateHeight}
        onClose={closeContractInfo}
        visible={contractModalVisible}
        data={contractModalData}
        pendingAnswerId={pendingAnswerId}
        options={options}
      />
    </>
  )
}

export default NetworkGraph
