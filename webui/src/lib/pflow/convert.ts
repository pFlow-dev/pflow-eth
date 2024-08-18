export function uint256HexToNumeric(val: string): bigint {
    // Remove leading '0x' if present
    if (val.startsWith('0x')) {
        val = val.slice(2);
    }

    // Convert hex substrings to numeric values and sum them up
    const valA = BigInt('0x' + val.slice(0, 8)) * (2n ** 224n);
    const valB = BigInt('0x' + val.slice(8, 22)) * (2n ** 168n);
    const valC = BigInt('0x' + val.slice(22, 36)) * (2n ** 112n);
    const valD = BigInt('0x' + val.slice(36, 50)) * (2n ** 56n);
    const valE = BigInt('0x' + val.slice(50, 64));

    return valA + valB + valC + valD + valE;
}

