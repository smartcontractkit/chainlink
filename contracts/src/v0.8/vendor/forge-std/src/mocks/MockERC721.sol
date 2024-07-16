// SPDX-License-Identifier: MIT
pragma solidity >=0.6.2 <0.9.0;

import {IERC721Metadata} from "../interfaces/IERC721.sol";

/// @notice This is a mock contract of the ERC721 standard for testing purposes only, it SHOULD NOT be used in production.
/// @dev Forked from: https://github.com/transmissions11/solmate/blob/0384dbaaa4fcb5715738a9254a7c0a4cb62cf458/src/tokens/ERC721.sol
contract MockERC721 is IERC721Metadata {
    /*//////////////////////////////////////////////////////////////
                         METADATA STORAGE/LOGIC
    //////////////////////////////////////////////////////////////*/

    string internal _name;

    string internal _symbol;

    function name() external view override returns (string memory) {
        return _name;
    }

    function symbol() external view override returns (string memory) {
        return _symbol;
    }

    function tokenURI(uint256 id) public view virtual override returns (string memory) {}

    /*//////////////////////////////////////////////////////////////
                      ERC721 BALANCE/OWNER STORAGE
    //////////////////////////////////////////////////////////////*/

    mapping(uint256 => address) internal _ownerOf;

    mapping(address => uint256) internal _balanceOf;

    function ownerOf(uint256 id) public view virtual override returns (address owner) {
        require((owner = _ownerOf[id]) != address(0), "NOT_MINTED");
    }

    function balanceOf(address owner) public view virtual override returns (uint256) {
        require(owner != address(0), "ZERO_ADDRESS");

        return _balanceOf[owner];
    }

    /*//////////////////////////////////////////////////////////////
                         ERC721 APPROVAL STORAGE
    //////////////////////////////////////////////////////////////*/

    mapping(uint256 => address) internal _getApproved;

    mapping(address => mapping(address => bool)) internal _isApprovedForAll;

    function getApproved(uint256 id) public view virtual override returns (address) {
        return _getApproved[id];
    }

    function isApprovedForAll(address owner, address operator) public view virtual override returns (bool) {
        return _isApprovedForAll[owner][operator];
    }

    /*//////////////////////////////////////////////////////////////
                               INITIALIZE
    //////////////////////////////////////////////////////////////*/

    /// @dev A bool to track whether the contract has been initialized.
    bool private initialized;

    /// @dev To hide constructor warnings across solc versions due to different constructor visibility requirements and
    /// syntaxes, we add an initialization function that can be called only once.
    function initialize(string memory name_, string memory symbol_) public {
        require(!initialized, "ALREADY_INITIALIZED");

        _name = name_;
        _symbol = symbol_;

        initialized = true;
    }

    /*//////////////////////////////////////////////////////////////
                              ERC721 LOGIC
    //////////////////////////////////////////////////////////////*/

    function approve(address spender, uint256 id) public payable virtual override {
        address owner = _ownerOf[id];

        require(msg.sender == owner || _isApprovedForAll[owner][msg.sender], "NOT_AUTHORIZED");

        _getApproved[id] = spender;

        emit Approval(owner, spender, id);
    }

    function setApprovalForAll(address operator, bool approved) public virtual override {
        _isApprovedForAll[msg.sender][operator] = approved;

        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function transferFrom(address from, address to, uint256 id) public payable virtual override {
        require(from == _ownerOf[id], "WRONG_FROM");

        require(to != address(0), "INVALID_RECIPIENT");

        require(
            msg.sender == from || _isApprovedForAll[from][msg.sender] || msg.sender == _getApproved[id],
            "NOT_AUTHORIZED"
        );

        // Underflow of the sender's balance is impossible because we check for
        // ownership above and the recipient's balance can't realistically overflow.
        _balanceOf[from]--;

        _balanceOf[to]++;

        _ownerOf[id] = to;

        delete _getApproved[id];

        emit Transfer(from, to, id);
    }

    function safeTransferFrom(address from, address to, uint256 id) public payable virtual override {
        transferFrom(from, to, id);

        require(
            !_isContract(to)
                || IERC721TokenReceiver(to).onERC721Received(msg.sender, from, id, "")
                    == IERC721TokenReceiver.onERC721Received.selector,
            "UNSAFE_RECIPIENT"
        );
    }

    function safeTransferFrom(address from, address to, uint256 id, bytes memory data)
        public
        payable
        virtual
        override
    {
        transferFrom(from, to, id);

        require(
            !_isContract(to)
                || IERC721TokenReceiver(to).onERC721Received(msg.sender, from, id, data)
                    == IERC721TokenReceiver.onERC721Received.selector,
            "UNSAFE_RECIPIENT"
        );
    }

    /*//////////////////////////////////////////////////////////////
                              ERC165 LOGIC
    //////////////////////////////////////////////////////////////*/

    function supportsInterface(bytes4 interfaceId) public view virtual override returns (bool) {
        return interfaceId == 0x01ffc9a7 // ERC165 Interface ID for ERC165
            || interfaceId == 0x80ac58cd // ERC165 Interface ID for ERC721
            || interfaceId == 0x5b5e139f; // ERC165 Interface ID for ERC721Metadata
    }

    /*//////////////////////////////////////////////////////////////
                        INTERNAL MINT/BURN LOGIC
    //////////////////////////////////////////////////////////////*/

    function _mint(address to, uint256 id) internal virtual {
        require(to != address(0), "INVALID_RECIPIENT");

        require(_ownerOf[id] == address(0), "ALREADY_MINTED");

        // Counter overflow is incredibly unrealistic.

        _balanceOf[to]++;

        _ownerOf[id] = to;

        emit Transfer(address(0), to, id);
    }

    function _burn(uint256 id) internal virtual {
        address owner = _ownerOf[id];

        require(owner != address(0), "NOT_MINTED");

        _balanceOf[owner]--;

        delete _ownerOf[id];

        delete _getApproved[id];

        emit Transfer(owner, address(0), id);
    }

    /*//////////////////////////////////////////////////////////////
                        INTERNAL SAFE MINT LOGIC
    //////////////////////////////////////////////////////////////*/

    function _safeMint(address to, uint256 id) internal virtual {
        _mint(to, id);

        require(
            !_isContract(to)
                || IERC721TokenReceiver(to).onERC721Received(msg.sender, address(0), id, "")
                    == IERC721TokenReceiver.onERC721Received.selector,
            "UNSAFE_RECIPIENT"
        );
    }

    function _safeMint(address to, uint256 id, bytes memory data) internal virtual {
        _mint(to, id);

        require(
            !_isContract(to)
                || IERC721TokenReceiver(to).onERC721Received(msg.sender, address(0), id, data)
                    == IERC721TokenReceiver.onERC721Received.selector,
            "UNSAFE_RECIPIENT"
        );
    }

    /*//////////////////////////////////////////////////////////////
                                HELPERS
    //////////////////////////////////////////////////////////////*/

    function _isContract(address _addr) private view returns (bool) {
        uint256 codeLength;

        // Assembly required for versions < 0.8.0 to check extcodesize.
        assembly {
            codeLength := extcodesize(_addr)
        }

        return codeLength > 0;
    }
}

interface IERC721TokenReceiver {
    function onERC721Received(address, address, uint256, bytes calldata) external returns (bytes4);
}
