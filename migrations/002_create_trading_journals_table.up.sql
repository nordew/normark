CREATE TABLE IF NOT EXISTS trading_journals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_trading_journals_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS trading_journal_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    journal_id UUID NOT NULL,
    day TIMESTAMP NOT NULL,
    asset VARCHAR(20) NOT NULL,
    ltf VARCHAR(10) NOT NULL,
    htf VARCHAR(10) NOT NULL,
    entry_charts TEXT[] DEFAULT '{}',
    session VARCHAR(20) NOT NULL,
    trade_type VARCHAR(20) NOT NULL,
    setup TEXT NULL,
    direction VARCHAR(10) NOT NULL,
    entry_type VARCHAR(10) NOT NULL,
    realized DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    max_rr DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    result VARCHAR(5) NOT NULL,
    notes TEXT DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT fk_trading_journal_entries_journal
        FOREIGN KEY (journal_id)
        REFERENCES trading_journals(id)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_trading_journals_user_id ON trading_journals(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journals_created_at ON trading_journals(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_trading_journals_deleted_at ON trading_journals(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_journal_id ON trading_journal_entries(journal_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_day ON trading_journal_entries(day DESC);
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_asset ON trading_journal_entries(asset) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_session ON trading_journal_entries(session) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_result ON trading_journal_entries(result) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_created_at ON trading_journal_entries(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_deleted_at ON trading_journal_entries(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_journal_day ON trading_journal_entries(journal_id, day DESC) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_trading_journal_entries_journal_asset ON trading_journal_entries(journal_id, asset) WHERE deleted_at IS NULL;

CREATE TRIGGER update_trading_journals_updated_at
    BEFORE UPDATE ON trading_journals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_trading_journal_entries_updated_at
    BEFORE UPDATE ON trading_journal_entries
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE trading_journal_entries
    ADD CONSTRAINT check_session CHECK (session IN ('asia', 'london', 'new_york'));

ALTER TABLE trading_journal_entries
    ADD CONSTRAINT check_trade_type CHECK (trade_type IN ('swing', 'intraday'));

ALTER TABLE trading_journal_entries
    ADD CONSTRAINT check_direction CHECK (direction IN ('buy', 'sell'));

ALTER TABLE trading_journal_entries
    ADD CONSTRAINT check_entry_type CHECK (entry_type IN ('market', 'limit'));

ALTER TABLE trading_journal_entries
    ADD CONSTRAINT check_result CHECK (result IN ('TP', 'SL', 'BE'));