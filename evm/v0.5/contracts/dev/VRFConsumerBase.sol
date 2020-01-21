pragma solidity 0.5.0;

////////////////////////////////////////////////////////////////////////////////
////////// DO NOT USE THIS IN PRODUCTION UNTIL IT HAS BEEN AUDITED /////////////
////////////////////////////////////////////////////////////////////////////////

import "../interfaces/LinkTokenInterface.sol";
import "./VRFRequestIDBase.sol";

/** ****************************************************************************
 * @notice Interface for contracts using VRF randomness
 * *****************************************************************************
 * @dev PURPOSE
 *
 * @dev Reggie the Random Oracle (not his real job) wants to provide randomness
 * @dev to Vera the verifier in such a way that Vera can be sure he's not
 * @dev making his output up to suit himself. Reggie provides Vera a public key
 * @dev to which he knows the secret key. Each time Vera provides a seed to
 * @dev Reggie, he gives back a value which is computed completely
 * @dev deterministically from the seed and the secret key.
 *
 * @dev Reggie provides a proof by which Vera can verify that the output was
 * @dev correctly computed once Reggie tells it to her, but without that proof,
 * @dev the output is indistinguishable to her from a uniform random sample
 * @dev from the output space.
 *
 * @dev The purpose of this contract is to make it easy for unrelated contracts
 * @dev to talk to Vera the verifier about the work Reggie is doing, to provide
 * @dev simple access to a verifiable source of randomness.
 * *****************************************************************************
 * @dev USAGE
 *
 * @dev Calling contracts must inherit from VRFConsumerInterface, and can
 * @dev initialize VRFConsumerInterface's attributes in their constructor as
 * @dev shown:
 *
 * @dev   contract VRFConsumer {
 * @dev     constuctor(<other arguments>, address _vrfCoordinator, address _link)
 * @dev       VRFConsumerBase(_vrfCoordinator, _link) public {
 * @dev         <initialization with other arguments goes here>
 * @dev       }
 * @dev   }
 *
 * @dev The oracle will have given you an ID for the VRF keypair they have
 * @dev committed to, call it keyHash, and have told you the minimum LINK price
 * @dev for VRF service. Make sure your contract has sufficient LINK, and call
 * @dev requestRandomness(keyHash, fee, seed), where seed is the input you want
 * @dev to generate randomness from.
 *
 * @dev Once the VRFCoordinator has received and validated the oracle's response
 * @dev to your request, it will call your contract's fulfillRandomness method.
 *
 * @dev The randomness argument to fulfillRandomness is the actual random value
 * @dev generated from your seed.
 *
 * @dev The requestId argument is generated from the keyHash and the seed by
 * @dev makeRequestId(keyHash, seed). If your contract could have concurrent
 * @dev requests open, you can use the requestId to track which seed is
 * @dev associated with which randomness. Collision of requestId's is
 * @dev cryptographically impossible. See VRFRequestIDBase.sol for more details.
 * *****************************************************************************
 * @dev SECURITY CONSIDERATIONS
 *
 * @dev To increase trust in your contract, the source of your seeds should be
 * @dev hard for anyone to influence. Any party who can influence them could in
 * @dev principle collude with the oracle (who can instantly compute the VRF
 * @dev output for any given seed) to bias the outcomes from your contract in
 * @dev their favor. For instance, the block hash is a natural choice of seed
 * @dev for many applications, but miners in control of a substantial fraction
 * @dev of hashing power and with access to VRF outputs could check the result
 * @dev of prospective block hashes as they are mined, and decide not to publish
 * @dev a block if they don't like the outcome it will lead to.
 *
 * @dev On the other hand, using block hashes as the seed makes it particularly
 * @dev easy to estimate the economic cost to a miner for this kind of cheating
 * @dev (namely, the block reward and transaction fees they forgo by refraining
 * @dev from publishing a block.)
 */
contract VRFConsumerBase is VRFRequestIDBase {
  /**
   * @notice fulfillRandomness handles the VRF response. Your contract must
   * @notice implement it.
   *
   * @dev The VRFCoordinator expects a calling contract to have a method with
   * @dev this signature, and will call it once it has verified the proof
   * @dev associated with the randomness.
   *
   * @param requestId keccak256(abi.encodePacked(keyHash, seed))
   * @param randomness the VRF output
   */
  function fulfillRandomness(bytes32 requestId, uint256 randomness) external;
  /**
   * @notice requestRandomness initiates a request for VRF output given _seed
   *
   * @dev The fulfillRandomness method receives the output, once it's provided
   * @dev by the Oracle, and verified by the vrfCoordinator.
   *
   * @dev The _keyHash must already be registered with the VRFCoordinator, and
   * @dev the _fee must exceed the fee specified during registration of the
   * @dev _keyHash.
   *
   * @param _keyHash ID of public key against which randomness is generated
   * @param _fee The amount of LINK to send with the request
   * @param _seed Random seed to input to VRF, from which output is determined
   */
  function requestRandomness(bytes32 _keyHash, uint256 _fee, uint256 _seed)
    external
  {
    LINK.transferAndCall(vrfCoordinator, _fee, abi.encode(_keyHash, _seed));
  }

  LinkTokenInterface LINK;
  address vrfCoordinator;

  constructor(address _vrfCoordinator, address _link) public {
    vrfCoordinator = _vrfCoordinator;
    LINK = LinkTokenInterface(_link);
  }
}
