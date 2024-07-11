# Migration Guide for TigerBeagle

![TigerBeagle Logo](../assets/logo.png)


This document outlines the process for migrating accounts and transactions using the TigerBeagle tool. It includes the JSON format for input files and the CLI commands to perform the migrations.

## Migrating Accounts

### JSON Format for Accounts

To migrate accounts, prepare a JSON file with an array of account objects. Each account object should have the following structure:

```json
[
  {
    "id": 123456789,
    "user_id": 0,
    "ledger": 700,
    "code": 10,
    "flags": 0,
    "debits_pending": 0,
    "debits_posted": 0,
    "credits_pending": 0,
    "credits_posted": 0
  },
  {
    "id": 987654321,
    "user_id": 0,
    "ledger": 700,
    "code": 10,
    "flags": 0,
    "debits_pending": 0,
    "debits_posted": 0,
    "credits_pending": 0,
    "credits_posted": 0
  }
]
```

Note:
- `id` should be a unique identifier for each account.
- `debitsPending`, `debitsPosted`, `creditsPending`, `creditsPosted`, and `userData128` are strings representing uint128 values.
- `userData64` and `userData32` are integer values.
- `ledger` is typically set to 1 unless you're using multiple ledgers.
- `code` is a user-defined value, often used to categorize accounts.
- `flags` is a 16-bit integer where each bit represents a boolean flag. For example, 1 represents `DebitsMustNotExceedCredits`.

### CLI Command for Account Migration

To migrate accounts from a JSON file, use the following command:

```
tigerbeagle migrate-accounts <path_to_json_file>
```

For example:

```
tigerbeagle migrate-accounts ./accounts_to_migrate.json
```

## Migrating Transfers

### JSON Format for Transfers

To migrate transfers, prepare a JSON file with an array of transfer objects. Each transfer object should have the following structure:

```json
[
  {
    "id": 1234567890,
    "debit_account_id": 123456789,
    "credit_account_id": 987654321,
    "amount": 100000,
    "pending_id": 0,
    "user_data_128": 0,
    "user_data_64": 0,
    "user_data_32": 0,
    "timeout": 0,
    "ledger": 700,
    "code": 10,
    "flags": 0
  },
  {
    "id": 9876543210,
    "debit_account_id": 987654321,
    "credit_account_id": 123456789,
    "amount": 50000,
    "pending_id": 0,
    "user_data_128": 0,
    "user_data_64": 0,
    "user_data_32": 0,
    "timeout": 0,
    "ledger": 700,
    "code": 10,
    "flags": 0
  }
]
```

Note:
- `id`, `debitAccountId`, `creditAccountId`, `amount`, `pendingId`, and `userData128` are strings representing uint128 values.
- `userData64` and `userData32` are integer values.
- `timeout` is an unsigned 32-bit integer representing the timestamp after which a two-phase transfer is considered timed out.
- `ledger` is typically set to 1 unless you're using multiple ledgers.
- `code` is a user-defined value, often used to categorize transfers.
- `flags` is a 16-bit integer where each bit represents a boolean flag.

### CLI Command for Transfer Migration

To migrate transfers from a JSON file, use the following command:

```
tigerbeagle migrate-transfers <path_to_json_file>
```

For example:

```
tigerbeagle migrate-transfers ./transfers_to_migrate.json
```

## Additional Notes

1. Ensure that the TigerBeetle server is running and accessible before starting the migration process.

2. It's recommended to migrate accounts before migrating transfers to ensure that all necessary accounts exist in the system.

3. The TigerBeagle tool will provide feedback on the migration process, including any errors or warnings for individual accounts or transfers. If an error occurs during migration, the process will stop and report which batch encountered the error. Successfully migrated batches before the error will remain in the system. You may need to restart the process with a fresh TB database.

4. For large datasets, tigerbeagle will handle the batching so you don't need to worry about the maximum batch size.

5. Always test the migration process in a non-production environment before applying it to a production system.

6. Make sure to keep backups of your data before starting the migration process.

For more information on using the TigerBeagle tool, refer to the main README file or use the `--help` flag with any command.
