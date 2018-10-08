export default data => {
  const normalizedData = Object.assign({}, data)
  if (typeof normalizedData.minimumContractPayment === 'number') {
    normalizedData.minimumContractPayment = normalizedData.minimumContractPayment.toString()
  }

  return normalizedData
}
