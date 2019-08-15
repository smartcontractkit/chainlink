import * as models from 'core/store/models'

export default (data: models.BridgeTypeRequest) => {
  const normalizedData = Object.assign({}, data)
  if (typeof normalizedData.minimumContractPayment === 'number') {
    normalizedData.minimumContractPayment = normalizedData.minimumContractPayment.toString()
  }

  return normalizedData
}
