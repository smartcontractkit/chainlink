**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [tx-origin](#tx-origin) (1 results) (Medium)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
## tx-origin
Impact: Medium
Confidence: Medium
 - [ ] ID-0
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) uses tx.origin for authorization: [tx.origin != address(0) && tx.origin != address(0x1111111111111111111111111111111111111111)](src/v0.8/newer_automation/AutomationBase.sol#L13)

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-1
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) is never used and should be removed

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/AutomationBase.sol#L2)

src/v0.8/newer_automation/AutomationBase.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [tx-origin](#tx-origin) (1 results) (Medium)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
 - [unimplemented-functions](#unimplemented-functions) (1 results) (Informational)
## tx-origin
Impact: Medium
Confidence: Medium
 - [ ] ID-0
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) uses tx.origin for authorization: [tx.origin != address(0) && tx.origin != address(0x1111111111111111111111111111111111111111)](src/v0.8/newer_automation/AutomationBase.sol#L13)

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-1
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) is never used and should be removed

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/AutomationBase.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/AutomationCompatible.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)

src/v0.8/newer_automation/AutomationBase.sol#L2


## unimplemented-functions
Impact: Informational
Confidence: High
 - [ ] ID-3
[AutomationCompatible](src/v0.8/newer_automation/AutomationCompatible.sol#L7) does not implement functions:
	- [AutomationCompatibleInterface.checkUpkeep(bytes)](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L23)
	- [AutomationCompatibleInterface.performUpkeep(bytes)](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L41)

src/v0.8/newer_automation/AutomationCompatible.sol#L7


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (3 results) (Informational)
 - [naming-convention](#naming-convention) (5 results) (Informational)
## pragma
Impact: Informational
Confidence: High
 - [ ] ID-0
3 different versions of Solidity are used:
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L2)
	- Version constraint ^0.8.4 is used by:
		-[^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationV21PlusCommon.sol#L2)
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/newer_automation/interfaces/ILogAutomation.sol#L2)

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-1
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/ILogAutomation.sol#L2)

src/v0.8/newer_automation/interfaces/ILogAutomation.sol#L2


 - [ ] ID-2
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationV21PlusCommon.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationV21PlusCommon.sol#L2


 - [ ] ID-3
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L2)

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L2


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-4
Function [AutomationCompatibleUtils._log(Log)](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L16) is not in mixedCase

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L16


 - [ ] ID-5
Function [AutomationCompatibleUtils._logTrigger(IAutomationV21PlusCommon.LogTrigger)](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L12) is not in mixedCase

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L12


 - [ ] ID-6
Function [AutomationCompatibleUtils._logTriggerConfig(IAutomationV21PlusCommon.LogTriggerConfig)](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L10) is not in mixedCase

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L10


 - [ ] ID-7
Function [AutomationCompatibleUtils._report(IAutomationV21PlusCommon.Report)](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L8) is not in mixedCase

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L8


 - [ ] ID-8
Function [AutomationCompatibleUtils._conditionalTrigger(IAutomationV21PlusCommon.ConditionalTrigger)](src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L14) is not in mixedCase

src/v0.8/newer_automation/AutomationCompatibleUtils.sol#L14


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [missing-zero-check](#missing-zero-check) (2 results) (Low)
 - [assembly](#assembly) (2 results) (Informational)
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (2 results) (Informational)
 - [immutable-states](#immutable-states) (1 results) (Optimization)
## missing-zero-check
Impact: Low
Confidence: Medium
 - [ ] ID-0
[AutomationForwarder.constructor(address,address,address).logic](src/v0.8/newer_automation/AutomationForwarder.sol#L23) lacks a zero-check on :
		- [i_logic = logic](src/v0.8/newer_automation/AutomationForwarder.sol#L26)

src/v0.8/newer_automation/AutomationForwarder.sol#L23


 - [ ] ID-1
[AutomationForwarder.constructor(address,address,address).target](src/v0.8/newer_automation/AutomationForwarder.sol#L23) lacks a zero-check on :
		- [i_target = target](src/v0.8/newer_automation/AutomationForwarder.sol#L25)

src/v0.8/newer_automation/AutomationForwarder.sol#L23


## assembly
Impact: Informational
Confidence: High
 - [ ] ID-2
[AutomationForwarder.fallback()](src/v0.8/newer_automation/AutomationForwarder.sol#L66-L91) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/AutomationForwarder.sol#L70-L90)

src/v0.8/newer_automation/AutomationForwarder.sol#L66-L91


 - [ ] ID-3
[AutomationForwarder.forward(uint256,bytes)](src/v0.8/newer_automation/AutomationForwarder.sol#L35-L60) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/AutomationForwarder.sol#L39-L57)

src/v0.8/newer_automation/AutomationForwarder.sol#L35-L60


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-4
2 different versions of Solidity are used:
	- Version constraint ^0.8.16 is used by:
		-[^0.8.16](src/v0.8/newer_automation/AutomationForwarder.sol#L2)
	- Version constraint ^0.8.4 is used by:
		-[^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/AutomationForwarder.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-5
Version constraint ^0.8.16 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- StorageWriteRemovalBeforeConditionalTermination.
It is used by:
	- [^0.8.16](src/v0.8/newer_automation/AutomationForwarder.sol#L2)

src/v0.8/newer_automation/AutomationForwarder.sol#L2


 - [ ] ID-6
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2


## immutable-states
Impact: Optimization
Confidence: High
 - [ ] ID-7
[AutomationForwarder.s_registry](src/v0.8/newer_automation/AutomationForwarder.sol#L21) should be immutable 

src/v0.8/newer_automation/AutomationForwarder.sol#L21


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (3 results) (Informational)
## pragma
Impact: Informational
Confidence: High
 - [ ] ID-0
3 different versions of Solidity are used:
	- Version constraint ^0.8.16 is used by:
		-[^0.8.16](src/v0.8/newer_automation/AutomationForwarderLogic.sol#L2)
	- Version constraint ^0.8.4 is used by:
		-[^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/shared/interfaces/ITypeAndVersion.sol#L2)

src/v0.8/newer_automation/AutomationForwarderLogic.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-1
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/shared/interfaces/ITypeAndVersion.sol#L2)

src/v0.8/shared/interfaces/ITypeAndVersion.sol#L2


 - [ ] ID-2
Version constraint ^0.8.16 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- StorageWriteRemovalBeforeConditionalTermination.
It is used by:
	- [^0.8.16](src/v0.8/newer_automation/AutomationForwarderLogic.sol#L2)

src/v0.8/newer_automation/AutomationForwarderLogic.sol#L2


 - [ ] ID-3
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [locked-ether](#locked-ether) (1 results) (Medium)
 - [missing-zero-check](#missing-zero-check) (1 results) (Low)
 - [assembly](#assembly) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
 - [naming-convention](#naming-convention) (1 results) (Informational)
## locked-ether
Impact: Medium
Confidence: High
 - [ ] ID-0
Contract locking ether found:
	Contract [Chainable](src/v0.8/newer_automation/Chainable.sol#L9-L61) has payable functions:
	 - [Chainable.fallback()](src/v0.8/newer_automation/Chainable.sol#L34-L60)
	But does not have a function to withdraw the ether

src/v0.8/newer_automation/Chainable.sol#L9-L61


## missing-zero-check
Impact: Low
Confidence: Medium
 - [ ] ID-1
[Chainable.constructor(address).fallbackAddress](src/v0.8/newer_automation/Chainable.sol#L18) lacks a zero-check on :
		- [i_FALLBACK_ADDRESS = fallbackAddress](src/v0.8/newer_automation/Chainable.sol#L19)

src/v0.8/newer_automation/Chainable.sol#L18


## assembly
Impact: Informational
Confidence: High
 - [ ] ID-2
[Chainable.fallback()](src/v0.8/newer_automation/Chainable.sol#L34-L60) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/Chainable.sol#L38-L59)

src/v0.8/newer_automation/Chainable.sol#L34-L60


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-3
Version constraint ^0.8.16 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- StorageWriteRemovalBeforeConditionalTermination.
It is used by:
	- [^0.8.16](src/v0.8/newer_automation/Chainable.sol#L2)

src/v0.8/newer_automation/Chainable.sol#L2


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-4
Variable [Chainable.i_FALLBACK_ADDRESS](src/v0.8/newer_automation/Chainable.sol#L13) is not in mixedCase

src/v0.8/newer_automation/Chainable.sol#L13


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [tx-origin](#tx-origin) (1 results) (Medium)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
## tx-origin
Impact: Medium
Confidence: Medium
 - [ ] ID-0
[ExecutionPrevention._preventExecution()](src/v0.8/newer_automation/ExecutionPrevention.sol#L11-L16) uses tx.origin for authorization: [tx.origin != address(0)](src/v0.8/newer_automation/ExecutionPrevention.sol#L13)

src/v0.8/newer_automation/ExecutionPrevention.sol#L11-L16


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-1
[ExecutionPrevention._preventExecution()](src/v0.8/newer_automation/ExecutionPrevention.sol#L11-L16) is never used and should be removed

src/v0.8/newer_automation/ExecutionPrevention.sol#L11-L16


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/ExecutionPrevention.sol#L2)

src/v0.8/newer_automation/ExecutionPrevention.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [tx-origin](#tx-origin) (1 results) (Medium)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
## tx-origin
Impact: Medium
Confidence: Medium
 - [ ] ID-0
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) uses tx.origin for authorization: [tx.origin != address(0) && tx.origin != address(0x1111111111111111111111111111111111111111)](src/v0.8/newer_automation/AutomationBase.sol#L13)

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-1
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) is never used and should be removed

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/AutomationBase.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/KeeperBase.sol#L5)

src/v0.8/newer_automation/AutomationBase.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [tx-origin](#tx-origin) (1 results) (Medium)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (1 results) (Informational)
 - [unimplemented-functions](#unimplemented-functions) (1 results) (Informational)
## tx-origin
Impact: Medium
Confidence: Medium
 - [ ] ID-0
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) uses tx.origin for authorization: [tx.origin != address(0) && tx.origin != address(0x1111111111111111111111111111111111111111)](src/v0.8/newer_automation/AutomationBase.sol#L13)

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-1
[AutomationBase._preventExecution()](src/v0.8/newer_automation/AutomationBase.sol#L11-L16) is never used and should be removed

src/v0.8/newer_automation/AutomationBase.sol#L11-L16


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/AutomationBase.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/AutomationCompatible.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/KeeperCompatible.sol#L5)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)

src/v0.8/newer_automation/AutomationBase.sol#L2


## unimplemented-functions
Impact: Informational
Confidence: High
 - [ ] ID-3
[AutomationCompatible](src/v0.8/newer_automation/AutomationCompatible.sol#L7) does not implement functions:
	- [AutomationCompatibleInterface.checkUpkeep(bytes)](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L23)
	- [AutomationCompatibleInterface.performUpkeep(bytes)](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L41)

src/v0.8/newer_automation/AutomationCompatible.sol#L7


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/UpkeepFormat.sol#L3)

src/v0.8/newer_automation/UpkeepFormat.sol#L3


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/interfaces/TypeAndVersionInterface.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/UpkeepFormat.sol#L3)
	- [^0.8.0](src/v0.8/newer_automation/UpkeepTranscoder.sol#L3)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/UpkeepTranscoderInterface.sol#L2)

src/v0.8/interfaces/TypeAndVersionInterface.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [shadowing-state](#shadowing-state) (2 results) (High)
 - [unused-return](#unused-return) (1 results) (Medium)
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (3 results) (Informational)
## shadowing-state
Impact: High
Confidence: High
 - [ ] ID-0
[ArbitrumModule.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L20) shadows:
	- [ChainModuleBase.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L8)

src/v0.8/newer_automation/chains/ArbitrumModule.sol#L20


 - [ ] ID-1
[ArbitrumModule.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L19) shadows:
	- [ChainModuleBase.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L7)

src/v0.8/newer_automation/chains/ArbitrumModule.sol#L19


## unused-return
Impact: Medium
Confidence: Medium
 - [ ] ID-2
[ArbitrumModule.getMaxL1Fee(uint256)](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L38-L41) ignores return value by [(None,perL1CalldataByte,None,None,None,None) = ARB_GAS.getPricesInWei()](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L39)

src/v0.8/newer_automation/chains/ArbitrumModule.sol#L38-L41


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-3
3 different versions of Solidity are used:
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L2)
		-[0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)
	- Version constraint >=0.4.21<0.9.0 is used by:
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

src/v0.8/newer_automation/chains/ArbitrumModule.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-4
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

src/v0.8/newer_automation/interfaces/IChainModule.sol#L2


 - [ ] ID-5
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/chains/ArbitrumModule.sol#L2)
	- [0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)

src/v0.8/newer_automation/chains/ArbitrumModule.sol#L2


 - [ ] ID-6
Version constraint >=0.4.21<0.9.0 is too complex.
It is used by:
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (2 results) (Informational)
## pragma
Impact: Informational
Confidence: High
 - [ ] ID-0
2 different versions of Solidity are used:
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-1
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)

src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2


 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

src/v0.8/newer_automation/interfaces/IChainModule.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [shadowing-state](#shadowing-state) (2 results) (High)
 - [shadowing-local](#shadowing-local) (1 results) (Low)
 - [pragma](#pragma) (1 results) (Informational)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (3 results) (Informational)
 - [naming-convention](#naming-convention) (8 results) (Informational)
## shadowing-state
Impact: High
Confidence: High
 - [ ] ID-0
[OptimismModule.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/OptimismModule.sol#L17) shadows:
	- [ChainModuleBase.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L8)

src/v0.8/newer_automation/chains/OptimismModule.sol#L17


 - [ ] ID-1
[OptimismModule.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/OptimismModule.sol#L16) shadows:
	- [ChainModuleBase.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L7)

src/v0.8/newer_automation/chains/OptimismModule.sol#L16


## shadowing-local
Impact: Low
Confidence: High
 - [ ] ID-2
[OVM_GasPriceOracle.constructor(address)._owner](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41) shadows:
	- [Ownable._owner](node_modules/@openzeppelin/contracts/access/Ownable.sol#L21) (state variable)

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-3
3 different versions of Solidity are used:
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
		-[^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
		-[^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
		-[0.8.19](src/v0.8/newer_automation/chains/OptimismModule.sol#L2)
	- Version constraint ^0.8.9 is used by:
		-[^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-4
[Context._msgData()](node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23) is never used and should be removed

node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-5
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
	- [0.8.19](src/v0.8/newer_automation/chains/OptimismModule.sol#L2)

src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2


 - [ ] ID-6
Version constraint ^0.8.9 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation.
It is used by:
	- [^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2


 - [ ] ID-7
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
	- [^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-8
Contract [OVM_GasPriceOracle](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162) is not in CapWords

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162


 - [ ] ID-9
Parameter [OVM_GasPriceOracle.getL1GasUsed(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150


 - [ ] ID-10
Parameter [OVM_GasPriceOracle.setGasPrice(uint256)._gasPrice](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64


 - [ ] ID-11
Parameter [OVM_GasPriceOracle.setL1BaseFee(uint256)._baseFee](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74


 - [ ] ID-12
Parameter [OVM_GasPriceOracle.setDecimals(uint256)._decimals](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104


 - [ ] ID-13
Parameter [OVM_GasPriceOracle.setScalar(uint256)._scalar](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94


 - [ ] ID-14
Parameter [OVM_GasPriceOracle.getL1Fee(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117


 - [ ] ID-15
Parameter [OVM_GasPriceOracle.setOverhead(uint256)._overhead](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [shadowing-state](#shadowing-state) (2 results) (High)
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (3 results) (Informational)
## shadowing-state
Impact: High
Confidence: High
 - [ ] ID-0
[ScrollModule.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ScrollModule.sol#L18) shadows:
	- [ChainModuleBase.FIXED_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L7)

src/v0.8/newer_automation/chains/ScrollModule.sol#L18


 - [ ] ID-1
[ScrollModule.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ScrollModule.sol#L19) shadows:
	- [ChainModuleBase.PER_CALLDATA_BYTE_GAS_OVERHEAD](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L8)

src/v0.8/newer_automation/chains/ScrollModule.sol#L19


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-2
3 different versions of Solidity are used:
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
		-[0.8.19](src/v0.8/newer_automation/chains/ScrollModule.sol#L2)
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)
	- Version constraint ^0.8.16 is used by:
		-[^0.8.16](src/v0.8/vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol#L2)

src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-3
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

src/v0.8/newer_automation/interfaces/IChainModule.sol#L2


 - [ ] ID-4
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2)
	- [0.8.19](src/v0.8/newer_automation/chains/ScrollModule.sol#L2)

src/v0.8/newer_automation/chains/ChainModuleBase.sol#L2


 - [ ] ID-5
Version constraint ^0.8.16 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- StorageWriteRemovalBeforeConditionalTermination.
It is used by:
	- [^0.8.16](src/v0.8/vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol#L2)

src/v0.8/vendor/@scroll-tech/contracts/src/L2/predeploys/IScrollL1GasPriceOracle.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [reentrancy-no-eth](#reentrancy-no-eth) (1 results) (Medium)
 - [unused-return](#unused-return) (1 results) (Medium)
 - [shadowing-local](#shadowing-local) (1 results) (Low)
 - [calls-loop](#calls-loop) (1 results) (Low)
 - [reentrancy-events](#reentrancy-events) (1 results) (Low)
 - [assembly](#assembly) (1 results) (Informational)
 - [pragma](#pragma) (1 results) (Informational)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (4 results) (Informational)
 - [naming-convention](#naming-convention) (11 results) (Informational)
 - [cache-array-length](#cache-array-length) (1 results) (Optimization)
## reentrancy-no-eth
Impact: Medium
Confidence: Medium
 - [ ] ID-0
Reentrancy in [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187):
	External calls:
	- [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)
	State variables written after the call(s):
	- [s_feedMapping[feedId].bid = report.bid](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L174)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].ask = report.ask](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L175)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].price = report.price](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L176)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].observationsTimestamp = report.observationsTimestamp](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L177)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


## unused-return
Impact: Medium
Confidence: Medium
 - [ ] ID-1
[ChainSpecificUtil._getL1CalldataGasCost(uint256)](src/v0.8/ChainSpecificUtil.sol#L104-L115) ignores return value by [(None,l1PricePerByte,None,None,None,None) = ARBGAS.getPricesInWei()](src/v0.8/ChainSpecificUtil.sol#L107)

src/v0.8/ChainSpecificUtil.sol#L104-L115


## shadowing-local
Impact: Low
Confidence: High
 - [ ] ID-2
[OVM_GasPriceOracle.constructor(address)._owner](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41) shadows:
	- [Ownable._owner](node_modules/@openzeppelin/contracts/access/Ownable.sol#L21) (state variable)

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41


## calls-loop
Impact: Low
Confidence: Medium
 - [ ] ID-3
[MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187) has external calls inside a loop: [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


## reentrancy-events
Impact: Low
Confidence: Medium
 - [ ] ID-4
Reentrancy in [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187):
	External calls:
	- [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)
	Event emitted after the call(s):
	- [FeedUpdated(report.observationsTimestamp,report.price,report.bid,report.ask,feedId)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L180)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


## assembly
Impact: Informational
Confidence: High
 - [ ] ID-5
[MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L139-L141)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-6
4 different versions of Solidity are used:
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
		-[^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
		-[^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)
		-[^0.8.0](src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2)
		-[^0.8.0](src/v0.8/shared/access/ConfirmedOwner.sol#L2)
		-[^0.8.0](src/v0.8/shared/access/ConfirmedOwnerWithProposal.sol#L2)
		-[^0.8.0](src/v0.8/shared/interfaces/IOwnable.sol#L2)
	- Version constraint ^0.8.9 is used by:
		-[^0.8.9](src/v0.8/ChainSpecificUtil.sol#L2)
		-[^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1)
	- Version constraint >=0.4.21<0.9.0 is used by:
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-7
[Context._msgData()](node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23) is never used and should be removed

node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-8
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1


 - [ ] ID-9
Version constraint >=0.4.21<0.9.0 is too complex.
It is used by:
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5


 - [ ] ID-10
Version constraint ^0.8.9 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation.
It is used by:
	- [^0.8.9](src/v0.8/ChainSpecificUtil.sol#L2)
	- [^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)

src/v0.8/ChainSpecificUtil.sol#L2


 - [ ] ID-11
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
	- [^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2)
	- [^0.8.0](src/v0.8/shared/access/ConfirmedOwner.sol#L2)
	- [^0.8.0](src/v0.8/shared/access/ConfirmedOwnerWithProposal.sol#L2)
	- [^0.8.0](src/v0.8/shared/interfaces/IOwnable.sol#L2)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-12
Contract [OVM_GasPriceOracle](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162) is not in CapWords

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162


 - [ ] ID-13
Variable [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79


 - [ ] ID-14
Parameter [OVM_GasPriceOracle.getL1GasUsed(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150


 - [ ] ID-15
Variable [MercuryRegistry.s_verifier](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L74) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L74


 - [ ] ID-16
Variable [MercuryRegistry.s_feeds](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L78) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L78


 - [ ] ID-17
Parameter [OVM_GasPriceOracle.setGasPrice(uint256)._gasPrice](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64


 - [ ] ID-18
Parameter [OVM_GasPriceOracle.setL1BaseFee(uint256)._baseFee](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74


 - [ ] ID-19
Parameter [OVM_GasPriceOracle.setDecimals(uint256)._decimals](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104


 - [ ] ID-20
Parameter [OVM_GasPriceOracle.setScalar(uint256)._scalar](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94


 - [ ] ID-21
Parameter [OVM_GasPriceOracle.getL1Fee(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117


 - [ ] ID-22
Parameter [OVM_GasPriceOracle.setOverhead(uint256)._overhead](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84


## cache-array-length
Impact: Optimization
Confidence: High
 - [ ] ID-23
Loop condition [i < s_feeds.length](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L260) should use cached array length instead of referencing `length` member of the storage array.
 
src/v0.8/newer_automation/dev/MercuryRegistry.sol#L260


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [reentrancy-no-eth](#reentrancy-no-eth) (1 results) (Medium)
 - [unused-return](#unused-return) (3 results) (Medium)
 - [shadowing-local](#shadowing-local) (1 results) (Low)
 - [events-maths](#events-maths) (1 results) (Low)
 - [calls-loop](#calls-loop) (2 results) (Low)
 - [reentrancy-events](#reentrancy-events) (1 results) (Low)
 - [assembly](#assembly) (2 results) (Informational)
 - [pragma](#pragma) (1 results) (Informational)
 - [dead-code](#dead-code) (1 results) (Informational)
 - [solc-version](#solc-version) (4 results) (Informational)
 - [naming-convention](#naming-convention) (14 results) (Informational)
 - [cache-array-length](#cache-array-length) (1 results) (Optimization)
## reentrancy-no-eth
Impact: Medium
Confidence: Medium
 - [ ] ID-0
Reentrancy in [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187):
	External calls:
	- [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)
	State variables written after the call(s):
	- [s_feedMapping[feedId].bid = report.bid](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L174)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].ask = report.ask](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L175)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].price = report.price](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L176)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)
	- [s_feedMapping[feedId].observationsTimestamp = report.observationsTimestamp](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L177)
	[MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) can be used in cross function reentrancies:
	- [MercuryRegistry._updateFeed(string,string,int192,uint32)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L276-L286)
	- [MercuryRegistry.addFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L235-L251)
	- [MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145)
	- [MercuryRegistry.getLatestFeedData(string[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L95-L102)
	- [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187)
	- [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79)
	- [MercuryRegistry.setFeeds(string[],string[],int192[],uint32[])](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L253-L274)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


## unused-return
Impact: Medium
Confidence: Medium
 - [ ] ID-1
[ChainSpecificUtil._getL1CalldataGasCost(uint256)](src/v0.8/ChainSpecificUtil.sol#L104-L115) ignores return value by [(None,l1PricePerByte,None,None,None,None) = ARBGAS.getPricesInWei()](src/v0.8/ChainSpecificUtil.sol#L107)

src/v0.8/ChainSpecificUtil.sol#L104-L115


 - [ ] ID-2
[MercuryRegistryBatchUpkeep.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L57-L62) ignores return value by [i_registry.checkCallback(values,lookupData)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L61)

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L57-L62


 - [ ] ID-3
[MercuryRegistryBatchUpkeep.checkUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54) ignores return value by [i_registry.revertForFeedLookup(feeds)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L53)

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54


## shadowing-local
Impact: Low
Confidence: High
 - [ ] ID-4
[OVM_GasPriceOracle.constructor(address)._owner](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41) shadows:
	- [Ownable._owner](node_modules/@openzeppelin/contracts/access/Ownable.sol#L21) (state variable)

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L41


## events-maths
Impact: Low
Confidence: Medium
 - [ ] ID-5
[MercuryRegistryBatchUpkeep.updateBatchingWindow(uint256,uint256)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L77-L87) should emit an event for: 
	- [s_batchStart = batchStart](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L85) 
	- [s_batchEnd = batchEnd](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L86) 

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L77-L87


## calls-loop
Impact: Low
Confidence: Medium
 - [ ] ID-6
[MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187) has external calls inside a loop: [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


 - [ ] ID-7
[MercuryRegistryBatchUpkeep.checkUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54) has external calls inside a loop: [f = i_registry.s_feeds(i)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L36-L40)

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54


## reentrancy-events
Impact: Low
Confidence: Medium
 - [ ] ID-8
Reentrancy in [MercuryRegistry.performUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187):
	External calls:
	- [report = abi.decode(s_verifier.verify(values[i]),(Report))](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L160)
	Event emitted after the call(s):
	- [FeedUpdated(report.observationsTimestamp,report.price,report.bid,report.ask,feedId)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L180)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L156-L187


## assembly
Impact: Informational
Confidence: High
 - [ ] ID-9
[MercuryRegistryBatchUpkeep.checkUpkeep(bytes)](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L49-L51)

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L26-L54


 - [ ] ID-10
[MercuryRegistry.checkCallback(bytes[],bytes)](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145) uses assembly
	- [INLINE ASM](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L139-L141)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L118-L145


## pragma
Impact: Informational
Confidence: High
 - [ ] ID-11
4 different versions of Solidity are used:
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
		-[^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
		-[^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)
		-[^0.8.0](src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2)
		-[^0.8.0](src/v0.8/shared/access/ConfirmedOwner.sol#L2)
		-[^0.8.0](src/v0.8/shared/access/ConfirmedOwnerWithProposal.sol#L2)
		-[^0.8.0](src/v0.8/shared/interfaces/IOwnable.sol#L2)
	- Version constraint ^0.8.9 is used by:
		-[^0.8.9](src/v0.8/ChainSpecificUtil.sol#L2)
		-[^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)
	- Version constraint 0.8.19 is used by:
		-[0.8.19](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1)
		-[0.8.19](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L1)
	- Version constraint >=0.4.21<0.9.0 is used by:
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
		-[>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## dead-code
Impact: Informational
Confidence: Medium
 - [ ] ID-12
[Context._msgData()](node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23) is never used and should be removed

node_modules/@openzeppelin/contracts/utils/Context.sol#L21-L23


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-13
Version constraint 0.8.19 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess.
It is used by:
	- [0.8.19](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1)
	- [0.8.19](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L1)

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L1


 - [ ] ID-14
Version constraint >=0.4.21<0.9.0 is too complex.
It is used by:
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5)
	- [>=0.4.21<0.9.0](src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbSys.sol#L5)

src/v0.8/vendor/@arbitrum/nitro-contracts/src/precompiles/ArbGasInfo.sol#L5


 - [ ] ID-15
Version constraint ^0.8.9 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- VerbatimInvalidDeduplication
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation.
It is used by:
	- [^0.8.9](src/v0.8/ChainSpecificUtil.sol#L2)
	- [^0.8.9](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L2)

src/v0.8/ChainSpecificUtil.sol#L2


 - [ ] ID-16
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](node_modules/@openzeppelin/contracts/access/Ownable.sol#L4)
	- [^0.8.0](node_modules/@openzeppelin/contracts/utils/Context.sol#L4)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2)
	- [^0.8.0](src/v0.8/shared/access/ConfirmedOwner.sol#L2)
	- [^0.8.0](src/v0.8/shared/access/ConfirmedOwnerWithProposal.sol#L2)
	- [^0.8.0](src/v0.8/shared/interfaces/IOwnable.sol#L2)

node_modules/@openzeppelin/contracts/access/Ownable.sol#L4


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-17
Variable [MercuryRegistryBatchUpkeep.s_batchStart](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L16) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L16


 - [ ] ID-18
Variable [MercuryRegistryBatchUpkeep.s_batchEnd](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L17) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L17


 - [ ] ID-19
Contract [OVM_GasPriceOracle](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162) is not in CapWords

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L18-L162


 - [ ] ID-20
Variable [MercuryRegistry.s_feedMapping](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L79


 - [ ] ID-21
Parameter [OVM_GasPriceOracle.getL1GasUsed(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L150


 - [ ] ID-22
Variable [MercuryRegistry.s_verifier](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L74) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L74


 - [ ] ID-23
Variable [MercuryRegistry.s_feeds](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L78) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistry.sol#L78


 - [ ] ID-24
Parameter [OVM_GasPriceOracle.setGasPrice(uint256)._gasPrice](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L64


 - [ ] ID-25
Parameter [OVM_GasPriceOracle.setL1BaseFee(uint256)._baseFee](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L74


 - [ ] ID-26
Parameter [OVM_GasPriceOracle.setDecimals(uint256)._decimals](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L104


 - [ ] ID-27
Parameter [OVM_GasPriceOracle.setScalar(uint256)._scalar](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L94


 - [ ] ID-28
Parameter [OVM_GasPriceOracle.getL1Fee(bytes)._data](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L117


 - [ ] ID-29
Variable [MercuryRegistryBatchUpkeep.i_registry](src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L14) is not in mixedCase

src/v0.8/newer_automation/dev/MercuryRegistryBatchUpkeep.sol#L14


 - [ ] ID-30
Parameter [OVM_GasPriceOracle.setOverhead(uint256)._overhead](src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84) is not in mixedCase

src/v0.8/vendor/@eth-optimism/contracts/v0.8.9/contracts/L2/predeploys/OVM_GasPriceOracle.sol#L84


## cache-array-length
Impact: Optimization
Confidence: High
 - [ ] ID-31
Loop condition [i < s_feeds.length](src/v0.8/newer_automation/dev/MercuryRegistry.sol#L260) should use cached array length instead of referencing `length` member of the storage array.
 
src/v0.8/newer_automation/dev/MercuryRegistry.sol#L260


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)

src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [pragma](#pragma) (1 results) (Informational)
 - [solc-version](#solc-version) (2 results) (Informational)
## pragma
Impact: Informational
Confidence: High
 - [ ] ID-0
2 different versions of Solidity are used:
	- Version constraint ^0.8.0 is used by:
		-[^0.8.0](src/v0.8/newer_automation/interfaces/IAutomationForwarder.sol#L2)
		-[^0.8.0](src/v0.8/shared/interfaces/ITypeAndVersion.sol#L2)
	- Version constraint ^0.8.4 is used by:
		-[^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationForwarder.sol#L2


## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-1
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2


 - [ ] ID-2
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IAutomationForwarder.sol#L2)
	- [^0.8.0](src/v0.8/shared/interfaces/ITypeAndVersion.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationForwarder.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationRegistryConsumer.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/IAutomationV21PlusCommon.sol#L2)

src/v0.8/newer_automation/interfaces/IAutomationV21PlusCommon.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/IChainModule.sol#L2)

src/v0.8/newer_automation/interfaces/IChainModule.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/ILogAutomation.sol#L2)

src/v0.8/newer_automation/interfaces/ILogAutomation.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/KeeperCompatibleInterface.sol#L5)

src/v0.8/newer_automation/interfaces/AutomationCompatibleInterface.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/UpkeepFormat.sol#L3)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/MigratableKeeperRegistryInterface.sol#L3)

src/v0.8/newer_automation/UpkeepFormat.sol#L3


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/MigratableKeeperRegistryInterfaceV2.sol#L2)

src/v0.8/newer_automation/interfaces/MigratableKeeperRegistryInterfaceV2.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2)

src/v0.8/newer_automation/interfaces/StreamsLookupCompatibleInterface.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/UpkeepFormat.sol#L3)
	- [^0.8.0](src/v0.8/newer_automation/interfaces/UpkeepTranscoderInterface.sol#L2)

src/v0.8/newer_automation/UpkeepFormat.sol#L3


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/UpkeepTranscoderInterfaceV2.sol#L2)

src/v0.8/newer_automation/interfaces/UpkeepTranscoderInterfaceV2.sol#L2


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/v2_1/IKeeperRegistryMaster.sol#L4)

src/v0.8/newer_automation/interfaces/v2_1/IKeeperRegistryMaster.sol#L4


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
 - [naming-convention](#naming-convention) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/v2_2/IAutomationRegistryMaster.sol#L4)

src/v0.8/newer_automation/interfaces/v2_2/IAutomationRegistryMaster.sol#L4


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-1
Contract [AutomationRegistryBase2_2](src/v0.8/newer_automation/interfaces/v2_2/IAutomationRegistryMaster.sol#L270-L290) is not in CapWords

src/v0.8/newer_automation/interfaces/v2_2/IAutomationRegistryMaster.sol#L270-L290


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
 - [naming-convention](#naming-convention) (2 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.4 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables.
It is used by:
	- [^0.8.4](src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L4)

src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L4


## naming-convention
Impact: Informational
Confidence: High
 - [ ] ID-1
Contract [AutomationRegistryBase2_3](src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L309-L384) is not in CapWords

src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L309-L384


 - [ ] ID-2
Contract [IAutomationRegistryMaster2_3](src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L6-L307) is not in CapWords

src/v0.8/newer_automation/interfaces/v2_3/IAutomationRegistryMaster2_3.sol#L6-L307


**THIS CHECKLIST IS NOT COMPLETE**. Use `--show-ignored-findings` to show all the results.
Summary
 - [solc-version](#solc-version) (1 results) (Informational)
## solc-version
Impact: Informational
Confidence: High
 - [ ] ID-0
Version constraint ^0.8.0 contains known severe issues (https://solidity.readthedocs.io/en/latest/bugs.html)
	- FullInlinerNonExpressionSplitArgumentEvaluationOrder
	- MissingSideEffectsOnSelectorAccess
	- AbiReencodingHeadOverflowWithStaticArrayCleanup
	- DirtyBytesArrayToStorage
	- DataLocationChangeInInternalOverride
	- NestedCalldataArrayAbiReencodingSizeValidation
	- SignedImmutables
	- ABIDecodeTwoDimensionalArrayMemory
	- KeccakCaching.
It is used by:
	- [^0.8.0](src/v0.8/newer_automation/interfaces/v2_3/IWrappedNative.sol#L2)
	- [^0.8.0](src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol#L4)
	- [^0.8.0](src/v0.8/vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/extensions/IERC20Metadata.sol#L4)

src/v0.8/newer_automation/interfaces/v2_3/IWrappedNative.sol#L2


