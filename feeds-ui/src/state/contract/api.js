import { HeartbeatContract, provider } from './contract'
import { ethers } from 'ethers'
import { formatEthPrice } from './utils'

export async function oracleAddresses() {
  const addresses = []

  for (let i = 0; i <= 15; i++) {
    try {
      const address = await HeartbeatContract.oracles(i)
      addresses.push(address)
    } catch (err) {
      break
    }
  }
  return addresses
}

export async function currentAnswer() {
  const currentAnswer = await HeartbeatContract.currentAnswer()
  return formatEthPrice(currentAnswer)
}

export async function updateHeight() {
  const updatedHeight = await HeartbeatContract.updatedHeight()
  const block = await provider.getBlock(updatedHeight.toNumber())
  return {
    block: updatedHeight.toNumber(),
    timestamp: block.timestamp
  }
}

export async function nextAnswerId() {
  const answerCounter = await provider.getStorageAt(
    '0x79fEbF6B9F76853EDBcBc913e6aAE8232cFB9De9',
    13
  )
  let bigNumberify = ethers.utils.bigNumberify(answerCounter)
  return bigNumberify.toNumber()
}

export async function latestCompletedAnswer() {
  const currentAnswer = await HeartbeatContract.latestCompletedAnswer()
  return currentAnswer.toNumber()
}

export async function minimumResponses() {
  const minimumResponses = await HeartbeatContract.minimumResponses()
  return minimumResponses.toNumber()
}
