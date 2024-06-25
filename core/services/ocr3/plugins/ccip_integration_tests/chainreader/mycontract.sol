// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.18;

contract SimpleContract {
    event SimpleEvent(uint256 value);
    uint256 public eventCount;
    uint[] public numbers;

    struct Person {
        string name;
        uint age;
    }

    function emitEvent() public {
        eventCount++;
        numbers.push(eventCount);
        emit SimpleEvent(eventCount);
    }

    function getEventCount() public view returns (uint256) {
        return eventCount;
    }

    function getNumbers() public view returns (uint256[] memory) {
        return numbers;
    }

    function getPerson() public pure returns (Person memory) {
        return Person("Dim", 18);
    }
}
