/**
 * @packageDocumentation
 *
 * An extension to ether's bignumber library that manually
 * polyfills any methods we need for tests by converting the
 * numbers back and forth between ethers.utils.BigNumber and
 * bn.js. If we end up having to replace a ton of methods in the
 * future this way, it might be worth creating a proxy object
 * that automatically does these method polyfills for us.
 */
import { ethers } from 'ethers';
declare module 'ethers' {
    namespace ethers {
        namespace utils {
            interface BigNumber {
                isEven(): boolean;
                umod(val: ethers.utils.BigNumber): ethers.utils.BigNumber;
                shrn(val: number): ethers.utils.BigNumber;
                invm(val: ethers.utils.BigNumber): ethers.utils.BigNumber;
            }
        }
    }
}
export declare function extend(bignumber: typeof ethers.utils.BigNumber): void;
//# sourceMappingURL=BigNumber.d.ts.map