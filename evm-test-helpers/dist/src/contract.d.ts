/**
 * @packageDocumentation
 *
 * This file deals with contract helpers to deal with ethers.js contract abstractions
 */
import { ethers, Signer, ContractTransaction } from 'ethers';
import { Provider } from 'ethers/providers';
export * from './generated/factories/LinkToken__factory';
/**
 * The type of any function that is deployable
 */
declare type Deployable = {
    deploy: (...deployArgs: any[]) => Promise<any>;
};
/**
 * Get the return type of a function, and unbox any promises
 */
export declare type Instance<T extends Deployable> = T extends {
    deploy: (...deployArgs: any[]) => Promise<infer U>;
} ? U : never;
declare type Override<T> = {
    [K in keyof T]: T[K] extends (...args: any[]) => Promise<ContractTransaction> ? (...args: any[]) => Promise<any> : T[K];
};
export declare type CallableOverrideInstance<T extends Deployable> = T extends {
    deploy: (...deployArgs: any[]) => Promise<infer ContractInterface>;
} ? Omit<Override<ContractInterface>, 'connect'> & {
    connect(signer: string | Signer | Provider): CallableOverrideInstance<T>;
} : never;
export declare function callable(oldContract: ethers.Contract, methods: string[]): any;
//# sourceMappingURL=contract.d.ts.map