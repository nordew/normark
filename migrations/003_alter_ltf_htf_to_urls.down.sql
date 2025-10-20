-- Revert LTF and HTF columns back to VARCHAR(10)
ALTER TABLE trading_journal_entries
    ALTER COLUMN ltf TYPE VARCHAR(10),
    ALTER COLUMN htf TYPE VARCHAR(10);
