const {
    time,
    loadFixture,
} = require("@nomicfoundation/hardhat-network-helpers");
require("@nomiclabs/hardhat-web3");
const {ethers} = require("hardhat");
// require chai
const {expect} = require("chai");


describe("KonamiCode", function () {
    async function deployTestProxy() {
        // Contracts are deployed using the first signer/account by default
        const [owner, p0, p1] = await ethers.getSigners();

        const contractFactory = await ethers.getContractFactory("contracts/KonamiCode.sol:KonamiCode", {
            libraries: {},
        });
        const contract = await contractFactory.deploy(); // REVIEW: swap to test access ctl
        // console.log(api.from);

        return {contract, p0, p1};
    }

    it("should accept the konami-code sequence", async function () {
        const { contract, p0} = await loadFixture(deployTestProxy);
        const k = await contract.connect(p0); // FYI K is the shorthand symbol for contract

        async function getState() {
            return {
                TwoUps: await k.state(0),
                TwoDowns: await k.state(1),
                TwoLefts: await k.state(2),
                TwoRights: await k.state(3),
                ThenRight: await k.state(4),
                ThenSelect: await k.state(5),
                ThenStart: await k.state(6),
                ThenA: await k.state(7),
            }
        }

        let seq = 0;
        // contract data in hex
        const SESSION1 = '0x0000000000000000000000000000000000000000000000000000000000000001';

        async function call(func) {
            seq++;
            return func()
                .then(async (tx) => {
                    const receipt = await tx.wait();
                    expect(receipt.events[0].topics[1]).to.equal(SESSION1);
                    return getState();
                })
                .then((state) => {
                    console.log({seq, state})
                    return state;
                });
        }

        console.log({seq, state: await getState()});
        //await call(k.Up); // FIXME: contract should allow for an odd number of "Ups"
        await call(k.Up);
        await call(k.Up); // REVIEW: causes session to reset
        await call(k.Down)
        await call(k.Down);
        await call(k.Left);
        await call(k.Right);
        await call(k.Left);
        await call(k.Right);
        await call(k.B);
        await call(k.A);
        await call(k.Select); // REVIEW: this is the 2-player version
        await call(k.Start);
    });

});
