function install_schema() {
  if [[ -x `which psql` ]] ; then 
    echo "installing schema '${1}'"
    echo "create schema ${1} ; set search_path to ${1} ; " | cat - *.sql | psql -h localhost -U pflow
  else
    echo "psql client not installed"
  fi
}

install_schema hardhat
