/**
 * @packageDocumentation
 *
 * This file contains functionality for ease of creating ethereum account abstractions
 * based on ethers.js. Useful for creating many accounts for testing purposes only.
 */
import { ethers } from 'ethers'
import { JsonRpcProvider } from 'ethers/providers'
import { makeDebug } from './debug'
import { isOVM } from './helpers'

interface RCreateFundedWallet {
  /**
   * The created wallet
   */
  wallet: ethers.Wallet
  /**
   * The receipt of the tx that funded the created wallet
   */
  receipt?: ethers.providers.TransactionReceipt
}
/**
 * Create a pre-funded wallet with all defaults
 *
 * @param provider The provider to connect to the created wallet and to withdraw funds from
 * @param accountIndex The account index of the corresponding wallet derivation path
 */
export async function createFundedWallet(
  provider: JsonRpcProvider,
  accountIndex: number,
): Promise<RCreateFundedWallet> {
  const wallet = createWallet(provider, accountIndex)
  const receipt = await fundWallet(wallet, provider)
  return { wallet, receipt }
}

/**
 * Create an ethers.js wallet instance that is connected to the given provider
 *
 * @param provider A compatible ethers.js provider such as the one returned by `ganache.provider()` to connect the wallet to
 * @param accountIndex The account index to derive from the mnemonic phrase
 */
export function createWallet(
  provider: ethers.providers.JsonRpcProvider,
  accountIndex: number,
): ethers.Wallet {
  const debug = makeDebug('wallet:createWallet')
  if (accountIndex < 0) {
    throw Error(`Account index must be greater than 0, got ${accountIndex}`)
  }

  /**
   * THIS IS FOR TESTING PURPOSES ONLY
   */
  const mnemonicPhrase =
    'dose weasel clever culture letter volume endorse used harvest ripple circle install'

  const path = `m/44'/60'/${accountIndex}'/0/0`
  debug('created wallet with parameters: %o', { mnemonicPhrase, path })

  return ethers.Wallet.fromMnemonic(mnemonicPhrase, path).connect(provider)
}

/**
 * Fund a wallet with unlocked accounts available from the given provider
 *
 * @param wallet The ethers wallet to fund
 * @param provider The provider which has control over unlocked, funded accounts to transfer funds from
 * @param overrides Transaction parameters to override when sending the funding transaction
 */
export async function fundWallet(
  wallet: ethers.Wallet,
  provider: ethers.providers.JsonRpcProvider,
  overrides?: Omit<ethers.providers.TransactionRequest, 'to' | 'from'>,
): Promise<ethers.providers.TransactionReceipt | undefined> {
  const debug = makeDebug('wallet:fundWallet')

  // OVM has no native ETH/value, so wallet funding is not supported
  if (isOVM()) return

  debug('funding wallet')

  debug('retreiving accounts...')

  const nodeOwnedAccounts = await provider.listAccounts()
  debug('retreived accounts: %o', nodeOwnedAccounts)

  const signer = provider.getSigner(nodeOwnedAccounts[0])

  const txParams: ethers.providers.TransactionRequest = {
    to: wallet.address,
    value: ethers.utils.parseEther('10'),
    ...overrides,
  }
  debug('sending tx with the following parameters: %o', txParams)
  const tx = await signer.sendTransaction(txParams)

  debug('waiting on tx %s to complete...', tx.hash)
  const receipt = await tx.wait()
  debug('tx %s confirmed with tx receipt %o', tx.hash, receipt)
  return receipt
}
