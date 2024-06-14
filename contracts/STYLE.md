# Solidity Style Guide

## Background

Our starting point is the [official Solidity Style Guide](https://solidity.readthedocs.io/en/v0.8.0/style-guide.html) and [ConsenSys's Secure Development practices](https://consensys.github.io/smart-contract-best-practices/), but we deviate in some ways. We lean heavily on [Prettier](https://github.com/smartcontractkit/chainlink/blob/develop/contracts/.prettierrc) for formatting, and if you have to set up a new Solidity project we recommend starting with [our prettier config](https://github.com/smartcontractkit/chainlink/blob/develop/.prettierrc.js). We are trying to automate as much of this styleguide with Solhint as possible.

### Code Organization

- Group functionality together. E.g. Declare structs, events, and helper functions near the functions that use them. This is helpful when reading code because the related pieces are localized. It is also consistent with inheritance and libraries, which are separate pieces of code designed for a specific goal.
    - Why not follow the Solidity recommendation of grouping by visibility? Visibility is clearly defined next to the method signature, making it trivial to check. However, searching can be deceiving because of inherited methods. Given this inconsistency in grouping, we find it easier to read and more consistent to organize code around functionality. Additionally, we recommend testing the public interface for any Solidity contract to ensure it only exposes expected methods.

### Delineate Unaudited Code

- In a large repo it is worthwhile to keep code that has not yet been audited separate from the code that has been audited. This allows you to easily keep track of which files need to be reviewed.
    - E.g. we keep unaudited code in a directory named `dev`. Only once it has been audited we move the audited files out of `dev` and only then is it considered safe to deploy.

## Variables

### Visibility

- All contract variables should be private. Getters should be explicitly written and documented when you want to expose a variable publicly. Whether a getter function reads from storage, a constant, or calculates a value from somewhere else, that’s all implementation details that should not be exposed to the consumer by casing or other conventions.

Examples:

Good:

```javascript
uint256 private s_myVar;

function getMyVar() external view returns(uint256){
    return s_myVar;
}
```

Bad:

```javascript
uint256 public s_myVar;
```

### Naming and Casing

- Function arguments are named like this: `argumentName`. No leading or trailing underscores necessary.
- Storage variables prefixed with an `s_` to make it clear that they live in storage and are expensive to read and write: `s_variableName`. They should always be private, and you should write explicit getters if you want to expose a storage variable.
- Immutable variables should be prefixed with an `i_` to make it clear that they are immutable. E.g. `i_decimalPlaces`. They should always be private, and you should write explicit getters if you want to expose an immutable variable.
- Internal/private constants should be all caps with underscores: `FOO_BAR`. Like other contract variables, constants should not be public. Create getter methods if you want to publicly expose constants.
- Explicitly declare variable size: `uint256` not just `uint`. In addition to being explicit, it matches the naming used to calculate function selectors.

Examples:

Good:

```javascript
uint256 private s_myVar;
uint256 private immutable i_myImmutVar;
uint256 private constant MY_CONST_VAR;

function multiplyMyVar(uint256 multiplier) external view returns(uint256){
    return multiplier * s_myVar;
}
```

Bad:

```javascript
uint private s_myVar;
uint256 private immutable myImmutVar;
uint256 private constant s_myConstVar;

function multiplyMyVar_(uint _multiplier) external view returns(uint256){
    return _mutliplier * s_myVar;
}
```

### Types

- If you are storing an address and know/expect it to be of a type(or interface), make the variable that type. This more clearly documents the behavior of this variable than the `address` type and often leads to less casting code whenever the address is used.

Examples:

Good:

```javascript
import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
// .
// .
// .
AggregatorV3Interface private s_priceFeed;

constructor(address priceFeed) {
    s_priceFeed = AggregatorV3Interface(priceFeed);
}
```

Bad:

```javascript
import "@chainlink/contracts/src/v0.8/interfaces/AggregatorV3Interface.sol";
// .
// .
// .
address private s_priceFeed;

constructor(address priceFeed) {
    s_priceFeed = priceFeed;
}
```

## Functions

### Visibility

- Method visibility should always be explicitly declared. Contract’s [public interfaces should be tested](https://github.com/smartcontractkit/chainlink/blob/master/contracts/test/test-helpers/helpers.ts#L221) to enforce this and make sure that internal logic is not accidentally exposed.

### Naming

- Function names should start with imperative verbs, not nouns or other tenses.
    - `requestData` not `dataRequest`
    - `approve` not `approved`
- Prefix private and internal methods with an underscore. There should never be a publicly callable method starting with an underscore.
    - E.g. `_setOwner(address)`
- Prefix your public getters with `get` and your public setters with `set`.
    - `getConfig` and `setConfig`.

## Modifiers

- Only extract a modifier once a check is duplicated in multiple places. Modifiers arguably hurt readability, so we have found that they are not worth extracting until there is duplication.
- Modifiers should be treated as if they are view functions. They should not change state, only read it. While it is possible to change state in a modifier, it is unconventional and surprising.

### Naming

There are two common classes of modifiers, and their name should be prefixed accordingly to quickly represent their behavior:

- Control flow modifiers: Prefix the modifier name with `if` in the case that a modifier only enables or disables the subsequent code in the modified method, but does not revert.
- Reverting modifiers: Prefix the modifier name with `validate` in the case that a modifier reverts if a condition is not met.

### Return Values

- If an address is cast as a contract type, return the type, do not cast back to the address type. This prevents the consumer of the method signature from having to cast again, but presents an equivalent API for off-chain APIs. Additionally it is a more declarative API, providing more context if we return a type.

## Events

- Events should only be triggered on state changes. If the value is set but not changed, we prefer avoiding a log emission indicating a change. (e.g. Either do not emit a log, or name the event `ConfigSet` instead of `ConfigUpdated`.)

### Naming

- When possible event names should correspond to the method they are in or the action that is being taken. Events preferably follow the format <subject><actionPerformed>, where the action performed is the past tense of the imperative verb in the method name.  e.g. calling `setConfig` should emit an event called `ConfigSet`, not `ConfigUpdated` in a method named `setConfig`.

## Errors

### Use Custom Errors

Whenever possible (Solidity v0.8+) use [custom errors](https://blog.soliditylang.org/2021/04/21/custom-errors/) instead of emitting strings. This saves contract code size and simultaneously provides more informative error messages.

### Expose Errors

It is common to call a contract and then check the call succeeded:

```javascript
(bool success, ) = to.call(data);
require(success, "Contract call failed");
```

While this may look descriptive it swallows the error. Instead bubble up the error:

```javascript
error YourError(bytes response);

(bool success, bytes memory response) = to.call(data);
if (!success) { revert YourError(response); }
```

This will cost slightly more gas to copy the response into memory, but will ultimately make contract usage more understandable and easier to debug. Whether it is worth the extra gas is a judgement call you’ll have to make based on your needs.

The original error will not be human readable in an off-chain explorer because it is RLP hex encoded but is easily decoded with standard Solidity ABI decoding tools, or a hex to UTF-8 converter and some basic ABI knowledge.

## Control Flow

### `if` Statements

Always wrap the result statement of your `if` conditions in a closure, even if it is only one line.

Bad:

```javascript
  if (condition) statement;
```

Good:

```javascript
  if (condition) { statement; }
```

## Interfaces

### Scope

- Interfaces should be as concise as reasonably possible. Break it up into smaller composable interfaces when that is sensible.

### Naming

- Up through Solidity version 0.8: Interfaces should be named `FooInterface`, this follows our historical naming pattern.
- Starting in Solidity v0.9: Interfaces should be named `IFoo` instead of `FooInterface`. This follows the patterns of popular [libraries like OpenZeppelin’s](https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/IERC20.sol#L9).

## Vendor Dependencies

- That’s it, vendor your Solidity dependencies. Supply chain attacks are all the rage these days. There is not yet a silver bullet for best way to vendor, it depends on the size of your project and your needs. You should be as explicit as possible about where the code comes from and make sure that this is enforced in some way; e.g. reference a hash. Some options:
    - NPM packages work for repos already in the JavaScript ecosystem. If you go this route you should lock to a hash of the repo or use a proxy registry like GitHub Packages.
    - Git submodules are great if you don’t mind git submodules.
    - Copy and paste the code into a `vendor` directory. Record attribution of the code and license in the repo along with the commit or version that you pulled the code from.

## Common Behaviors

### Transferring Ownership

- When transferring control, whether it is of a token or a role in a contract, prefer "safe ownership" transfer patterns where the recipient must accept ownership. This avoids accidentally burning the control. This is also inline with the secure pattern of [prefer pull over push](https://consensys.github.io/smart-contract-best-practices/recommendations/#favor-pull-over-push-for-external-calls).

### Use Factories

- If you expect multiple instances of a contract to be deployed, it is probably best to [build a factory](https://www.quicknode.com/guides/solidity/how-to-create-a-smart-contract-factory-in-solidity-using-hardhat) as this allows for simpler deployments later. Additionally it reduces the burden of verifying the correctness of the contract deployment. If many people have to deploy an instance of a contract then doing so with a contract makes it much easier for verification because instead of checking the code hash and/or the compiler and maybe the source code, you only have to check that the contract was deployed through a factory.
- Factories can add some pain when deploying with immutable variables. In general it is difficult to parse those out immutable variables from internal transactions. There is nothing inherently wrong with contracts deployed in this manner  but at the time of writing they may not easily verify on Etherscan.

### Call with Exact Gas

- `call` accepts a gas parameter, but that parameter is a ceiling on gas usage. If a transaction does not have enough gas, `call` will simply provide as much gas as it safely can. This is unintuitive and can lead to transactions failing for unexpected reasons. We have [an implementation of `callWithExactGas`](https://github.com/smartcontractkit/chainlink/blob/075f3e2caf61b8685d2dc78714f1ee39764fda17/contracts/src/v0.8/KeeperRegistry.sol#L792) to ensure the precise gas amount requested is provided.

## Picking a Pragma

- If a contract or library is expected to be imported by outside parties then the pragma should be kept as loose as possible without sacrificing safety. We publish versions for every minor semver version of Solidity, and maintain a corresponding set of tests for each published version.
    - Examples: libraries, interfaces, abstract contracts, and contracts expected to be inherited from
- Otherwise, Solidity contracts should have a pragma which is locked to a specific version.
    - Example: Most concrete contracts.
- Avoid changing pragmas after audit. Unless there is a bug that has affects your contract, then you should try to stick to a known good pragma. In practice this means we typically only support one (occasionally two) pragma for any “major”(minor by semver naming) Solidity version.
