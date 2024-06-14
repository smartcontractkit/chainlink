// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {FunctionsRequest} from "../../../dev/v1_0_0/libraries/FunctionsRequest.sol";

contract FunctionsTestHelper {
  using FunctionsRequest for FunctionsRequest.Request;

  FunctionsRequest.Request private s_req;

  event RequestData(bytes data);

  function closeEvent() public {
    emit RequestData(s_req.encodeCBOR());
  }

  function initializeRequestForInlineJavaScript(string memory sourceCode) public {
    FunctionsRequest.Request memory r;
    r.initializeRequestForInlineJavaScript(sourceCode);
    storeRequest(r);
  }

  function addSecretsReference(bytes memory secrets) public {
    FunctionsRequest.Request memory r = s_req;
    r.addSecretsReference(secrets);
    storeRequest(r);
  }

  function addEmptyArgs() public pure {
    FunctionsRequest.Request memory r;
    string[] memory args;
    r.setArgs(args);
  }

  function addTwoArgs(string memory arg1, string memory arg2) public {
    string[] memory args = new string[](2);
    args[0] = arg1;
    args[1] = arg2;
    FunctionsRequest.Request memory r = s_req;
    r.setArgs(args);
    storeRequest(r);
  }

  function storeRequest(FunctionsRequest.Request memory r) private {
    s_req.codeLocation = r.codeLocation;
    s_req.language = r.language;
    s_req.source = r.source;
    s_req.args = r.args;
    s_req.secretsLocation = r.secretsLocation;
    s_req.encryptedSecretsReference = r.encryptedSecretsReference;
  }
}
