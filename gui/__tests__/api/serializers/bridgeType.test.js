import serialize from 'api/serializers/bridgeType'

describe('Serializers - Bridge Type', () => {
  it('transforms minimumContractPayment from a number to a string', () => {
    expect(serialize({})).toEqual({})

    const obj = { minimumContractPayment: 10.2 }

    expect(serialize(obj)).toEqual({
      minimumContractPayment: '10.2'
    })
  })
})
