#!/usr/bin/env bash
export PGPASSWORD="postgres"
psql -U postgres -d postgres -c "drop database testme"
psql -U postgres -d postgres -c "create database testme"
