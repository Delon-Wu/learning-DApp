// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Counter {
    uint private _num;

    event Increased();
    event Decreased();

    function increase() public {
        _num++;
        emit Increased();
    }

    function decrease() public {
        require(_num > 0, "Num must greater than 0");
        _num--;
        emit Decreased();
    }
    function getNumber() public view returns (uint) {
        return _num;
    }
}