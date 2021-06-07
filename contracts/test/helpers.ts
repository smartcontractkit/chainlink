import { BigNumber, Contract } from "ethers";
import type { providers } from "ethers";
import { assert } from "chai";
import { ethers } from "hardhat";

export const constants = {
  ZERO_ADDRESS: "0x0000000000000000000000000000000000000000",
  ZERO_BYTES32: "0x0000000000000000000000000000000000000000000000000000000000000000",
  MAX_UINT256: BigNumber.from("2").pow(BigNumber.from("256")).sub(BigNumber.from("1")),
  MAX_INT256: BigNumber.from("2").pow(BigNumber.from("255")).sub(BigNumber.from("1")),
  MIN_INT256: BigNumber.from("2").pow(BigNumber.from("255")).mul(BigNumber.from("-1")),
};

/**
 * Convert an Ether value to a wei amount
 *
 * @param args Ether value to convert to an Ether amount
 */
export function toWei(...args: Parameters<typeof ethers.utils.parseEther>): ReturnType<typeof ethers.utils.parseEther> {
  return ethers.utils.parseEther(...args);
}

/**
 * Increase the current time within the evm to "n" seconds past the current time
 *
 * @param seconds The number of seconds to increase to the current time by
 * @param provider The ethers provider to send the time increase request to
 */
export async function increaseTimeBy(seconds: number, provider: providers.JsonRpcProvider) {
  await provider.send("evm_increaseTime", [seconds]);
}

/**
 * Instruct the provider to mine an additional block
 *
 * @param provider The ethers provider to instruct to mine an additional block
 */
export async function mineBlock(provider: providers.JsonRpcProvider) {
  await provider.send("evm_mine", []);
}

/**
 * Parse out an evm word (32 bytes) into an address (20 bytes) representation
 *
 * @param hex The evm word in hex string format to parse the address
 * out of.
 */
export function evmWordToAddress(hex?: string): string {
  if (!hex) {
    throw Error("Input not defined");
  }

  assert.equal(hex.slice(0, 26), "0x000000000000000000000000");
  return ethers.utils.getAddress(hex.slice(26));
}

/**
 * Check that a contract's abi exposes the expected interface.
 *
 * @param contract The contract with the actual abi to check the expected exposed methods and getters against.
 * @param expectedPublic The expected public exposed methods and getters to match against the actual abi.
 */
export function publicAbi(contract: Contract, expectedPublic: string[]) {
  const actualPublic = [];
  for (const m in contract.functions) {
    if (!m.includes("(")) {
      actualPublic.push(m);
    }
  }

  for (const method of actualPublic) {
    const index = expectedPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is NOT expected to be public`);
  }

  for (const method of expectedPublic) {
    const index = actualPublic.indexOf(method);
    assert.isAtLeast(index, 0, `#${method} is expected to be public`);
  }
}
