import React from 'react';

import { poseidon } from '@iden3/js-crypto'
import { MerkleTree } from 'merkletreejs'

const Demo: React.FC = () => {
    // 0. existing credentials
    const credentials = ["test@gmail.com", "test2@gmail.com", "test3@gmail.com", "test4@gmail.com"];
    const encoder = new TextEncoder();
    const leaves = credentials.map(cred => poseidon.hashBytes(encoder.encode(cred)).toString(16));

    // console.log("Leaves:", leaves);
    // 1. hash credentials
    const newCredential = "test-new@gmail.com";
    const hashedCredential = poseidon.hashBytes(encoder.encode(newCredential)).toString(16);
    leaves.push(hashedCredential);
    console.log("Leaves:", leaves);
    // 2. Build Merkle tree with new leaf
    const tree = new MerkleTree([hashedCredential, hashedCredential], poseidon, { sortPairs: true });
    console.log("Tree:", tree);
    // const root = tree.getRoot();
    
    // 3. Generate proof (leaf, root, positions)
    // const proof = tree.getProof(hashedCredential);
    // const positions = proof.map(p => p.position === 'left' ? 0 : 1); // 0 = left, 1 = right

    // const verified = tree.verify(proof, hashedCredential, root);

    return (
        <div style={{ color: 'white' }}>
            <h1>ZK auth demo</h1>
            <p>Hash: {hashedCredential}</p>
            {/* <p>Root: {root.toString('hex')}</p>
            <p>Proof: {JSON.stringify(proof)}</p>
            <p>Verified: {verified ? 'true' : 'false'}</p>
            <p>Positions: {JSON.stringify(positions)}</p> */}
        </div>
    );
};

export default Demo;