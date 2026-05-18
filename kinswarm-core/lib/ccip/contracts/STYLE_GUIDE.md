# Structure

This guide is split into two sections: [Guidelines](#guidelines) and [Rules](#rules).
Guidelines are recommendations that should be followed but are hard to enforce in an automated way.
Rules are all enforced through CI, this can be through Solhint rules or other tools.

## Background

Our starting point is the [official Solidity Style Guide](https://docs.soliditylang.org/en/v0.8.21/style-guide.html) and [ConsenSys's Secure Development practices](https://consensys.github.io/smart-contract-best-practices/), but we deviate in some ways. We lean heavily on [Prettier](https://github.com/smartcontractkit/chainlink/blob/develop/contracts/.prettierrc) for formatting, and if you have to set up a new Solidity project we recommend starting with [our prettier config](https://github.com/smartcontractkit/chainlink/blob/develop/contracts/.prettierrc). We are trying to automate as much of this style guide with Solhint as possible.

This guide is not meant to be applied retroactively. There is no need to rewrite existing code to adhere to this guide, and when making (small) changes in existing files, it is not required to adhere to this guide if it conflicts with other practices in that existing file. Consistency is preferred.

We will be looking into `forge fmt`, but for now, we still use `prettier`.


# <a name="guidelines"></a>Guidelines

## Code Organization
- Group functionality together. E.g. Declare structs, events, and helper functions near the functions that use them. This is helpful when reading code because the related pieces are localized. It is also consistent with inheritance and libraries, which are separate pieces of code designed for a specific goal.
- 🤔Why not follow the Solidity recommendation of grouping by visibility? Visibility is clearly defined next to the method signature, making it trivial to check. However, searching can be deceiving because of inherited methods. Given this inconsistency in grouping, we find it easier to read and more consistent to organize code around functionality. Additionally, we recommend testing the public interface for any Solidity contract to ensure it only exposes expected methods.
- Follow the [Solidity folder structure CLIP](https://github.com/smartcontractkit/CLIPs/tree/main/clips/2023-04-13-solidity-folder-structure)

### Delineate Unaudited Code

- In a large repo, it is worthwhile to keep code that has not yet been audited separately from the code that has been audited. This allows you to easily keep track of which files need to be reviewed.
  - E.g. we keep unaudited code in a directory named `dev` that exists within each project's folder. Only once it has been audited we move the audited files out of `dev` and only then is it considered safe to deploy.
  - This `dev` folder also has implications for when code is valid for bug bounties, so be extra careful to move functionality out of a `dev` folder.


## Comments
- Besides comments above functions and structs, comments should live everywhere a reader might be confused.
  Don’t overestimate the reader of your contract, expect confusion in many places and document accordingly.
  This will help massively during audits and onboarding new team members.
- Headers should be used to group functionality, the following header style and length are recommended.
  - Don’t use headers for a single function, or to say “getters”. Group by functionality e.g. the `Tokens and pools`, or `fees` logic within the CCIP OnRamp.
- Comments should start with a capital letter and end with a period.

```solidity
  // ================================================================
  // │                      Tokens and pools                        │
  // ================================================================

....

  // ================================================================
  // │                             Fees                             │
  // ================================================================
```

## Variables

- Function arguments are named like this: `argumentName`. No leading or trailing underscores are necessary.
- Names should be explicit on the unit it contains, e.g. a network fee that is charged in USD cents

```solidity
uint256 fee; // bad
uint256 networkFee; // bad
uint256 networkFeeUSD; // bad
uint256 networkFeeUSDCents; // good
```

### Types

- If you are storing an address and know/expect it to be of a type(or interface), make the variable that type. This more clearly documents the behavior of this variable than the `address` type and often leads to less casting code whenever the address is used.

### Structs

- Structs should contain struct packing comments to indicate the storage slot layout
  - Using the exact characters from the example below will ensure visually appealing struct packing comments.
  - Notice there is no line on the unpacked last `fee` item.
- Struct should contain comments, clearly indicating the denomination of values e.g. 0.01 USD if the variable name doesn’t already do that (which it should).
  - A simple tool that could help with packing structs and adding comments: https://github.com/RensR/Spack

```solidity
/// @dev Struct to hold the fee configuration for a fee token, same as the FeeTokenConfig but with
/// token included so that an array of these can be passed in to setFeeTokenConfig to set the mapping
struct FeeTokenConfigArgs {
  address token; // ────────────╮ Token address
  uint32 networkFeeUSD; //      │ Flat network fee to charge for messages, multiples of 0.01 USD
  //                            │ multiline comments should work like this. More fee info
  uint64 gasMultiplier; // ─────╯ Price multiplier for gas costs, 1e18 based so 11e17 = 10% extra cost
  uint64 premiumMultiplier; // ─╮ Multiplier for fee-token-specific premiums
  bool enabled; // ─────────────╯ Whether this fee token is enabled
  uint256 fee; //                 The flat fee the user pays in juels
}
```
## Functions

### Naming

- Function names should start with imperative verbs, not nouns or other tenses.
  - `requestData` not `dataRequest`
  - `approve` not `approved`
  - `getFeeParameters` not `feeParameters`

- Prefix your public getters with `get` and your public setters with `set`.
  - `getConfig` and `setConfig`.

### Return Values

- If an address is cast as a contract type, return the type, do not cast back to the address type.
  This prevents the consumer of the method signature from having to cast again but presents an equivalent API for off-chain APIs.
  Additionally, it is a more declarative API, providing more context if we return a type.

## Modifiers

- Only extract a modifier once a check is duplicated in multiple places. Modifiers arguably hurt readability, so we have found that they are not worth extracting until there is duplication.
- Modifiers should be treated as if they are view functions. They should not change state, only read it. While it is possible to change the state in a modifier, it is unconventional and surprising.
- Modifiers tend to bloat contract size because the code is duplicated wherever the modifier is used.

## Events

- Events should only be triggered on state changes. If the value is set but not changed, we prefer avoiding a log emission indicating a change. (e.g. Either do not emit a log, or name the event `ConfigSet` instead of `ConfigUpdated`.)
- Events should be emitted for all state changes, not emitting should be an exception
- When possible event names should correspond to the method they are in or the action that is being taken. Events preferably follow the format <subject><actionPerformed>, where the action performed is the past tense of the imperative verb in the method name.  e.g. calling `setConfig` should emit an event called `ConfigSet`, not `ConfigUpdated` in a method named `setConfig`.


### Expose Errors

It is common to call a contract and then check the call succeeded:

```solidity
(bool success, ) = to.call(data);
require(success, "Contract call failed");
```

While this may look descriptive it swallows the error. Instead, bubble up the error:

```solidity
bool success;
retData = new bytes(maxReturnBytes);
assembly {
  // call and return whether we succeeded. ignore return data
  // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
  success := call(gasLimit, target, 0, add(payload, 0x20), mload(payload), 0, 0)

  // limit our copy to maxReturnBytes bytes
  let toCopy := returndatasize()
  if gt(toCopy, maxReturnBytes) {
    toCopy := maxReturnBytes
  }
  // Store the length of the copied bytes
  mstore(retData, toCopy)
  // copy the bytes from retData[0:_toCopy]
  returndatacopy(add(retData, 0x20), 0, toCopy)
}
return (success, retData);
```

This will cost slightly more gas to copy the response into memory, but will ultimately make contract usage more understandable and easier to debug. Whether it is worth the extra gas is a judgment call you’ll have to make based on your needs.

The original error will not be human-readable in an off-chain explorer because it is RLP hex encoded but is easily decoded with standard Solidity ABI decoding tools, or a hex to UTF-8 converter and some basic ABI knowledge.


## Interfaces

- Interfaces should be as concise as reasonably possible. Break it up into smaller composable interfaces when that is sensible.

## Dependencies

- Prefer not reinventing the wheel, especially if there is an Openzeppelin wheel.
- The `shared` folder can be treated as a first-party dependency and it is recommended to check if some functionality might already be in there before either writing it yourself or adding a third party dependency.
- When we have reinvented the wheel already (like with ownership), it is OK to keep using these contracts. If there are clear benefits of using another standard like OZ, we can deprecate the custom implementation and start using the new standard in all new projects. Migration will not be required unless there are serious issues with the old implementation.
- When the decision is made to use a new standard, it is no longer allowed to use the old standard for new projects.

### Vendor dependencies

- That’s it, vendor your Solidity dependencies. Supply chain attacks are all the rage these days. There is not yet a silver bullet for the best way to vendor, it depends on the size of your project and your needs. You should be as explicit as possible about where the code comes from and make sure that this is enforced in some way; e.g. reference a hash. Some options:
  - NPM packages work for repos already in the JavaScript ecosystem. If you go this route you should lock to a hash of the repo or use a proxy registry like GitHub Packages.
  - Copy and paste the code into a `vendor` directory. Record attribution of the code and license in the repo along with the commit or version that you pulled the code from.
  - Foundry uses git submodules for its dependencies. We only use the `forge-std` lib through submodules, we don’t import any non-Foundry-testing code through this method.


## Common Behaviors

### Transferring Ownership

- When transferring control, whether it is of a token or a role in a contract, prefer "safe ownership" transfer patterns where the recipient must accept ownership. This avoids accidentally burning the control. This is also in line with the secure pattern of [prefer pull over push](https://consensys.github.io/smart-contract-best-practices/recommendations/#favor-pull-over-push-for-external-calls).

### Call with Exact Gas

- `call` accepts a gas parameter, but that parameter is a ceiling on gas usage. If a transaction does not have enough gas, `call` will simply provide as much gas as it safely can. This is unintuitive and can lead to transactions failing for unexpected reasons. We have [an implementation of `callWithExactGas`](https://github.com/smartcontractkit/chainlink/blob/075f3e2caf61b8685d2dc78714f1ee39764fda17/contracts/src/v0.8/KeeperRegistry.sol#L792) to ensure the precise gas amount requested is provided.

### Sending tokens

- Prefer [ERC20.safeTransfer](https://docs.openzeppelin.com/contracts/2.x/api/token/erc20#SafeERC20) over ERC20.transfer

### Gas golfing

- Golf your code. Make it cheap, within reason.
  - Focus on the hot path
- Most of the cost of executing Solidity is related to reading/writing storage
- Calling other contracts will also be costly
- Common types to safely use are
  - uint40 for timestamps (or uint32 if you really need the space)
  - uint96 for LINK, as there are only 1b LINK tokens
- prefer `++i` over `i++`
- If you’re unsure about golfing, ask in the #tech-solidity channel

## Testing

- Test using Foundry.
- Aim for at least 90% *useful* coverage as a baseline, but (near) 100% is achievable in Solidity. Always 100% test the critical path.
  - Make sure to test for each event emitted
  - Test each reverting path
- Consider fuzzing, start with stateless (very easy in Foundry) and if that works, try stateful fuzzing.
- Consider fork testing if applicable

### Foundry

- Create a Foundry profile for each project folder in `foundry.toml`
- Foundry tests live in the project folder in `src`, not in the `contracts/test/` folder
- Set the block number and timestamp. It is preferred to set these values to some reasonable value close to reality.
- There should be no code between `vm.expectEmit`/`vm.expectRevert` and the function call

## Picking a Pragma

- If a contract or library is expected to be imported by outside parties then the pragma should be kept as loose as possible without sacrificing safety. We publish versions for every minor Semver version of Solidity and maintain a corresponding set of tests for each published version.
  - Examples: libraries, interfaces, abstract contracts, and contracts expected to be inherited from
- Otherwise, Solidity contracts should have a pragma that is locked to a specific version.
  - Example: Most concrete contracts.
- Avoid changing pragmas after the audit. Unless there is a bug that affects your contract, then you should try to stick to a known good pragma. In practice, this means we typically only support one (occasionally two) pragma for any “major”(minor by Semver naming) Solidity version.
- The current advised pragma is `0.8.24`, lower versions should be avoided when starting a new project. Newer versions can be considered.
  - Explicitly use the `Paris` hardfork when compiling with 0.8.24 to keep the bytecode compatible with all chains.
- All contracts should have an SPDX license identifier. If unsure about which one to pick, please consult with legal. Most older contracts have been MIT, but some of the newer products have been using BUSL-1.1


## Versioning

Contracts should implement the following interface

```solidity
interface ITypeAndVersion {
  function typeAndVersion() external pure returns (string memory);
}
```

Here are some examples of what this should look like:

```solidity
contract AccessControlledFoo is Foo {
  string public constant override typeAndVersion = "AccessControlledFoo 1.0.0";
}

contract OffchainAggregator is ITypeAndVersion {
   string public constant override typeAndVersion = "OffchainAggregator 1.0.0";

    function getData() public returns(uint256) {
        return 4;
    }
}

// Next version of Aggregator contract
contract SuperDuperAggregator is ITypeAndVersion {
    /// This is a new contract that has not been released yet, so we
    /// add a `-dev` suffix to the typeAndVersion.

    // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
    string public constant override typeAndVersion = "SuperDuperAggregator 1.1.0-dev";

    function getData() public returns(uint256) {
      return 5;
    }
}
```

All contracts will expose a `typeAndVersion` constant.
The string has the following format: `<contract name><SPACE><semver>-<dev>` with the `-dev` part only being applicable to contracts that have not been fully released.
Try to fit it into 32 bytes to keep the impact on contract sizes minimal.

Note that `ITypeAndVersion` should be used, not `TypeAndVersionInterface`.









# <a name="rules"></a>Rules

All rules have a `rule` tag which indicates how the rule is enforced.


## Comments

- Comments should be in the `//` (default) or `///` (Natspec) format, not the `/*  */` format.
  - rule: `tbd`
- Comments should follow [NatSpec](https://docs.soliditylang.org/en/latest/natspec-format.html)
  - rule: `tbd`

## Imports

- Imports should always be explicit
  - rule: `no-global-import`
- Imports have to follow the following format:
  - rule: `tbd`

```solidity
import {IInterface} from "../interfaces/IInterface.sol";

import {AnythingElse} from "../code/AnythingElse.sol";

import {ThirdPartyCode} from "../../vendor/ThirdPartyCode.sol";
```

- An example would be

```solidity
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IPool} from "../interfaces/pools/IPool.sol";

import {AggregateRateLimiter} from "../AggregateRateLimiter.sol";
import {Client} from "../libraries/Client.sol";

import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
```

## Variables

### Visibility

All contract variables should be private by default. Getters should be explicitly written and documented when you want to expose a variable publicly.
Whether a getter function reads from storage, a constant, or calculates a value from somewhere else, that’s all implementation details that should not be exposed to the consumer by casing or other conventions.

rule: `tbd`

### Naming and Casing

- Storage variables prefixed with an `s_` to make it clear that they live in storage and are expensive to read and write: `s_variableName`. They should always be private, and you should write explicit getters if you want to expose a storage variable.
  - rule: `chainlink-solidity/prefix-storage-variables-with-s-underscore`
- Immutable variables should be prefixed with an `i_` to make it clear that they are immutable. E.g. `i_decimalPlaces`. They should always be private, and you should write explicit getters if you want to expose an immutable variable.
  - rule: `chainlink-solidity/prefix-immutable-variables-with-i`
- Internal/private constants should be all caps with underscores: `FOO_BAR`. Like other contract variables, constants should not be public. Create getter methods if you want to publicly expose constants.
  - rule: `chainlink-solidity/all-caps-constant-storage-variables`
- Explicitly declare variable size: `uint256` not just `uint`. In addition to being explicit, it matches the naming used to calculate function selectors.
  - rule: `explicit-types`
- Mapping should always be named if Solidity allows it (≥0.8.18)
  - rule: `tbd`


## Functions

### Visibility

- Method visibility should always be explicitly declared.
  - rule: `state-visibility`

- Prefix private and internal methods with an underscore. There should never be a publicly callable method starting with an underscore.
  - E.g. `_setOwner(address)`
  - rule: `chainlink-solidity/prefix-internal-functions-with-underscore`

### Return values

- Returned values should always be explicit. Using named return values and then returning with an empty return should be avoided
  - rule: `chainlink-solidity/explicit-returns`

```solidity
// Bad
function getNum() external view returns (uint64 num) {
  num = 4;
  return;
}

// Good
function getNum() external view returns (uint64 num) {
  num = 4;
  return num;
}

// Good
function getNum() external view returns (uint64 num) {
  return 4;
}
```

## Errors

Use [custom errors](https://blog.soliditylang.org/2021/04/21/custom-errors/) instead of emitting strings. This saves contract code size and simultaneously provides more informative error messages.

rule: `gas-custom-errors`

## Interfaces

### Purpose

Interfaces separate NatSpec from contract logic, requiring readers to do more work to understand the code. For this reason, you shouldn’t create an interface by default.

If created, interfaces should have a documented purpose. For example, an interface is useful if 3rd party on-chain contracts will interact with your contract. CCIP’s [`IRouterClient` interface](https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/interfaces/IRouterClient.sol) is a good example here.

### Naming

Interfaces should be named `IFoo` instead of `FooInterface`. This follows the patterns of popular [libraries like OpenZeppelin’s](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/IERC20.sol#L9). 

rule: `interface-starts-with-i`

## Structs

- All structs should be packed to have the lowest memory footprint to reduce gas usage. Even structs that will never be written to storage should be packed.
  - A contract can be considered a struct; it should also be packed to reduce gas cost.

rule: `gas-struct-packing`

- Structs should be constructed with named arguments. This prevents accidental assignment to the wrong field and makes the code more readable.

```solidity
// Good
function setConfig(uint64 _foo, uint64 _bar, uint64 _baz) external {
  config = Config({
    foo: _foo,
    bar: _bar,
    baz: _baz
  });
}

// Bad
function setConfig(uint64 _foo, uint64 _bar, uint64 _baz) external {
  config = Config(_foo, _bar, _baz);
}
```

rule: `tbd`
