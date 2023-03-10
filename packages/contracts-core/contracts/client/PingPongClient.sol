// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

import { IMessageRecipient } from "../interfaces/IMessageRecipient.sol";
import { IOrigin } from "../interfaces/IOrigin.sol";
import { Tips } from "../libs/Tips.sol";
import { TypeCasts } from "../libs/TypeCasts.sol";

contract PingPongClient is IMessageRecipient {
    using TypeCasts for address;

    struct PingPongMessage {
        uint256 pingId;
        bool isPing;
        uint16 counter;
    }

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                              IMMUTABLES                              ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    // local chain Origin: used for sending messages
    address public immutable origin;

    // local chain Destination: used for receiving messages
    address public immutable destination;

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                               STORAGE                                ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    uint256 public random;

    /// @notice Amount of "Ping" messages sent.
    uint256 public pingsSent;

    /// @notice Amount of "Ping" messages received.
    /// Every received Ping message leads to sending a Pong message back to initial sender.
    uint256 public pingsReceived;

    /// @notice Amount of "Pong" messages received.
    /// When all messages are delivered, should be equal to `pingsSent`
    uint256 public pongsReceived;

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                                EVENTS                                ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    /// @notice Emitted when a Ping message is sent.
    /// Triggered externally, or by receveing a Pong message with instructions to do more pings.
    event PingSent(uint256 pingId);

    /// @notice Emitted when a Ping message is received.
    /// Will always send a Pong message back.
    event PingReceived(uint256 pingId);

    /// @notice Emitted when a Pong message is sent.
    /// Triggered whenever a Ping message is received.
    event PongSent(uint256 pingId);

    /// @notice Emitted when a Pong message is received.
    /// Will initiate a new Ping, if the counter in the message is non-zero.
    event PongReceived(uint256 pingId);

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                             CONSTRUCTOR                              ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    constructor(address _origin, address _destination) {
        origin = _origin;
        destination = _destination;
        // Initiate "random" value
        random = uint256(keccak256(abi.encode(block.number)));
    }

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                          RECEIVING MESSAGES                          ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    /// @notice Called by Destination upon executing the message.
    function handle(
        uint32 _origin,
        uint32,
        bytes32 _sender,
        uint256,
        bytes memory _message
    ) external {
        require(msg.sender == destination, "PingPongClient: !destination");
        PingPongMessage memory _msg = abi.decode(_message, (PingPongMessage));
        if (_msg.isPing) {
            // Ping is received
            ++pingsReceived;
            emit PingReceived(_msg.pingId);
            // Send Pong back
            _pong(_origin, _sender, _msg);
        } else {
            // Pong is received
            ++pongsReceived;
            emit PongReceived(_msg.pingId);
            // Send extra ping, if initially requested
            if (_msg.counter != 0) {
                _ping(_origin, _sender, _msg.counter - 1);
            }
        }
    }

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                           SENDING MESSAGES                           ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    function doPings(
        uint16 _pingCount,
        uint32 _destination,
        address _recipient,
        uint16 _counter
    ) external {
        for (uint256 i = 0; i < _pingCount; ++i) {
            _ping(_destination, _recipient.addressToBytes32(), _counter);
        }
    }

    /// @notice Send a Ping message to destination chain.
    /// Upon receiving a Ping, a Pong message will be sent back.
    /// If `_counter > 0`, this process will be repeated when the Pong message is received.
    /// @param _destination Chain to send Ping message to
    /// @param _recipient   Recipient of Ping message
    /// @param _counter     Additional amount of Ping-Pong rounds to conclude
    function doPing(
        uint32 _destination,
        address _recipient,
        uint16 _counter
    ) external {
        _ping(_destination, _recipient.addressToBytes32(), _counter);
    }

    function nextOptimisticPeriod() public view returns (uint32 period) {
        // Use random optimistic period up to one minute
        return uint32(random % 1 minutes);
    }

    /*╔══════════════════════════════════════════════════════════════════════╗*\
    ▏*║                            INTERNAL LOGIC                            ║*▕
    \*╚══════════════════════════════════════════════════════════════════════╝*/

    /// @dev Returns a random optimistic period value from 0 to 59 seconds.
    function _optimisticPeriod() internal returns (uint32 period) {
        // Use random optimistic period up to one minute
        period = nextOptimisticPeriod();
        // Adjust "random" value
        random = uint256(keccak256(abi.encode(random)));
    }

    /**
     * @dev Send a "Ping" or "Pong" message.
     * @param _destination  Domain of destination chain
     * @param _recipient    Message recipient on destination chain
     * @param _msg          Ping-pong message
     */
    function _sendMessage(
        uint32 _destination,
        bytes32 _recipient,
        PingPongMessage memory _msg
    ) internal {
        bytes memory tips = Tips.emptyTips();
        bytes memory message = abi.encode(_msg);
        IOrigin(origin).dispatch(_destination, _recipient, _optimisticPeriod(), tips, message);
    }

    /// @dev Initiate a new Ping-Pong round.
    function _ping(
        uint32 _destination,
        bytes32 _recipient,
        uint16 _counter
    ) internal {
        uint256 pingId = pingsSent++;
        _sendMessage(
            _destination,
            _recipient,
            PingPongMessage({ pingId: pingId, isPing: true, counter: _counter })
        );
        emit PingSent(pingId);
    }

    /// @dev Send a Pong message back.
    function _pong(
        uint32 _destination,
        bytes32 _recipient,
        PingPongMessage memory _msg
    ) internal {
        _sendMessage(
            _destination,
            _recipient,
            PingPongMessage({ pingId: _msg.pingId, isPing: false, counter: _msg.counter })
        );
        emit PongSent(_msg.pingId);
    }
}