import React from 'react'
import { ethers } from 'ethers'
import { Form, Select, Input, Button, InputNumber } from 'antd'
import { withRouter } from 'react-router'
import { ROPSTEN_ID, MAINNET_ID } from 'utils'
import { stringifyQuery } from 'utils'

const { Option } = Select

const formItemLayout = {
  labelCol: { span: 4, offset: 5 },
  wrapperCol: { span: 6 },
}
const formTailLayout = {
  wrapperCol: { span: 8, offset: 9 },
}

export const isAddress = () => (rule, value, callback) => {
  try {
    ethers.utils.getAddress(value)
    callback()
  } catch (error) {
    return callback('Wrong Contract Address')
  }
}

const Create = ({ form, history }) => {
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
        <h2>Create Visualisation</h2>
      </Form.Item>

      <Form {...formItemLayout}>
        <Form.Item label="Contract Address">
          {getFieldDecorator('contractAddress', {
            validateFirst: true,
            validateTrigger: 'onBlur',
            rules: [
              { required: true, message: 'Contract address is required!' },
              { validator: isAddress(form) },
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

        <Form.Item label="Counter (seconds)">
          {getFieldDecorator('counter')(
            <InputNumber placeholder="600" style={{ width: '100%' }} />,
          )}
        </Form.Item>

        <Form.Item label="Network">
          {getFieldDecorator('networkId', {
            rules: [{ required: true }],
            initialValue: MAINNET_ID,
          })(
            <Select placeholder="Select a Network">
              <Option value={MAINNET_ID}>Mainnet</Option>
              <Option value={ROPSTEN_ID}>Ropsten</Option>
            </Select>,
          )}
        </Form.Item>

        <Form.Item label="Answer decimal places">
          {getFieldDecorator('decimalPlaces', {
            rules: [{ required: true }],
            initialValue: 6,
          })(<InputNumber placeholder="6" style={{ width: '100%' }} />)}
        </Form.Item>

        <Form.Item label="Answer multiply">
          {getFieldDecorator('multiply', {
            rules: [{ required: true }],
            initialValue: 1000000000000000000,
          })(<Input placeholder="1000000000000000000" />)}
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
