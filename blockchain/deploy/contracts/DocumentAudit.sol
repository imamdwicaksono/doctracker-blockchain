// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

/**
 * AUDIT LEDGER ONLY
 * - append-only
 * - no delete
 * - no update
 */
contract DocumentAudit {

    event DocumentEvent(
        string indexed docId,
        bytes32 hash,
        string action,
        string actor,
        uint256 timestamp
    );

    function record(
        string calldata docId,
        bytes32 hash,
        string calldata action,
        string calldata actor
    ) external {
        emit DocumentEvent(
            docId,
            hash,
            action,
            actor,
            block.timestamp
        );
    }
}
