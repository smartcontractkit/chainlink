// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

import {Functions} from "../dev/functions/Functions.sol";

contract FunctionsTestHelper {
  using Functions for Functions.Request;

  error EmptySource();
  error EmptyUrl();
  error EmptySecrets();
  error EmptyArgs();

  Functions.Request private s_req;

  event RequestData(bytes data);

  function closeEvent() public {
    emit RequestData(s_req.encodeCBOR());
  }

  function initializeRequestForInlineJavaScript(string memory sourceCode) public {
    Functions.Request memory r;
    r.initializeRequestForInlineJavaScript(sourceCode);
    storeRequest(r);
  }

  function addSecrets(bytes memory secrets) public {
    Functions.Request memory r = s_req;
    r.addRemoteSecrets(secrets);
    storeRequest(r);
  }

  function addEmptyArgs() public pure {
    Functions.Request memory r;
    string[] memory args;
    r.addArgs(args);
  }

  function addTwoArgs(string memory arg1, string memory arg2) public {
    string[] memory args = new string[](2);
    args[0] = arg1;
    args[1] = arg2;
    Functions.Request memory r = s_req;
    r.addArgs(args);
    storeRequest(r);
  }

  function storeRequest(Functions.Request memory r) private {
    s_req.codeLocation = r.codeLocation;
    s_req.language = r.language;
    s_req.source = r.source;
    s_req.args = r.args;
    s_req.secretsLocation = r.secretsLocation;
    s_req.secrets = r.secrets;
  }
}
