function validate_and_install_schema() {
  schema_name="$1"
  if [[ "$schema_name" != "hardhat" && "$schema_name" != "sepolia_optimism" && "$schema_name" != "optimism" ]]; then
    echo "Invalid schema name. Please use 'hardhat', 'sepolia_optimism', or 'optimism'."
    exit 1
  fi
}

# NOTE use ./docker/add_user.sh to create db user 'pflow'
function install_schema() {
  if [[ -x `which psql` ]] ; then
    echo "installing schema '${1}'"
    echo "create schema ${1} ; set search_path to ${1} ; " | cat - *.sql | psql -h localhost -U pflow
  else
    echo "psql client not installed"
  fi
}

validate_and_install_schema "$1"
install_schema "$1"
