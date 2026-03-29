// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract Bank is Ownable {
    IERC20 public immutable token;

    mapping(address => uint256) public balanceOf;

    event Deposit(address indexed user, uint256 amount, uint256 timestamp);
    event Withdrawal(address indexed user, uint256 amount, uint256 timestamp);
    event Transfer(address indexed from, address indexed to, uint256 amount, uint256 timestamp);

    constructor(address _token) Ownable(msg.sender) {
        require(_token != address(0), "Invalid token address");
        token = IERC20(_token);
    }

    function deposit(uint256 amount) external {
        require(amount > 0, "Amount must be greater than 0");
        require(
            token.transferFrom(msg.sender, address(this), amount),
            "Transfer failed"
        );

        balanceOf[msg.sender] += amount;

        emit Deposit(msg.sender, amount, block.timestamp);
    }

    function withdraw(uint256 amount) external {
        require(amount > 0, "Amount must be greater than 0");
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");

        balanceOf[msg.sender] -= amount;
        require(token.transfer(msg.sender, amount), "Transfer failed");

        emit Withdrawal(msg.sender, amount, block.timestamp);
    }

    function transfer(address to, uint256 amount) external {
        require(to != address(0), "Invalid recipient");
        require(amount > 0, "Amount must be greater than 0");
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");

        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;

        emit Transfer(msg.sender, to, amount, block.timestamp);
    }

    function getBalance(address user) external view returns (uint256) {
        return balanceOf[user];
    }
}
