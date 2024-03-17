mkdir -p build
forge flatten ./test/examples/TicTacToe.sol  -o ./build/TicTacToe.sol
forge flatten ./contracts/MetaModel.sol  -o ./build/MetaModel.sol
forge flatten ./contracts/PflowRouter.sol  -o ./build/PflowRouter.sol
