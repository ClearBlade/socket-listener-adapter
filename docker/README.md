# Docker Image Creation

## Prerequisites

- Building the image requires internet acess

### Creating the Docker image for the socket-listener-adapter adapter

Clone this repository and execute the following commands to create a docker image for the socket-listener-adapter:  

- ```docker build --no-cache -f docker/Dockerfile -t clearblade/socket-listener-adapter:{version} -t clearblade/socket-listener-adapter:latest .```

docker buildx build --no-cache --platform linux/amd64 -t clearblade/socket-listener-adapter:1.0.0 -t clearblade/socket-listener-adapter:latest -f docker/Dockerfile .

#### Clean docker build cache
- ```docker builder prune```

# Using the adapter

## Deploying the adapter image

When the docker image has been created, it will need to be saved and imported into the runtime environment. Execute the following steps to save and deploy the adapter image

- On the machine where the ```docker build``` command was executed, execute ```docker save clearblade/socket-listener-adapter:{version} -o socket-listener-adapter.tar```. Optionally add ```| gzip > socket-listener-adapter.tar.gz``` to create a .tar.gz image.

- On the server where docker is running, execute ```docker load -i socket-listener-adapter.tar```

## Executing the adapter

Once you create the docker image, start the socket-listener-adapter using the following command:


```docker run -d --name socket-listener-adapter --network cb-net --restart always clearblade/socket-listener-adapter:{version} -systemKey=<SYSTEM_KEY> -systemSecret=<SYSTEM_SECRET> -platformURL=<PLATFORM_URL> -messagingURL=<MESSAGING_URL> -deviceName=<DEVICE_NAME> -password=<DEVICE_ACTIVE_KEY> -adapterConfigCollection=<COLLECTION_NAME> -logLevel=<LOG_LEVEL>```

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
```docker run -d --name clearblade/socket-listener-adapter:1.0.0 --network cb-net --restart always socket-listener-adapter -systemKey=cc9d8bba0bfeeed78595c4dfbb0b -systemSecret=CC9D8BBA0BB4F1C5AD8994E6D41B -platformURL=https://platform.clearblade.com -messagingURL=platform.clearblade.com:8901 -deviceId=socket-listener-adapter -password=01234567890 -logLevel=debug```
