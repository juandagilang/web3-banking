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
- Backend API implemented (Go + Gin + PostgreSQL)
  - Auth endpoints (nonce, login)
  - Account endpoint (balance)
  - Transaction endpoint (paginated history)
  - Event listener (blockchain sync)

### 🔄 In Progress
- Frontend implementation (Vue 3 + Vite)

### ⏳ Pending
- Frontend-backend integration
- End-to-end testing

---

## Project Structure

```
Web3-Banking/
├── Backend/                    # ✅ Complete (Go + Gin + PostgreSQL)
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
    │   ├── api/              # Axios client, endpoints
    │   ├── components/        # Reusable UI components
    │   ├── composables/      # Vue composables (useWeb3, useAuth, useContract)
    │   ├── pages/            # Route pages (Login, Dashboard)
    │   ├── store/            # Pinia stores (auth, wallet)
    │   └── types/            # TypeScript interfaces
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
- [ ] Initialize Vue 3 + Vite project with TailwindCSS
- [ ] Setup viem + web3modal for wallet connection
- [ ] Add Pinia stores for auth and wallet state
- [ ] Create Login page (wallet connect + signature)
- [ ] Create Dashboard page
  - Display balance
  - Deposit form
  - Withdraw form
  - Transfer form
- [ ] Add transaction history with pagination
- [ ] Connect to Go backend API

### Frontend Stack
| Package | Purpose |
|---------|---------|
| Vue 3 + Vite | Framework & build |
| TypeScript | Type safety |
| Pinia | State management |
| Vue Router | Navigation |
| viem | Blockchain interaction |
| web3modal | Wallet connection |
| Axios | HTTP client |
| TailwindCSS | Styling |

### Phase 3: Backend (✅ COMPLETE)
- [x] Setup Go project with Gin
- [x] Setup PostgreSQL connection
- [x] Create migrations
- [x] Implement repositories
- [x] Implement usecases
- [x] Implement HTTP handlers
- [x] Implement event listener (sync Deposit/Withdrawal/Transfer events to DB)
- [x] Endpoints:
  - `POST /api/v1/auth/nonce` - Get nonce for wallet
  - `POST /api/v1/auth/login` - Verify signature
  - `GET /api/v1/account/:address` - Get balance
  - `GET /api/v1/transactions?page=1&limit=20` - Paginated history
  - `GET /api/v1/contract/info` - Token info

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

## Auto-Commit Skill

**CRITICAL**: After completing ANY feature, change, fix, or code modification, the agent MUST:

1. **Detect changes** → `git status`
2. **Create feature branch** → `git checkout -b <type>/<name>`
3. **Stage files** (atomic - only related files)
4. **Commit** → `git commit -m "<type>(<scope>): <description>"`
5. **Push** → `git push -u origin <branch-name>`
6. **Create PR** → `gh pr create`
7. **Merge** → `gh pr merge --squash --delete-branch`

### Branch & Commit Format
| Scenario | Branch | Commit |
|----------|--------|--------|
| New feature | `feat/<name>` | `feat(contracts): add feature` |
| Bug fix | `fix/<name>` | `fix(frontend): resolve bug` |
| Test | `test/<name>` | `test(contracts): add tests` |
| Config | `chore/<name>` | `chore(config): update settings` |

### Skill Files
- `.skills/auto-commit.md` - Main workflow
- `.skills/commit-types.md` - Commit type reference
- `.skills/conventions.md` - Project conventions

**Rules**:
- Always commit (no override)
- One logical change per commit
- Push after every commit
- Squash merge to main

---

## Notes for Future Agent

- Hardhat Node must be running before deploying contracts or testing: `npx hardhat node`
- Frontend uses viem for blockchain interaction, NOT ethers.js
- Backend user implements themselves - agent only assists
- All blockchain operations are on-chain (no off-chain ledger)
- Transaction history pagination uses DB indexing, not direct contract queries
- **ALWAYS use auto-commit skill after every feature**

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

# Frontend
cd Frontend
npm install
npm run dev

# Backend
cd Backend
cp .env.example .env
go mod download
go run cmd/api/main.go
```
