-- Change LTF and HTF columns from VARCHAR(10) to TEXT to store URLs
ALTER TABLE trading_journal_entries
    ALTER COLUMN ltf TYPE TEXT,
    ALTER COLUMN htf TYPE TEXT;
