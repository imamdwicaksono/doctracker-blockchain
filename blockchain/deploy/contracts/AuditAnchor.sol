// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract AuditAnchor {
    mapping(uint256 => bytes32) public blockHashes;

    event BlockAnchored(uint256 height, bytes32 hash);

    function anchorBlock(uint256 height, bytes32 hash) external {
        require(blockHashes[height] == bytes32(0), "ALREADY_ANCHORED");
        blockHashes[height] = hash;
        emit BlockAnchored(height, hash);
    }

    function verify(uint256 height, bytes32 hash) external view returns (bool) {
        return blockHashes[height] == hash;
    }
}
