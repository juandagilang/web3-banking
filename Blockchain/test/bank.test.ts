import { ethers } from "hardhat";
import { expect } from "chai";

describe("BankToken", () => {
  let token: any;

  beforeEach(async () => {
    const Token = await ethers.getContractFactory("BankToken");
    token = await Token.deploy();
    await token.waitForDeployment();
  });

  it("should have correct name and symbol", async () => {
    expect(await token.name()).to.equal("Web3Bank Token");
    expect(await token.symbol()).to.equal("W3B");
  });

  it("should have correct total supply", async () => {
    const totalSupply = await token.totalSupply();
    expect(totalSupply).to.equal(ethers.parseEther("1000000000"));
  });

  it("should mint tokens to owner", async () => {
    const [owner] = await ethers.getSigners();
    const balance = await token.balanceOf(owner.address);
    expect(balance).to.equal(ethers.parseEther("1000000000"));
  });
});

describe("Bank", () => {
  let bank: any;
  let token: any;
  let owner: any;
  let user1: any;
  let user2: any;

  beforeEach(async () => {
    [owner, user1, user2] = await ethers.getSigners();

    const Token = await ethers.getContractFactory("BankToken");
    token = await Token.deploy();
    await token.waitForDeployment();

    const Bank = await ethers.getContractFactory("Bank");
    bank = await Bank.deploy(await token.getAddress());
    await bank.waitForDeployment();

    // Transfer tokens to user1 for testing
    await token.transfer(user1.address, ethers.parseEther("1000"));
  });

  it("should accept deposits", async () => {
    const amount = ethers.parseEther("100");
    
    // Approve bank to spend tokens
    await token.connect(user1).approve(await bank.getAddress(), amount);
    
    // Deposit
    await bank.connect(user1).deposit(amount);

    // Check balance
    const balance = await bank.balanceOf(user1.address);
    expect(balance).to.equal(amount);
  });

  it("should allow withdrawals", async () => {
    const amount = ethers.parseEther("100");
    const initialBalance = await token.balanceOf(user1.address);

    // Approve and deposit
    await token.connect(user1).approve(await bank.getAddress(), amount);
    await bank.connect(user1).deposit(amount);

    // Withdraw
    await bank.connect(user1).withdraw(amount);

    // Check balance
    const bankBalance = await bank.balanceOf(user1.address);
    expect(bankBalance).to.equal(0);
    
    const tokenBalance = await token.balanceOf(user1.address);
    expect(tokenBalance).to.equal(initialBalance);
  });

  it("should allow transfers", async () => {
    const amount = ethers.parseEther("100");

    // Approve and deposit
    await token.connect(user1).approve(await bank.getAddress(), amount);
    await bank.connect(user1).deposit(amount);

    // Transfer to user2
    await bank.connect(user1).transfer(user2.address, amount);

    // Check balances
    expect(await bank.balanceOf(user1.address)).to.equal(0);
    expect(await bank.balanceOf(user2.address)).to.equal(amount);
  });

  it("should fail on insufficient balance", async () => {
    const amount = ethers.parseEther("100");

    await expect(bank.connect(user1).withdraw(amount)).to.be.revertedWith(
      "Insufficient balance"
    );
  });
});
