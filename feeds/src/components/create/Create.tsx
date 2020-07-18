import React from 'react'
import { ethers } from 'ethers'
import { Form, Select, Input, Button, InputNumber } from 'antd'
import { FormComponentProps } from 'antd/lib/form/Form'
import { withRouter, RouteComponentProps } from 'react-router'
import { Networks, stringifyQuery } from 'utils'

const { Option } = Select

const formItemLayout = {
  labelCol: { span: 4, offset: 5 },
  wrapperCol: { span: 6 },
}
const formTailLayout = {
  wrapperCol: { span: 8, offset: 9 },
}

const isAddress = () => (_rule: any, value: string, callback: any) => {
  try {
    ethers.utils.getAddress(value)
    callback()
  } catch (error) {
    return callback('Wrong Contract Address')
  }
}

interface CreateProps extends RouteComponentProps, FormComponentProps {}

const Create: React.FC<CreateProps> = ({ form, history }) => {
  const handleSubmit = () => {
    form.validateFields((err, values) => {
      if (!err) {
        history.push({
          pathname: 'custom',
          search: `?${stringifyQuery(values)}`,
        })
      }
    })
  }

  const { getFieldDecorator } = form

  return (
    <>
      <Form.Item {...formTailLayout}>
        <h2>Create Visualization</h2>
      </Form.Item>

      <Form {...formItemLayout}>
        <Form.Item label="Contract Address">
          {getFieldDecorator('contractAddress', {
            validateFirst: true,
            validateTrigger: 'onBlur',
            rules: [
              { required: true, message: 'Contract address is required!' },
              { validator: isAddress() },
            ],
          })(
            <Input placeholder="0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9" />,
          )}
        </Form.Item>

        <Form.Item label="Name">
          {getFieldDecorator('name', {
            rules: [{ required: true, message: 'Name is required!' }],
          })(<Input placeholder="ETH / USD" />)}
        </Form.Item>

        <Form.Item label="Value Prefix">
          {getFieldDecorator('valuePrefix')(<Input placeholder="$" />)}
        </Form.Item>

        <Form.Item label="Heartbeat (seconds)">
          {getFieldDecorator('heartbeat')(
            <InputNumber placeholder="600" style={{ width: '100%' }} />,
          )}
        </Form.Item>

        <Form.Item label="Network">
          {getFieldDecorator('networkId', {
            rules: [{ required: true }],
            initialValue: Networks.MAINNET,
          })(
            <Select placeholder="Select a Network">
              <Option value={Networks.MAINNET}>Mainnet</Option>
              <Option value={Networks.ROPSTEN}>Ropsten</Option>
            </Select>,
          )}
        </Form.Item>

        <Form.Item label="Answer decimal places">
          {getFieldDecorator('decimalPlaces', {
            rules: [{ required: true }],
            initialValue: 6,
          })(<InputNumber placeholder="6" style={{ width: '100%' }} />)}
        </Form.Item>

        <Form.Item label="Format decimal places">
          {getFieldDecorator('formatDecimalPlaces', {
            rules: [{ required: true }],
            initialValue: 0,
          })(<InputNumber placeholder="6" style={{ width: '100%' }} />)}
        </Form.Item>

        <Form.Item label="Answer multiply">
          {getFieldDecorator('multiply', {
            rules: [{ required: true }],
            initialValue: 100000000,
          })(<Input placeholder="100000000" />)}
        </Form.Item>

        <Form.Item label="History">
          {getFieldDecorator('history', {
            rules: [{ required: true }],
            initialValue: 'true',
          })(
            <Select placeholder="Select">
              <Option value={'true'}>Yes</Option>
              <Option value={'false'}>No</Option>
            </Select>,
          )}
        </Form.Item>

        <Form.Item label="History days">
          {getFieldDecorator('historyDays', {
            rules: [{ required: false }],
            initialValue: 1,
          })(<Input placeholder="1" />)}
        </Form.Item>

        <Form.Item {...formTailLayout}>
          <Button type="primary" onClick={() => handleSubmit()}>
            Create
          </Button>
        </Form.Item>
      </Form>
    </>
  )
}

const WrappedComponent = Form.create({ name: 'create' })(withRouter(Create))

export default WrappedComponent
