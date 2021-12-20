// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.8.4 <0.9.0;

import {BytesLib} from "./bytes.sol";
import {BLS} from "./bls.sol";
import {BN256G2} from "./bn256g2.sol";

contract MixinProcess {
    using BytesLib for bytes;
    using BLS for uint256[4];
    using BLS for uint256[2];
    using BLS for bytes;
    using BLS for uint256;
    using BN256G2 for uint256;

    struct Decimal {
        uint256 value;
        int32 exp;
        bool sign;
    }

    struct PriceData {
        uint64 timestamp;
        Decimal price;
    }

    event PriceEvent(uint128 asset, uint64 timestamp, PriceData price);

    uint64 GENESIS_TS = 1639500000;
    uint32 public THRESHOLD = 4;
    uint256[4][] public ORACLE_GROUP = [
        [
            // 135ca425d4051e4009040a03cd231212286fcb89a750998ed75683c56a418490016fece0d474b23b14ee82f2e04506ab385df5e757189ed05de50059fd2ac7b309056867016387e4d8380ea59723188b6dc8c3a1edaab2d167dc196b2a64d11829568534001f61090d2f5f1a63106478beee9ef8382f9c813db67417bb1e3bab
            0x016fece0d474b23b14ee82f2e04506ab385df5e757189ed05de50059fd2ac7b3, // x.y
            0x135ca425d4051e4009040a03cd231212286fcb89a750998ed75683c56a418490, // x.x
            0x29568534001f61090d2f5f1a63106478beee9ef8382f9c813db67417bb1e3bab, // y.y
            0x09056867016387e4d8380ea59723188b6dc8c3a1edaab2d167dc196b2a64d118 // y.x
        ],
        [
            // 22abab734a3883f004b6e48794317bfcba954fab493c4504e46c927b2cbb1a0907ed31630479243dda29b94799194c5a1436cc61f527ea0f77531426892a7af71f6dc0c459dc537ee40d16fbbc555d01cb3e95d658418640e0689952e282b24810303d9e4fc88902543ff5baa120fd0a2e59a23528961e7e582c3c7284cc73d8
            0x07ed31630479243dda29b94799194c5a1436cc61f527ea0f77531426892a7af7,
            0x22abab734a3883f004b6e48794317bfcba954fab493c4504e46c927b2cbb1a09,
            0x10303d9e4fc88902543ff5baa120fd0a2e59a23528961e7e582c3c7284cc73d8,
            0x1f6dc0c459dc537ee40d16fbbc555d01cb3e95d658418640e0689952e282b248
        ],
        [
            // 180da0ffabd02f4df7001a6c7732a692cd4306a77bc98651bb6672ef976f5213023926a33c3da7868f9cf12e1b0b61bca58e6452c7f77f09d62250a3ffbc790714d560e8b11790caee9650dfd5399c8954dc65cfdb157764df2d9727bb0a6cfa11c574c61c40b428bb2bfe5451ffea378e060f619f4dd0e5d2dc42c86961a92a
            0x023926a33c3da7868f9cf12e1b0b61bca58e6452c7f77f09d62250a3ffbc7907,
            0x180da0ffabd02f4df7001a6c7732a692cd4306a77bc98651bb6672ef976f5213,
            0x11c574c61c40b428bb2bfe5451ffea378e060f619f4dd0e5d2dc42c86961a92a,
            0x14d560e8b11790caee9650dfd5399c8954dc65cfdb157764df2d9727bb0a6cfa
        ],
        [
            // 1a9d8c29bb6a9a18e65bcc5509382a18ae4ca3faf1b66b068047d4fdc35974da0616da23b96a256699e0517817b1f52d1291f0537e9824709343a6389ca239b420b73184cc5a91efac874e5542aee3e0ec5ec58feaba929339b414bc1c68b32a216d106575ce6b47d742249cff7aa0860b56c71d941b17704dc121228fb7f9f6
            0x0616da23b96a256699e0517817b1f52d1291f0537e9824709343a6389ca239b4,
            0x1a9d8c29bb6a9a18e65bcc5509382a18ae4ca3faf1b66b068047d4fdc35974da,
            0x216d106575ce6b47d742249cff7aa0860b56c71d941b17704dc121228fb7f9f6,
            0x20b73184cc5a91efac874e5542aee3e0ec5ec58feaba929339b414bc1c68b32a
        ],
        [
            // 2724f455ae874c8d4e32c7ce6b20e3365f4ea38edb4c93488ad0a9d90d3a3b450519a0e7b748ab056072966a63b76be42af513fb29e8754d46415cbe3235e8620075b07147db40b927904c07e8fec7e3e17b2a82c03291d63ed0e726a01e91052115f68d2ab8f29193dc9f99f14f40647485a240e685255a0598d6ce54eab18f
            0x0519a0e7b748ab056072966a63b76be42af513fb29e8754d46415cbe3235e862,
            0x2724f455ae874c8d4e32c7ce6b20e3365f4ea38edb4c93488ad0a9d90d3a3b45,
            0x2115f68d2ab8f29193dc9f99f14f40647485a240e685255a0598d6ce54eab18f,
            0x0075b07147db40b927904c07e8fec7e3e17b2a82c03291d63ed0e726a01e9105
        ],
        [
            // 04d9352637cf7f98d2b40fd3912bda1df94065c46fa2d4933a444d7120cc644f24bc963073ba6ea6decda4c80b4086ff041507234db65927927e50cba56657451bd8bf1bb42c940a8064c0ea2cffaaf538542316032ba262b1f950e173f258ae2a2ecf6fa0a69ec567af3135c875c43692a6f4c3ae10dd192873f6eb3936bfa6
            0x24bc963073ba6ea6decda4c80b4086ff041507234db65927927e50cba5665745,
            0x04d9352637cf7f98d2b40fd3912bda1df94065c46fa2d4933a444d7120cc644f,
            0x2a2ecf6fa0a69ec567af3135c875c43692a6f4c3ae10dd192873f6eb3936bfa6,
            0x1bd8bf1bb42c940a8064c0ea2cffaaf538542316032ba262b1f950e173f258ae
        ]
    ];

    mapping(uint128 => PriceData) public prices;

    function work(bytes memory data) internal returns (bool) {
        uint256 offset = 0;

        require(
            data.length >= 65 || data.length <= 162,
            "invalid extra length"
        );

        uint8 size = data.toUint8(offset);
        require(size < 8, "invalid timestamp size");
        offset += 1;

        uint64 timestamp = uint64(toInt64(data, offset, size));

        require(timestamp > GENESIS_TS, "invalid timestamp");
        offset += size;

        require(data.toUint8(offset) == 16, "invalid asset");
        offset += 1;

        uint128 asset = data.toUint128(offset);
        offset += 16;

        require(
            prices[asset].timestamp < timestamp,
            "timestamp older than last price"
        );

        size = data.toUint8(offset);
        require(size >= 4 && size <= 37, "invalid price");
        offset += 1;

        PriceData memory price;
        price.timestamp = timestamp;
        (price.price.sign, price.price.exp, price.price.value) = toDecimal(
            data,
            offset,
            size
        );
        offset += size;

        uint256[2] memory message = data.slice(0, offset).hashToPoint();

        size = data.toUint8(offset);
        require(size == 36 || size == 97, "invalid cosi-signature");
        offset += 1;

        size = data.toUint8(offset);
        require(size > 0 && size <= 8, "invalid signature mask size");
        offset += 1;

        uint64 mask = toUint64(data, offset, size);
        offset += 1;

        uint256[4] memory pubkey = maskToPublicKey(mask);

        size = data.toUint8(offset);
        require(size == 33 || size == 64, "invalid signature size");
        offset += 1;

        uint256[2] memory sig = toSignature(data, offset, size);
        offset += size;

        require(sig.verifySingle(pubkey, message), "invalid price signature");

        prices[asset] = price;
        emit PriceEvent(asset, timestamp, price);
        return true;
    }

    // process || nonce || asset || amount || extra || timestamp || members || threshold || sig
    function mixin(bytes calldata raw) public returns (bool) {
        require(raw.length >= 141, "event data too small");

        uint256 size = 0;
        uint256 offset = 40;
        size = raw.toUint16(offset);
        require(size <= 32, "integer out of bounds");
        offset = offset + 2 + size;

        size = raw.toUint16(offset);
        offset = offset + 2;
        bytes memory extra = raw.slice(offset, size);
        return work(extra);
    }

    function getPrice(uint128 asset) public view returns (PriceData memory) {
        return prices[asset];
    }

    function toInt64(
        bytes memory data,
        uint256 offset,
        uint8 size
    ) internal pure returns (int64) {
        uint64 ux = toUint64(data, offset, size);
        int64 x = int64(ux >> 1);
        if (ux & 1 != 0) {
            x = ~x;
        }
        return x;
    }

    function toUint64(
        bytes memory data,
        uint256 offset,
        uint8 size
    ) internal pure returns (uint64) {
        uint64 x;
        uint64 s;
        for (uint256 i = offset; i < offset + size; i++) {
            uint8 b = data.toUint8(i);
            if (b < 0x80) {
                require(i != 10 || b < 1, "uint64 value overflow");
                return x | (uint64(b) << s);
            }
            x = x | ((uint64(b & 0x7f)) << s);
            s += 7;
        }
        return 0;
    }

    function toDecimal(
        bytes memory data,
        uint256 offset,
        uint8 size
    )
        internal
        pure
        returns (
            bool,
            int32,
            uint256
        )
    {
        int32 exp = data.toInt32(offset);
        offset += 4;

        if (size == 4) {
            return (true, exp, 0);
        }

        bool sign = (data.toUint8(offset) & 1) != 0;
        offset += 1;

        uint256 value = new bytes(32 - (size - 5))
            .concat(data.slice(offset, size - 5))
            .toUint256(0);

        return (sign, exp, value);
    }

    function toSignature(
        bytes memory data,
        uint256 offset,
        uint8 size
    ) internal view returns (uint256[2] memory) {
        uint256[2] memory sig;
        sig[0] = data.toUint256(offset);
        offset += 32;
        if (size == 64) {
            sig[1] = data.toUint256(offset + 32);
            offset += 32;
        } else {
            sig[1] = sig[0].sigToUncompresed(data.toUint8(offset));
            offset += 1;
        }

        require(sig.isValidSignature(), "invalid signature");
        return sig;
    }

    function maskToPublicKey(uint64 mask)
        internal
        view
        returns (uint256[4] memory)
    {
        uint256 pubcount = 0;
        for (uint256 i = 0; i < ORACLE_GROUP.length; i++) {
            if ((mask & (1 << (i + 1))) != 0) {
                pubcount += 1;
            }
        }

        require(pubcount >= THRESHOLD, "cosi signature has not enough signers");

        uint256[4] memory pubkey;
        for (uint256 i = 0; i < ORACLE_GROUP.length; i++) {
            if ((mask & (1 << (i + 1))) == 0) {
                continue;
            }

            (pubkey[0], pubkey[1], pubkey[2], pubkey[3]) = pubkey[0].ecTwistAdd(
                pubkey[1],
                pubkey[2],
                pubkey[3],
                ORACLE_GROUP[i][0],
                ORACLE_GROUP[i][1],
                ORACLE_GROUP[i][2],
                ORACLE_GROUP[i][3]
            );
        }
        require(pubkey.isValidPublicKey(), "invalid public key");
        return pubkey;
    }
}
