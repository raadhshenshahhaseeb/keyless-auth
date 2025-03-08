#!/bin/bash

function generate_groth16() {
    cd circuit && go run main.go -groth16
    cd ..
}

function generate_expander() {
    cd circuit && go run main.go -expander
    cd ..
}

function deploy_contract_merkle() {
    cd contracts && npx hardhat ignition deploy ignition/modules/MerkleStorage.ts --network sepolia --strategy create2
    cd ..
}

function deploy_contract_verifier() {
    cd contracts && npx hardhat ignition deploy ignition/modules/Verifier.ts --network sepolia --strategy create2
    cd ..
}

function deploy_contracts() {
    deploy_contract_merkle
    deploy_contract_verifier
}

function deploy_gcloud {
    # Check required environment variables are set
    required_vars=("APP_PORT" "REDIS_USERNAME" "REDIS_HOST" "REDIS_PORT" "REDIS_PASSWORD")
    
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            echo "Error: Required environment variable $var is not set"
            exit 1
        fi
    done
    cd pkg && gcloud run deploy zk-be --source . --set-env-vars="APP_PORT=$APP_PORT,REDIS_USERNAME=$REDIS_USERNAME,REDIS_HOST=$REDIS_HOST,REDIS_PORT=$REDIS_PORT,REDIS_PASSWORD=$REDIS_PASSWORD"
    cd ..
}

# Parse command line argument
case "$1" in
    "groth16")
        generate_groth16
        ;;
    "expander") 
        generate_expander
        ;;
    "merkle")
        deploy_contract_merkle
        ;;
    "verifier")
        deploy_contract_verifier
        ;;
    "deploy-all")
        deploy_contracts
        ;;
    "deploy-api")
        deploy_gcloud
        ;;
    *)
        echo "Usage: $0 [groth16|expander|merkle|verifier|deploy-all|gcloud]"
        echo "  groth16    - Generate Groth16 proof system"
        echo "  expander   - Generate expander"
        echo "  merkle     - Deploy Merkle storage contract"
        echo "  verifier   - Deploy verifier contract" 
        echo "  deploy-all - Deploy all contracts"
        echo "  gcloud     - Deploy to Google Cloud"
        exit 1
        ;;
esac




