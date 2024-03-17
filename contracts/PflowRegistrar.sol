// SPDX-License-Identifier: MIT
pragma solidity ^0.8.18;

import "./Metamodel.sol";

interface PflowInterface {
    function register(address _address) external;

    function getFlows() external view returns (PflowEth.FlowInfo[] memory);
}

abstract contract PflowEth  {

    event Registered(
        bytes32 flowId,
        address indexed model,
        uint256 indexed id
    );

    address[] public flows;
    mapping(address => uint256) public routeId;
    mapping(bytes32 => address) public routes;

    struct FlowInfo {
        address model;
        bytes32 flowId;
        uint256 id;
    }

    function calculateFlowId(address _address, uint256 id) external view returns (bytes32) {
        return keccak256(abi.encodePacked(address(this), "#", _address, "#", id));
    }

    function addFlow(bytes32 flowId, address _address) internal returns (bytes32) {
        emit Registered(flowId, _address, flows.length);
        routeId[_address] = flows.length;
        flows.push(_address);
        return flowId;
    }

    function setAddress(address _address) internal {
        Declaration.PetriNet memory m = ModelInterface(_address).declaration(); // test call
        require(m.places.length > 0 || m.transitions.length > 0, "ModelRegistry: invalid model");
        require(_address != address(0), "ModelRegistry: invalid address");
        routes[addFlow(this.calculateFlowId(_address, flows.length), _address)] = _address;
    }
}

contract PflowRegistrar is PflowEth, PflowInterface {

    function register(address _address) external {
        require(routeId[_address] <= 0, "address already registered");
        setAddress(_address);
    }

    function registerMany(address[] calldata _addresses) external {
        for (uint256 i = 0; i < _addresses.length; i++) {
            setAddress(_addresses[i]);
        }
    }

    function getFlows() external view returns (FlowInfo[] memory) {
        FlowInfo[] memory result = new FlowInfo[](flows.length);
        for (uint256 i = 0; i < flows.length; i++) {
            result[i] = FlowInfo(flows[i], this.calculateFlowId(flows[i], i), i);
        }
        return result;
    }

}
