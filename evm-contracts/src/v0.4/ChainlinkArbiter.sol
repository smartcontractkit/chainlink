pragma solidity 0.4.24;

library SafeMath {

    /**
    * @dev Multiplies two numbers, reverts on overflow.
    */
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
        // benefit is lost if 'b' is also tested.
        // See: https://github.com/OpenZeppelin/openzeppelin-solidity/pull/522
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        require(c / a == b);

        return c;
    }

    /**
    * @dev Integer division of two numbers truncating the quotient, reverts on division by zero.
    */
    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b > 0); // Solidity only automatically asserts when dividing by 0
        uint256 c = a / b;
        // assert(a == b * c + a % b); // There is no case in which this doesn't hold

        return c;
    }

    /**
    * @dev Subtracts two numbers, reverts on overflow (i.e. if subtrahend is greater than minuend).
    */
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b <= a);
        uint256 c = a - b;

        return c;
    }

    /**
    * @dev Adds two numbers, reverts on overflow.
    */
    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        require(c >= a);

        return c;
    }

    /**
    * @dev Divides two numbers and returns the remainder (unsigned integer modulo),
    * reverts when dividing by zero.
    */
    function mod(uint256 a, uint256 b) internal pure returns (uint256) {
        require(b != 0);
        return a % b;
    }
}

interface ArbiterInterface {
    function getOndutyOracle() external returns(address, string);
}

// Arbiter contract is used to get the current chainlink oracle id
contract Arbiter is ArbiterInterface{
    using SafeMath for uint256;
    // chainlink oracleAddress => jobID , every oracle only put one jobID which parse to string
    mapping(address => string) public oraclesAndJobIDs;
    // ETH Super Node address => oracleAddress
    mapping(address => address) public superNodeAndOracles;
    // ETH Super Node address existence map
    mapping(address => bool) public superNodes;
    // current on duty super node address
    address[] public onDutySuperNodeAddresses;
    // valid on duty oracles
    uint public validOracles;
    // the last update height
    uint public updateHeight;
    // on duty oracles hash
    bytes32 public updateHash;

    // on duty oracle
    address public currentOracle;
    // on duty oracle job id
    string public currentJobId;

    // get current on duty oraces
    function getOndutyOracle() external returns(address, string){
        refresh();
        address addr = onDutySuperNodeAddresses[block.number.sub(updateHeight).mod(validOracles)];
        address oracle = superNodeAndOracles[addr];
        require(oracle != 0x0,"super node not register oracle yet");
        string storage jobId = oraclesAndJobIDs[oracle];
        currentOracle = oracle;
        currentJobId = jobId;
        return (addr,jobId);
    }

    function registerArbiter(address oracle, string jobId) public {
        refresh();
        // require(superNodes[msg.sender] == true , "sender must be one of the super node account");
        superNodeAndOracles[msg.sender] = oracle;
        oraclesAndJobIDs[oracle] = jobId;
    }

    function refresh() public {
        (address[] memory dposArbiter,uint validLength) = getArbiters();
        bytes32 newHash = keccak256(abi.encodePacked(dposArbiter));
        if (updateHash != newHash ){
            updateHash = newHash;
            updateHeight = block.number;
            onDutySuperNodeAddresses = dposArbiter;
            validOracles = validLength;
            for (uint i=0;i<dposArbiter.length;i++){
                address node = dposArbiter[i];
                if (superNodes[node] == false) {
                    superNodes[node] = true;
                }
            }
        }
    }

    function getArbiters() public view returns(address[],uint) {
        address[] storage addr;
        uint validAddressLength = 36;
        uint256[36] memory p;
        uint  input;
        assembly {
            if iszero(staticcall(gas, 20, input, 0x00, p, 0xC0)) {
                revert(0,0)
            }
        }
        for (uint i = 0;i<p.length;i++){
            addr.push(address(p[i]));
            if (p[i] == 0x0 && validAddressLength == 36){
                validAddressLength = i;
            }
        }
        return (addr,validAddressLength);
    }
}