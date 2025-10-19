DROP TRIGGER IF EXISTS update_trading_journal_entries_updated_at ON trading_journal_entries;
DROP TRIGGER IF EXISTS update_trading_journals_updated_at ON trading_journals;

DROP INDEX IF EXISTS idx_trading_journal_entries_journal_asset;
DROP INDEX IF EXISTS idx_trading_journal_entries_journal_day;
DROP INDEX IF EXISTS idx_trading_journal_entries_deleted_at;
DROP INDEX IF EXISTS idx_trading_journal_entries_created_at;
DROP INDEX IF EXISTS idx_trading_journal_entries_result;
DROP INDEX IF EXISTS idx_trading_journal_entries_session;
DROP INDEX IF EXISTS idx_trading_journal_entries_asset;
DROP INDEX IF EXISTS idx_trading_journal_entries_day;
DROP INDEX IF EXISTS idx_trading_journal_entries_journal_id;

DROP INDEX IF EXISTS idx_trading_journals_deleted_at;
DROP INDEX IF EXISTS idx_trading_journals_created_at;
DROP INDEX IF EXISTS idx_trading_journals_user_id;

DROP TABLE IF EXISTS trading_journal_entries;
DROP TABLE IF EXISTS trading_journals;
