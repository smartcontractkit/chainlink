// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {Internal} from "../../../libraries/Internal.sol";
import {MultiOCR3Base} from "../../../ocr/MultiOCR3Base.sol";
import {OffRamp} from "../../../offRamp/OffRamp.sol";
import {OffRampHelper} from "../../helpers/OffRampHelper.sol";
import {OffRampSetup} from "./OffRampSetup.t.sol";

contract OffRamp_afterOC3ConfigSet is OffRampSetup {
  function test_RevertWhen_afterOCR3ConfigSet_SignatureVerificationDisabled() public {
    s_offRamp = new OffRampHelper(
      OffRamp.StaticConfig({
        chainSelector: DEST_CHAIN_SELECTOR,
        gasForCallExactCheck: GAS_FOR_CALL_EXACT_CHECK,
        rmnRemote: s_mockRMNRemote,
        tokenAdminRegistry: address(s_tokenAdminRegistry),
        nonceManager: address(s_inboundNonceManager)
      }),
      _generateDynamicOffRampConfig(address(s_feeQuoter)),
      new OffRamp.SourceChainConfigArgs[](0)
    );

    MultiOCR3Base.OCRConfigArgs[] memory ocrConfigs = new MultiOCR3Base.OCRConfigArgs[](1);
    ocrConfigs[0] = MultiOCR3Base.OCRConfigArgs({
      ocrPluginType: uint8(Internal.OCRPluginType.Commit),
      configDigest: s_configDigestCommit,
      F: F,
      isSignatureVerificationEnabled: false,
      signers: s_validSigners,
      transmitters: s_validTransmitters
    });

    vm.expectRevert(OffRamp.SignatureVerificationRequiredInCommitPlugin.selector);
    s_offRamp.setOCR3Configs(ocrConfigs);
  }
}
