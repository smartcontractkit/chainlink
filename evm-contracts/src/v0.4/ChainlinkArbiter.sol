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
    function score(uint data) external returns (uint);
}

// Arbiter contract is used to get the current chainlink oracle id
contract Arbiter is ArbiterInterface{
    using SafeMath for uint256;
    // chainlink oracleAddress => jobID , every oracle only put one jobID which parse to string
    mapping(address => string) public oraclesAndJobIDs;
    // ETH Super Node address => oracleAddress
    mapping(bytes32 => address) public superNodeAndOracles;
    // ETH Super Node pub key keccak256 hash existence map
    mapping(bytes32 => bool) public superNodes;
    // current on duty super node pub key keccak256 hash
    bytes32[] public onDutySuperNodes;
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
        bytes32 addr = onDutySuperNodes[block.number.sub(updateHeight).mod(validOracles)];
        address oracle = superNodeAndOracles[addr];
        require(oracle != 0x0,"super node not register oracle yet");
        string storage jobId = oraclesAndJobIDs[oracle];
        currentOracle = oracle;
        currentJobId = jobId;
        return (oracle,jobId);
    }

    function score(uint data) external returns (uint){
        //TODO to be implemented
        return data;
    }

    function registerArbiter(string publicKey, address oracle, string jobId, string signature) public {
        bytes32 oracleHash = keccak256(oracle);
        bytes32 jobIdHash = keccak256(jobId);
        bytes memory mergeHash = mergeBytes(oracleHash,jobIdHash);
        string memory data = toHex(mergeHash);
        require(p256_verify(publicKey,data,signature) == true,"verify signature error");
        refresh();
        bytes memory pubKeyBytes = hexStr2bytes(publicKey);
        bytes32 keyHash = keccak256(pubKeyBytes);
        require(superNodes[keyHash] == true , "sender must be one of the super node account");
        superNodeAndOracles[keyHash] = oracle;
        oraclesAndJobIDs[oracle] = jobId;
    }

    function refresh() public {
        uint validAddressLength = 36;
        bytes32[36] memory p;
        uint  input;
        assembly {
            if iszero(staticcall(gas, 20, input, 0x00, p, 0xC0)) {
                revert(0,0)
            }
        }
        for (uint i = 0;i<p.length;i++){
            if (p[i] == 0x0 && validAddressLength == 36){
                validAddressLength = i;
            }
        }
        bytes32 newHash = keccak256(abi.encodePacked(p));
        if (updateHash != newHash ){
            updateHash = newHash;
            updateHeight = block.number;
            onDutySuperNodes = p;
            validOracles = validAddressLength;
            for (uint j=0;j<p.length;j++){
                bytes32 node = p[j];
                if (superNodes[node] == false) {
                    superNodes[node] = true;
                }
            }
        }
    }

    function p256_verify(string pubkey, string data, string sig) public view returns(bool) {
        string memory i = strConcat(strConcat(pubkey, data), sig);
        bytes memory input = hexStr2bytes(i);
        uint256[1] memory p;

        assembly {
            if iszero(staticcall(gas, 21, input, 193, p, 0xc0)) {
                revert(0,0)
            }
        }

        return p[0] == 1;
    }

    function strConcat(string _a, string _b) internal returns (string){
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        string memory ret = new string(_ba.length + _bb.length);
        bytes memory bret = bytes(ret);
        uint k = 0;
        for (uint i = 0; i < _ba.length; i++)bret[k++] = _ba[i];
        for (i = 0; i < _bb.length; i++) bret[k++] = _bb[i];
        return string(ret);
    }

    function hexStr2bytes(string data) internal returns (bytes){

        bytes memory a = bytes(data);
        uint[] memory b = new uint[](a.length);

        for (uint i = 0; i < a.length; i++) {
            uint _a = uint(a[i]);

            if (_a > 96) {
                b[i] = _a - 97 + 10;
            }
            else if (_a > 66) {
                b[i] = _a - 65 + 10;
            }
            else {
                b[i] = _a - 48;
            }
        }

        bytes memory c = new bytes(b.length / 2);
        for (uint _i = 0; _i < b.length; _i += 2) {
            c[_i / 2] = byte(b[_i] * 16 + b[_i + 1]);
        }

        return c;
    }

    function toHex(bytes origin) returns(string) {
        bytes memory dst = new bytes(2 * origin.length);
        bytes memory hextable = "0123456789abcdef";
        uint j = 0;
        for (uint i= 0; i<origin.length;i++ ) {
            dst[j] = hextable[uint256(origin[i])>>4];
            dst[j+1] = hextable[uint256(origin[i])&0x0f];
            j = j+ 2;
        }
        return string(dst);
    }

    function mergeBytes(bytes32 first,bytes32 second) returns(bytes) {
        bytes memory merged = new bytes(first.length + second.length);
        uint k = 0;
        for (i = 0; i < first.length; i++) {
            merged[k] = first[i];
            k++;
        }

        for (uint i = 0; i < second.length; i++) {
            merged[k] = second[i];
            k++;
        }
        return merged;
    }
}