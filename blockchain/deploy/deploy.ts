import { ethers } from "hardhat";

async function main() {
  const Audit = await ethers.getContractFactory("DocumentAudit");
  const audit = await Audit.deploy();

  await audit.deployed();

  console.log("DocumentAudit deployed to:", audit.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
