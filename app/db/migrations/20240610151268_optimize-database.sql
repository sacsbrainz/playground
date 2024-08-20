-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied.

-- Apply PRAGMA settings that can be set inside a transaction
PRAGMA auto_vacuum = incremental;

PRAGMA page_size = 4096;

PRAGMA busy_timeout = 5000;

PRAGMA cache_size = -2048;

PRAGMA foreign_keys = ON;

PRAGMA incremental_vacuum = 1000;

PRAGMA temp_store = MEMORY;

-- Apply PRAGMA settings that must be set outside of a transaction
-- These will be run individually outside of the transaction block
COMMIT;

PRAGMA journal_mode = WAL;

PRAGMA synchronous = NORMAL;

PRAGMA mmap_size = 524288000;

BEGIN;

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back.

-- Revert PRAGMA settings that can be set inside a transaction
PRAGMA auto_vacuum = NONE;

PRAGMA page_size = 4096;

PRAGMA busy_timeout = 0;

PRAGMA cache_size = 0;

PRAGMA foreign_keys = OFF;

PRAGMA incremental_vacuum = 0;

PRAGMA temp_store = DEFAULT;

-- Revert PRAGMA settings that must be set outside of a transaction
-- These will be run individually outside of the transaction block
COMMIT;

-- PRAGMA journal_mode = DELETE;

PRAGMA synchronous = FULL;

PRAGMA mmap_size = 0;

BEGIN;