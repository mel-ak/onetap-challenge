-- Drop indexes
DROP INDEX IF EXISTS idx_bills_status;
DROP INDEX IF EXISTS idx_bills_due_date;
DROP INDEX IF EXISTS idx_bills_provider_id;
DROP INDEX IF EXISTS idx_bills_linked_account_id;
DROP INDEX IF EXISTS idx_linked_accounts_provider_id;
DROP INDEX IF EXISTS idx_linked_accounts_user_id;

-- Drop tables
DROP TABLE IF EXISTS bills;
DROP TABLE IF EXISTS linked_accounts;
DROP TABLE IF EXISTS providers;
DROP TABLE IF EXISTS users;