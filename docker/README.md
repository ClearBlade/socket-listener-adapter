# Docker Image Creation

## Prerequisites

- Building the image requires internet acess

### Creating the Docker image for the socket-listener-adapter adapter

Clone this repository and execute the following commands to create a docker image for the socket-listener-adapter:  

- ```cd socket-listener-adapter/docker```
- ```docker build --no-cache -f ../docker/Dockerfile -t socket_listener_adapter  ..```

#### Clean docker build cache
- ```docker builder prune```

# Using the adapter

## Deploying the adapter image

When the docker image has been created, it will need to be saved and imported into the runtime environment. Execute the following steps to save and deploy the adapter image

- On the machine where the ```docker build``` command was executed, execute ```docker save socket_listener_adapter:latest -o socket_listener_adapter.tar``` 

- On the server where docker is running, execute ```docker load -i socket_listener_adapter.tar```

## Executing the adapter

Once you create the docker image, start the socket_listener_adapter using the following command:


```docker run -d --name socket_listener_adapter --network cb-net --restart always socket_listener_adapter -systemKey=<SYSTEM_KEY> -systemSecret=<SYSTEM_SECRET> -platformURL=<PLATFORM_URL> -messagingURL=<MESSAGING_URL> -deviceName=<DEVICE_NAME> -password=<DEVICE_ACTIVE_KEY> -adapterConfigCollection=<COLLECTION_NAME> -logLevel=<LOG_LEVEL>```

```
-systemKey The System Key of your System on the ClearBlade Platform
-systemSecret The System Secret of your System on the ClearBlade Platform
-platformURL The address of the ClearBlade Platform (ex. https://platform.clearblade.com:443)
-messagingURL The MQTT broker address (ex. tcp://platform.clearblade.com:1883)
-deviceName The name of a device account created on the ClearBlade Platform
-password The active key of a device created on the ClearBlade Platform
-logLevel The level of runtime logging the adapter should provide (fatal, error, warn, info, debug)
```

Ex.
```docker run -d --name socket_listener_adapter --restart always socket_listener_adapter -systemKey=8ee8d7eb0b84b69bdd899984b84e -systemSecret=8EE8D7EB0BD69796B1AC9BB5D0CC01 -platformURL=https://community.clearblade.com -messagingURL=community.clearblade.com:8901 -deviceName=socket_listener_adapter -password=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI4ZWU4ZDdlYjBiODRiNjliZGQ4OTk5ODRiODRlIDo6IHNvY2tldF9saXN0ZW5lcl9hZGFwdGVyIiwic2lkIjoiMDUyMjgzNTktOGE1Zi00ZjY3LWFmYTYtZGI2MjU1NWNjOTk4IiwidXQiOjMsInR0IjoxLCJleHAiOi0xLCJpYXQiOjE2MzU0NDcwNzV9.iVYK9eMIo2j1kRQAhWCFjWCsalSaH8wbp80njHE5WF8 -logLevel=debug```
