# Commit Types Reference

## Types

### feat
New feature implementation.
```
feat(scope): add user registration functionality
feat(contracts): implement multi-signature support
feat(frontend): add transaction history component
```

### fix
Bug fix.
```
fix(contracts): resolve reentrancy vulnerability in withdraw
fix(frontend): fix balance not updating after deposit
fix(api): handle null response from blockchain node
```

### test
Adding or updating tests.
```
test(contracts): add unit tests for transfer function
test(frontend): add component tests for login form
test(contracts): add integration tests for Bank contract
```

### chore
Maintenance tasks, dependencies, build, config.
```
chore(config): update hardhat version to 2.22.0
chore(deps): upgrade OpenZeppelin contracts
chore(scripts): add deployment automation
chore(ci): add GitHub Actions workflow
```

### docs
Documentation only changes.
```
docs(readme): add deployment instructions
docs(contracts): document Bank contract interface
docs(api): add endpoint documentation
```

### refactor
Code restructuring without behavior change.
```
refactor(contracts): extract validation logic
refactor(frontend): optimize component structure
refactor(api): improve error handling pattern
```

### perf
Performance improvements.
```
perf(contracts): optimize gas usage in transfer
perf(frontend): implement virtual scrolling for large lists
```

### style
Code style changes (formatting, no logic change).
```
style(contracts): fix linting issues
style(frontend): standardize code formatting
```

## Scope Reference

| Scope | Usage |
|-------|-------|
| `contracts` | Smart contract code (`.sol` files) |
| `scripts` | Deployment and utility scripts |
| `config` | Configuration files |
| `frontend` | Vue/TypeScript frontend code |
| `backend` | Go backend code |
| `api` | HTTP API endpoints |
| `database` | Database migrations, schemas |
| `root` | Project root files (.gitignore, etc.) |

## Rules
1. Use imperative mood: "add" not "added"
2. Keep subject line under 72 characters
3. Use lowercase for type and scope
4. No period at end of message
5. Reference issues/tickets in body if needed
