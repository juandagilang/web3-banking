import { ethers } from "hardhat";

async function main() {
  console.log("Deploying contracts...");

  // Deploy BankToken
  const Token = await ethers.getContractFactory("BankToken");
  const token = await Token.deploy();
  await token.waitForDeployment();
  const tokenAddress = await token.getAddress();
  console.log(`BankToken deployed to: ${tokenAddress}`);

  // Deploy Bank
  const Bank = await ethers.getContractFactory("Bank");
  const bank = await Bank.deploy(tokenAddress);
  await bank.waitForDeployment();
  const bankAddress = await bank.getAddress();
  console.log(`Bank deployed to: ${bankAddress}`);

  // Transfer all tokens to Bank
  const totalSupply = await token.totalSupply();
  await token.transfer(bankAddress, totalSupply);
  console.log(`Transferred ${totalSupply} tokens to Bank`);

  // Save addresses
  console.log("\n=== Contract Addresses ===");
  console.log(`TOKEN_ADDRESS=${tokenAddress}`);
  console.log(`BANK_ADDRESS=${bankAddress}`);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
