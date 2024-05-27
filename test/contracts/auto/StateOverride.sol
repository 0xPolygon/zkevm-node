// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract StateOverride {
    uint256 number = 1;
    string text = "text";

    function addrBalance(address a) public view returns (uint256) {
        return address(a).balance;
    }

    function getNumber() public view returns (uint256) {
        return number;
    }

    function getText() public view returns (string memory) {
        return text;
    }
}
