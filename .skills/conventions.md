# Project Conventions

## Git Workflow

### Branch Strategy
- Each feature/fix = new branch from `main`
- Branch naming: `<type>/<kebab-case-name>`
- Types: `feat`, `fix`, `test`, `chore`, `refactor`, `docs`

### Commit Format
```
<type>(<scope>): <description>

[optional body]
[optional footer]
```

### Atomic Commits
- One logical change per commit
- Do not mix unrelated changes
- Example of WRONG:
  ```
  feat: add login and update dependencies
  ```
- Example of CORRECT:
  ```
  feat(frontend): add wallet login component
  chore(deps): upgrade pinia to latest version
  ```

## Code Conventions

### Smart Contracts (Solidity)
- Solidity version: `^0.8.28`
- Use OpenZeppelin contracts for standards
- Include NatSpec comments for public functions
- Events for all state changes
- Use `require` for validation with descriptive messages

### Frontend (Vue 3)
- Composition API with `<script setup>`
- TypeScript strict mode
- Pinia for state management
- viem for blockchain interaction (NOT ethers.js)

### Backend (Go)
- Gin framework for HTTP
- Clean architecture: domain → usecase → repository → delivery
- PostgreSQL for data persistence

## File Organization

```
Web3-Banking/
├── Blockchain/              # Hardhat project
│   ├── contracts/           # Solidity contracts
│   ├── scripts/             # Deployment scripts
│   ├── test/                # Test files
│   └── hardhat.config.ts    # Hardhat config
├── Frontend/                # Vue 3 + Vite
│   ├── src/
│   │   ├── components/      # Vue components
│   │   ├── pages/           # Route pages
│   │   ├── composables/     # Vue composables
│   │   ├── store/           # Pinia stores
│   │   └── api/             # API clients
│   └── vite.config.ts
├── Backend/                 # Go + Gin
│   ├── cmd/api/             # Entry point
│   ├── internal/
│   │   ├── domain/          # Entities
│   │   ├── usecase/         # Business logic
│   │   ├── repository/      # Data access
│   │   └── delivery/        # HTTP handlers
│   └── migrations/          # DB migrations
└── .skills/                 # Agent skills
```

## GitHub Integration
- Use `gh` CLI for all GitHub operations
- Authenticate with: `"/c/Program Files/GitHub CLI/gh.exe"`
- Create PR for every feature branch
- Squash merge to main

## Blockchain Configuration
- Network: Hardhat Node (local)
- Chain ID: 31337
- RPC: http://127.0.0.1:8545
