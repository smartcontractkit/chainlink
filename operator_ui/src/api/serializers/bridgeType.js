export default data => {
  const normalizedData = {
    ...data
  }
  if (typeof normalizedData.minimumContractPayment === 'number') {
    normalizedData.minimumContractPayment = normalizedData.minimumContractPayment.toString()
  }

  return normalizedData
}
