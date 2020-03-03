import React, { useEffect, useState } from 'react'
import { humanizeUnixTimestamp } from 'utils/'
import _ from 'lodash'
import { Table, Icon } from 'antd'

const OracleTable = ({
  networkGraphState,
  networkGraphNodes,
  fetchEthGasPrice,
  ethGasPrice,
}) => {
  const [data, setData] = useState()
  const [gasPrice, setGasPrice] = useState()

  useEffect(() => {
    fetchEthGasPrice()
  }, [fetchEthGasPrice])

  useEffect(() => {
    if (ethGasPrice) {
      setGasPrice(ethGasPrice.fast / 10)
    }
  }, [ethGasPrice])

  useEffect(() => {
    const mergedData = networkGraphNodes
      .filter(node => node.type === 'oracle')
      .map(oracle => {
        const state = _.find(networkGraphState, { sender: oracle.address })
        return {
          oracle,
          state,
          key: oracle.id,
        }
      })
    setData(mergedData)
  }, [networkGraphState, networkGraphNodes])

  const columns = [
    {
      title: 'Oracle',
      dataIndex: 'oracle.name',
      key: 'name',
      sorter: (a, b) =>
        a.oracle.name.localeCompare(b && b.oracle && b.oracle.name),
    },
    {
      title: 'Answer',
      dataIndex: 'state.responseFormatted',
      key: 'answer',
      sorter: (a, b) => {
        if (!a.state || !b.state) return

        return a.state.responseFormatted - b.state.responseFormatted
      },
    },
    {
      title: 'Gas Price (Gwei)',
      dataIndex: 'state.meta.gasPrice',
      key: 'gas',
      sorter: (a, b) => {
        if (!a.state || !b.state) return
        return a.state.meta.gasPrice - b.state.meta.gasPrice
      },
      defaultSortOrder: 'descend',
    },
    {
      title: 'Date',
      dataIndex: 'state.meta.timestamp',
      key: 'timestamp',
      render: timestamp => humanizeUnixTimestamp(timestamp),
    },
  ]

  return (
    <div className="oracle-table">
      <h2 className="oracle-table-header">Oracles data</h2>
      <div className="gas-price-info">
        <h4>
          <Icon type="check-circle" /> Recommended gas price:{' '}
          <b>{gasPrice} Gwei</b>
        </h4>
      </div>
      <Table
        dataSource={data}
        columns={columns}
        pagination={false}
        size={'middle'}
        locale={{ emptyText: <Icon type="loading" /> }}
      />
    </div>
  )
}

export default OracleTable
