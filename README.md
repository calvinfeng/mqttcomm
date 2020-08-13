# MQTT Communication Prototype

## Client ID

All clients are required to have a client name or ID. The client name is used by the MQTT broker to track subscriptions.
Client names must also be unique. If you attempt to connect to an MQTT broker with the same name as an existing client,
then the existing client connection is dropped. Because most MQTT clients will attempt to reconnect following a 
disconnect, this can result in a loop of disconnect and connect.

## Clean Sessions

MQTT clients by default establish a clean session with a broker. A clean session is one in which the broker isn't
expected to remember client subscriptions and may hold undelivered messages for the client.

## Last Will Messages

The idea of the last will message is to notify a subscriber that the publisher is unavailable due to network outage. The
last will message is set by the publishing client, and is set on a per topic basis which means that each topic can have
its own last will message. The message is stored on the broker and sent to any subscribing client (to that topic) if
the connection to the publisher fails.

## MQTT Topics

Topics are structured in a hierarchy similar to folders and files in a file system using forward slashes. Topic names are

- Case sensitive
- UTF-8 strings
- Must consist of at least one character to be valid

There are no topics created on a broker by default, except for the $SYS topic. All topics are created by a subscribing
or publishing client, and they are not permanent. A topic only exists if a client has subscribed to it, or a broker
has a retained or last will messages stored for that topic.

### $SYS

The `$SYS` is a reserved topic and is used by most MQTT brokers to publish information about the broker. They are
read-only topics for the MQTT clients.

## MQTT Security

http://www.steves-internet-guide.com/mqtt-security-mechanisms/

## Useful Examples

https://github.com/eclipse/paho.mqtt.golang/tree/master/cmd