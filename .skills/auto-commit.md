# Auto-Commit Skill

## Purpose
Automatically commits every feature/change using atomic commits with conventional format.

## Trigger
**After completing ANY feature, change, fix, or code modification.**

## Workflow

### 1. Detect Changes
```bash
git status
```
Identify all modified, added, or deleted files.

### 2. Analyze Changes
Determine:
- **Type**: feat | fix | test | chore | docs | refactor
- **Scope**: contracts | scripts | frontend | backend | config | root
- **Description**: Brief action/change description

### 3. Create Feature Branch
```bash
git checkout -b <type>/<feature-name>
```
- Features: `feat/<name>`
- Fixes: `fix/<name>`
- Refactor: `refactor/<name>`

### 4. Stage Files (Atomic)
Stage only files related to this specific change:
```bash
git add <specific-files>
```
Do NOT stage unrelated changes.

### 5. Commit
```bash
git commit -m "<type>(<scope>): <description>"
```

### 6. Push
```bash
git push -u origin <branch-name>
```

### 7. Create & Merge PR
```bash
gh pr create --title "<type>: <description>" --body "Auto-committed feature"
gh pr merge --squash --delete-branch
```

## Branch Naming
| Type | Branch Pattern | Example |
|------|----------------|---------|
| feat | `feat/<name>` | `feat/user-authentication` |
| fix | `fix/<name>` | `fix/balance-calculation` |
| refactor | `refactor/<name>` | `refactor/contract-storage` |
| test | `test/<name>` | `test/withdrawal-flow` |
| chore | `chore/<name>` | `chore/update-dependencies` |

## Auto-Generated Commit Messages
Format: `<type>(<scope>): <imperative-description>`

Examples:
- `feat(contracts): add withdrawal functionality`
- `fix(frontend): resolve balance display error`
- `test(contracts): add deposit unit tests`
- `chore(config): update hardhat network settings`
- `docs(readme): add deployment instructions`

## Scope Mapping
| Folder/Component | Scope |
|-----------------|-------|
| `Blockchain/contracts/` | `contracts` |
| `Blockchain/scripts/` | `scripts` |
| `Blockchain/test/` | `contracts` (tests) |
| `Blockchain/` config files | `config` |
| `Frontend/` | `frontend` |
| `Backend/` | `backend` |
| Root level | `root` |

## Notes
- **Always commit** - No exceptions, no override
- **Atomic commits** - One logical change per commit
- **Push immediately** - After each commit
- **Merge via PR** - Create PR and squash merge
