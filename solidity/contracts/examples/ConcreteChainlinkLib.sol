pragma solidity ^0.4.18;

import "../ChainlinkLib.sol";

contract ConcreteChainlinkLib {
  using ChainlinkLib for ChainlinkLib.Run;

  ChainlinkLib.Run run;

  event RunData(
    bytes names,
    bytes types,
    bytes values,
    bytes payload
  );

  event Log(bytes val);

  function fireEvent()
    public
  {
    emit RunData(run.names, run.types, run.values, run.payload());
  }

  function add(string _key, string _value)
    public
  {
    ChainlinkLib.Run memory r2 = run;
    r2.add(_key, _value);
    run = r2;
  }

  function addBytes32(string _key, bytes32 _value)
    public
  {
    ChainlinkLib.Run memory r2 = run;
    r2.addBytes32(_key, _value);
    run = r2;
  }

  function addBytes32Array(string _key, bytes32[] memory _values)
    public
  {
    ChainlinkLib.Run memory r2 = run;
    r2.addBytes32Array(_key, _values);
    run = r2;
  }

}
