/**
 * @packageDocumentation
 *
 * This file contains functionality for ease of creating ethereum account abstractions
 * based on ethers.js. Useful for creating many accounts for testing purposes only.
 */
import { ethers } from 'ethers';
import { JsonRpcProvider } from 'ethers/providers';
interface RCreateFundedWallet {
    /**
     * The created wallet
     */
    wallet: ethers.Wallet;
    /**
     * The receipt of the tx that funded the created wallet
     */
    receipt: ethers.providers.TransactionReceipt;
}
/**
 * Create a pre-funded wallet with all defaults
 *
 * @param provider The provider to connect to the created wallet and to withdraw funds from
 * @param accountIndex The account index of the corresponding wallet derivation path
 */
export declare function createFundedWallet(provider: JsonRpcProvider, accountIndex: number): Promise<RCreateFundedWallet>;
/**
 * Create an ethers.js wallet instance that is connected to the given provider
 *
 * @param provider A compatible ethers.js provider such as the one returned by `ganache.provider()` to connect the wallet to
 * @param accountIndex The account index to derive from the mnemonic phrase
 */
export declare function createWallet(provider: ethers.providers.JsonRpcProvider, accountIndex: number): ethers.Wallet;
/**
 * Fund a wallet with unlocked accounts available from the given provider
 *
 * @param wallet The ethers wallet to fund
 * @param provider The provider which has control over unlocked, funded accounts to transfer funds from
 * @param overrides Transaction parameters to override when sending the funding transaction
 */
export declare function fundWallet(wallet: ethers.Wallet, provider: ethers.providers.JsonRpcProvider, overrides?: Omit<ethers.providers.TransactionRequest, 'to' | 'from'>): Promise<ethers.providers.TransactionReceipt>;
export {};
//# sourceMappingURL=wallet.d.ts.map