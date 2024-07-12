# Socket Listener Adapter

The __socket-listener-adapter__ adapter provides the ability for the ClearBlade platform to receive data over a tcp or udp socket.

The adapter subscribes to MQTT topics which are used to interact with the socket. The adapter publishes any data received from the socket to MQTT topics so that the ClearBlade Platform is able to process the data as needed.

# MQTT Topic Structure
The __socket-listener-adapter__ adapter utilizes MQTT messaging to communicate with the ClearBlade Platform. The adapter will publish messages to MQTT topics in order to send socket data to the ClearBlade Platform/Edge. The topic structures utilized by the adapter are as follows:

  * {__TOPIC ROOT__}/{__PROTOCOL__}/{__LISTENER PORT__}/incoming-data

## ClearBlade Platform Dependencies
The socket-listener-adapter adapter was constructed to provide the ability to communicate with a _System_ defined in a ClearBlade Platform instance. Therefore, the adapter requires a _System_ to have been created within a ClearBlade Platform instance.

Once a System has been created, artifacts must be defined within the ClearBlade Platform system to allow the adapters to function properly. At a minimum: 

  * A device needs to be created in the Auth --> Devices collection. The device will represent the adapter, or more importantly, the Septentrio GNSS receiver or the gateway on which the adapter is executing. The _name_ and _active key_ values specified in the Auth --> Devices collection will be used by the adapter to authenticate to the ClearBlade Platform or ClearBlade Edge. 
  * An adapter configuration data collection needs to be created in the ClearBlade Platform _system_ and populated with the data appropriate to the socket-listener-adapter installation. The schema of the data collection should be as follows:


| Column Name      | Column Datatype |
| ---------------- | --------------- |
| adapter_name     | string          |
| topic_root       | string          |
| adapter_settings | string (json)   |

### adapter_settings
The adapter_settings column will need to contain an array of JSON objects containing the following attributes:

##### protocol
* Either __tcp__ or __udp__

##### listen_port
* The tcp/udp port the adapter should connect to

##### message_end_character
* The protocol specific character used to terminate the data stream

##### transform_to_hex
* A boolean flag indicating whether or not the received data should be transformed to hex

#### adapter_settings_examples

##### TCP connection type example
```json
[
  {
    "protocol": "tcp",
    "listen_port": "64010",
    "message_end_character": "\n"
  },
  {
    "protocol": "tcp",
    "listen_port": "64011",
    "message_end_character": "\n"
  },
  {
    "protocol": "udp",
    "listen_port": "2018",
    "message_end_character": "\n",
    "transform_to_hex": true
  }
]
```

## Usage

### Executing the adapter

`socket-listener-adapter -systemKey=<SYSTEM_KEY> -systemSecret=<SYSTEM_SECRET> -platformURL=<PLATFORM_URL> -messagingURL=<MESSAGING_URL> -deviceName=<DEVICE_NAME> -password=<DEVICE_ACTIVE_KEY> -adapterConfigCollection=<COLLECTION_NAME> -logLevel=<LOG_LEVEL>`

   __*Where*__ 

   __systemKey__
  * REQUIRED
  * The system key of the ClearBLade Platform __System__ the adapter will connect to

   __systemSecret__
  * REQUIRED
  * The system secret of the ClearBLade Platform __System__ the adapter will connect to
   
   __deviceName__
  * The device name the adapter will use to authenticate to the ClearBlade Platform
  * Requires the device to have been defined in the _Auth - Devices_ collection within the ClearBlade Platform __System__
  * OPTIONAL
  * Defaults to __septentrio-gnss-adapter__
   
   __password__
  * REQUIRED
  * The active key the adapter will use to authenticate to the platform
  * Requires the device to have been defined in the _Auth - Devices_ collection within the ClearBlade Platform __System__
   
   __platformUrl__
  * The url of the ClearBlade Platform instance the adapter will connect to
  * OPTIONAL
  * Defaults to __http://localhost:9000__

   __messagingUrl__
  * The MQTT url of the ClearBlade Platform instance the adapter will connect to
  * OPTIONAL
  * Defaults to __localhost:1883__

   __adapterConfigCollection__
  * REQUIRED 
  * The collection name of the data collection used to house adapter configuration data

   __logLevel__
  * The level of runtime logging the adapter should provide.
  * Available log levels:
    * fatal
    * error
    * warn
    * info
    * debug
  * OPTIONAL
  * Defaults to __info__

## Setup
---
The __socket-listener-adapter__  adapter is dependent upon the ClearBlade Go SDK and its dependent libraries being installed. The adapter adapter was written in Go and therefore requires Go to be installed (https://golang.org/doc/install).

### Adapter compilation
In order to compile the adapter for execution within linux, the following steps need to be performed:

 1. Retrieve the adapter source code  
    * ```git clone git@github.com:ClearBlade/socket-listener-adapter.git```
 2. Navigate to the adapter directory  
    * ```cd socket-listener-adapter```
 4. Compile the adapter
    * ```GOARCH=arm64 GOOS=linux go build```
    * ```GOARCH=arm GOARM=7 GOOS=linux go build```
