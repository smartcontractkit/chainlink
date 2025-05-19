-- +goose Up
-- BFC-2694 - fix search path so public takes precedence. No need for a downward migration.
SET search_path TO public,evm;

