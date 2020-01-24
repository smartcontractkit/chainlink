import React, { useEffect, useState } from 'react'
import _ from 'lodash'
import { Table, Icon } from 'antd'
import { humanizeUnixTimestamp } from 'utils'
/* import { DispatchBinding } from '@chainlink/ts-helpers' */

interface StateProps {
  networkGraphState: any
  networkGraphNodes: any
  fetchEthGasPrice: any
  ethGasPrice: any
}

interface DispatchProps {
  fetchEthGasPrice: any
  /* fetchEthGasPrice: DispatchBinding<typeof aggregationOperations.fetchEthGasPrice> */
}

export interface Props extends StateProps, DispatchProps {}

const OracleTable: React.FC<Props> = ({
  networkGraphState,
  networkGraphNodes,
  fetchEthGasPrice,
  ethGasPrice,
}) => {
  const [data, setData] = useState<any | undefined>()
  const [gasPrice, setGasPrice] = useState<any | undefined>()

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
      .filter((node: any) => node.type === 'oracle')
      .map((oracle: any) => {
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
      sorter: (a: any, b: any): number =>
        a.oracle.name.localeCompare(b && b.oracle && b.oracle.name),
    },

    {
      title: 'Answer',
      dataIndex: 'state.responseFormatted',
      key: 'answer',
      sorter: (a: any, b: any): number => {
        if (!a.state || !b.state) return 0
        return a.state.responseFormatted - b.state.responseFormatted
      },
    },
    {
      title: 'Gas Price (Gwei)',
      dataIndex: 'state.meta.gasPrice',
      key: 'gas',
      sorter: (a: any, b: any): number => {
        if (!a.state || !b.state) return 0
        return a.state.meta.gasPrice - b.state.meta.gasPrice
      },
      defaultSortOrder: 'descend' as 'descend',
    },
    {
      title: 'Date',
      dataIndex: 'state.meta.timestamp',
      key: 'timestamp',
      sorter: (a: any, b: any): number => {
        if (!a.state || !b.state) return 0
        return a.state.meta.gasPrice - b.state.meta.gasPrice
      },
      render: (timestamp: number) => humanizeUnixTimestamp(timestamp),
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
