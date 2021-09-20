/**
 * @packageDocumentation
 *
 * This file provides utility functions related to test setup, such as creating a test provider,
 * optimizing test times via snapshots, and making test accounts.
 */
import { ethers } from 'ethers';
/**
 * Create a test provider which uses an in-memory, in-process chain
 */
export declare function provider(): ethers.providers.JsonRpcProvider;
/**
 * This helper function allows us to make use of ganache snapshots,
 * which allows us to snapshot one state instance and revert back to it.
 *
 * This is used to memoize expensive setup calls typically found in beforeEach hooks when we
 * need to setup our state with contract deployments before running assertions.
 *
 * @param provider The provider that's used within the tests
 * @param cb The callback to execute that generates the state we want to snapshot
 */
export declare function snapshot(provider: ethers.providers.JsonRpcProvider, cb: () => Promise<void>): () => Promise<void>;
export interface Roles {
    defaultAccount: ethers.Wallet;
    oracleNode: ethers.Wallet;
    oracleNode1: ethers.Wallet;
    oracleNode2: ethers.Wallet;
    oracleNode3: ethers.Wallet;
    oracleNode4: ethers.Wallet;
    stranger: ethers.Wallet;
    consumer: ethers.Wallet;
}
export interface Personas {
    Default: ethers.Wallet;
    Carol: ethers.Wallet;
    Eddy: ethers.Wallet;
    Nancy: ethers.Wallet;
    Ned: ethers.Wallet;
    Neil: ethers.Wallet;
    Nelly: ethers.Wallet;
    Norbert: ethers.Wallet;
}
interface Users {
    roles: Roles;
    personas: Personas;
}
/**
 * Generate roles and personas for tests along with their correlated account addresses
 */
export declare function users(provider: ethers.providers.JsonRpcProvider): Promise<Users>;
export {};
//# sourceMappingURL=setup.d.ts.map