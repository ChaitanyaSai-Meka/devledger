# DevLedger

A local-first CLI tool for tracking and splitting shared infrastructure costs among developer teams. No cloud, no setup - runs entirely on your machine.

## Installation

### Option 1 - Install directly
```bash
go install github.com/ChaitanyaSai-Meka/devledger@latest
```

### Option 2 - Build from source
```bash
git clone https://github.com/ChaitanyaSai-Meka/devledger
cd devledger
go build -o devledger .
sudo mv devledger /usr/local/bin/
```

## Quick Start
```bash
# Create users
devledger user create --username "alice"
devledger user create --username "bob"
devledger user create --username "charlie"

# Create a group and add members
devledger group create --groupname "backend-team"
devledger group add-member --groupname "backend-team" --username "alice"
devledger group add-member --groupname "backend-team" --username "bob"
devledger group add-member --groupname "backend-team" --username "charlie"

# Add an expense
devledger expense add --groupname "backend-team" --description "AWS Bill" --amount "30.00" --paidby "alice"

# View balances
devledger balance group --groupname "backend-team"

# Get suggested settlements
devledger balance simplify --groupname "backend-team"
```

## Commands

### Users
```bash
devledger user create --username <name>
devledger user delete --username <name>
devledger user list
devledger user groups --username <name>
```

### Groups
```bash
devledger group create --groupname <name>
devledger group delete --groupname <name>
devledger group list
devledger group members --groupname <name>
devledger group add-member --groupname <name> --username <name>
devledger group remove-member --groupname <name> --username <name>
```

### Expenses
```bash
devledger expense add --groupname <name> --description <desc> --amount <amount> --paidby <name>
devledger expense list --groupname <name>
devledger expense list-by-user --username <name>
devledger expense detail --expenseid <id>
devledger expense delete --expenseid <id>
devledger expense settle --expenseid <id> --username <name>
devledger expense unsettled --username <name>
```

### Balances
```bash
devledger balance group --groupname <name>
devledger balance simplify --groupname <name>
```

## How it Works
```
CLI commands
    ↓ HTTP calls
REST API (localhost:38080)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (DB queries)
    ↓
SQLite (devledger.db)
```

The server starts automatically on the first command and shuts down when done.

## Amount Format

All amounts are entered in rupees with up to 2 decimal places:
```bash
--amount "30.00"   # ₹30.00
--amount "1500"    # ₹1500.00
--amount "99.99"   # ₹99.99
```

Internally stored in paise to avoid floating-point precision issues.

## Future Enhancements

- Custom expense splits (unequal amounts)
- Partial payment tracking
- Edit expense description and amount
- Database migrations using golang-migrate
- Idle timeout - auto shutdown server after inactivity
- Multi-currency support

## Tech Stack

- **Language** - Go
- **Router** - Chi
- **CLI** - Cobra
- **Database** - SQLite (modernc.org/sqlite - pure Go, no CGo)
- **Money parsing** - shopspring/decimal
