// This setup uses Hardhat Ignition to manage smart contract deployments.
// Learn more about it at https://hardhat.org/ignition

import { buildModule } from "@nomicfoundation/hardhat-ignition/modules";

const MerkleStorageModule = buildModule("MerkleStorageModule", (m) => {
  const merkleStorage = m.contract("MerkleRootStore");

  return { merkleStorage };
});

export default MerkleStorageModule;
