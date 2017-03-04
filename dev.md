## RNET Dev Notes

### Encryption, headers and packeting
I want to use the packeter for sending messages. I can leave the header
unencrypted. The only meta-data this exposes is the size of the message and the
message ID. Not a huge security flaw; that the hosts are communicating is
already known. Revealing the number of messages is not great, and potentially
another sender could muck things up by sending bad packets (since we're not
checking origin, but that can be faked anyway).

There's a bootstrap problem. I'd like all messages to be encrypted. But how does
someone initially send their key?

A --handshake--> B
A ==msg==> B
msg = packets
:.
A --P1--> B
A --P2--> B
...
A --Pn--> B

#### MitM
We shouldn't need to worry about MitM attacks. When A sends her key to B, she
MACs it with B's public key (she may already have it, it could be pulled from
the DHT or out-of-bands). An attacker would have to compromise A getting the
key. And because A can use a web of trust to validate the keys, an attacker
would need to execute a perfect MitM attack by simulating the entire network
from A's perspective. Not impossible, but difficult, particularly at a massive
scale.

#### Packet Headers
So we need a packet header. I think one byte would do.

0: Handshake packet
1: Message packet
2: Keep alive (should only be 1 byte, may not be necessary)
3: Ping
4: Pong

We don't have headers for hole punching because those will still be treated as
message packets

So internally, the workflow is

1) Get UDP message
2) Send to header-handler
3) Route to correct