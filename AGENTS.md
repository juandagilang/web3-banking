# Web3 Banking System - Project Agent Guide

## Goal

Build a decentralized banking system where users can:
- Connect wallet (Web3-only authentication)
- Deposit tokens to the bank (on-chain)
- Withdraw tokens from the bank (on-chain)
- Transfer tokens to other users (on-chain)
- View transaction history (paginated)

---

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go + Gin |
| Database | PostgreSQL |
| Blockchain | Solidity + Hardhat |
| Frontend | Vue 3 + Vite + Pinia |
| Web3 | viem + web3modal |
| Chain | Hardhat Node (local) |

---

## Current Status

### ✅ Completed
- Blockchain contracts deployed to Hardhat Node
- BankToken (ERC-20): `0x5FbDB2315678afecb367f032d93F642f64180aa3`
- Bank contract: `0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512`
- All 7 contract tests passing
- Hardhat Node running at `http://127.0.0.1:8545`

### 🔄 In Progress
- Frontend implementation

### ⏳ Pending
- Backend implementation (user will code with assistance)
- Frontend-backend integration

---

## Project Structure

```
Web3-Banking/
├── Backend/                    # Go + Gin + PostgreSQL (user implements)
│   ├── cmd/api/main.go
│   ├── config/
│   ├── internal/
│   │   ├── domain/entity/
│   │   ├── usecase/
│   │   ├── repository/
│   │   └── delivery/
│   ├── migrations/
│   ├── go.mod
│   └── .env
│
├── Blockchain/                 # ✅ Complete
│   ├── contracts/
│   │   ├── BankToken.sol
│   │   └── Bank.sol
│   ├── scripts/deploy.ts
│   ├── test/bank.test.ts
│   ├── hardhat.config.ts
│   └── package.json
│
└── Frontend/                   # Vue 3 + Vite (in progress)
    ├── src/
    │   ├── api/
    │   ├── components/
    │   ├── composables/
    │   ├── pages/
    │   ├── store/
    │   └── types/
    ├── vite.config.ts
    └── package.json
```

---

## Roadmap / Implementation Phases

### Phase 1: Blockchain (✅ COMPLETE)
- [x] Initialize Hardhat project
- [x] Create BankToken.sol (ERC-20, 1B supply)
- [x] Create Bank.sol (deposit, withdraw, transfer)
- [x] Write and run tests
- [x] Deploy to Hardhat Node

### Phase 2: Frontend
- [ ] Initialize Vue 3 + Vite project
- [ ] Install dependencies: viem, web3modal, pinia, vue-router
- [ ] Setup web3 connection (web3modal + viem)
- [ ] Create Login page (wallet connect + signature)
- [ ] Create Dashboard page
  - Display balance
  - Deposit form
  - Withdraw form
  - Transfer form
- [ ] Add transaction history with pagination
- [ ] Connect to Hardhat Node

### Phase 3: Backend (User Implements)
- [ ] Setup Go project with Gin
- [ ] Setup PostgreSQL connection
- [ ] Create migrations
- [ ] Implement repositories
- [ ] Implement usecases
- [ ] Implement HTTP handlers
- [ ] Implement event listener (sync Deposit/Withdrawal/Transfer events to DB)
- [ ] Endpoints:
  - `POST /api/auth/nonce` - Get nonce for wallet
  - `POST /api/auth/login` - Verify signature
  - `GET /api/account/:address` - Get balance
  - `GET /api/transactions?page=1&limit=20` - Paginated history
  - `GET /api/contract/info` - Token info

### Phase 4: Integration
- [ ] Connect frontend to backend
- [ ] Test full flow
- [ ] Fix bugs

---

## Important Contract Details

### Contract Addresses
```
TOKEN_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3
BANK_ADDRESS=0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512
```

### Contract Functions
```solidity
// Bank.sol
function deposit(uint256 amount) external
function withdraw(uint256 amount) external
function transfer(address to, uint256 amount) external
function balanceOf(address user) external view returns (uint256)
function token() external view returns (address)
```

### Events (for backend event listener)
```
Deposit(address indexed user, uint256 amount, uint256 timestamp)
Withdrawal(address indexed user, uint256 amount, uint256 timestamp)
Transfer(address indexed from, address indexed to, uint256 amount, uint256 timestamp)
```

### ABI Files
Located in:
- `Blockchain/artifacts/contracts/Bank.sol/Bank.json`
- `Blockchain/artifacts/contracts/BankToken.sol/BankToken.json`

---

## Key Design Decisions

1. **Web3-only auth**: No email/password, user signs a message with wallet
2. **1:1 token**: Each token in bank = 1 token deposited (no interest)
3. **On-chain**: All operations (deposit/withdraw/transfer) happen on blockchain
4. **PostgreSQL for indexing**: Backend indexes events for fast transaction history queries

---

## Notes for Future Agent

- Hardhat Node must be running before deploying contracts or testing: `npx hardhat node`
- Frontend uses viem for blockchain interaction, NOT ethers.js
- Backend user implements themselves - agent only assists
- All blockchain operations are on-chain (no off-chain ledger)
- Transaction history pagination uses DB indexing, not direct contract queries

---

## Commands Reference

```bash
# Blockchain
cd Blockchain
npm install
npx hardhat node              # Start local node (separate terminal)
npx hardhat compile           # Compile contracts
npx hardhat test              # Run tests
npx hardhat run scripts/deploy.ts --network hardhat  # Deploy

# Frontend (to be created)
cd Frontend
npm install
npm run dev

# Backend (to be created by user)
cd Backend
go mod init
go run cmd/api/main.go
```
