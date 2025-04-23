// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StudentPoints {
    struct UserInfo {
        string username;
        string phone;
        string roleName;
        uint256 points;
        string title;
    }

    mapping(address => UserInfo) public users;

    event UserInfoUpdated(address indexed user, string username, string phone, string roleName, uint256 points, string title);

function updateUserInfo(
    string memory username,
    string memory phone,
    string memory roleName,
    uint256 points,
    string memory title
) public {
    address userAddr = msg.sender;
    users[userAddr] = UserInfo(username, phone, roleName, points, title);
    emit UserInfoUpdated(userAddr, username, phone, roleName, points, title);
}


    function getUserInfo(address userAddr) public view returns (
        string memory, string memory, string memory, uint256, string memory
    ) {
        UserInfo memory u = users[userAddr];
        return (u.username, u.phone, u.roleName, u.points, u.title);
    }
}
