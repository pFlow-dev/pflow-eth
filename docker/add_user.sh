#!/usr/bin/env bash
docker-compose exec db bash -c \
"echo \"\
  create user pflow with superuser login password 'pflow';
  create database pflow with owner pflow;
\" | psql -h localhost -U supabase_admin postgres"
