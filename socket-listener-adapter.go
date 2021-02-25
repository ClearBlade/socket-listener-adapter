package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	adapter_library "github.com/clearblade/adapter-go-library"
)

const (
	adapterName = "socket-listener-adapter"

	msgPublishQOS              = 0
	defaultTopicRoot           = "socket-listener"
	defaultListenPort          = "12345"
	defaultMessageEndCharacter = ""
	defaultDeviceName          = "socket-listener-adapter"
	defaultPlatformURL         = "http://localhost:9000"
	defaultMessagingURL        = "localhost:1883"
	defaultLogLevel            = "info"
	defaultLogFilePath         = "stdout"
)

var (
	adapterConfig   *adapter_library.AdapterConfig
	adapterSettings *[]socketAdapterSettings
)

type socketAdapterSettings struct {
	Protocol            string `json:"protocol"`
	ListenPort          string `json:"listen_port"`
	MessageEndCharacter string `json:"message_end_character"`
	TransformToHex      bool   `json:"transform_to_hex"`
}

// func init() {
// 	flag.StringVar(&sysKey, "systemKey", "", "system key (required)")
// 	flag.StringVar(&sysSec, "systemSecret", "", "system secret (required)")
// 	flag.StringVar(&deviceName, "deviceName", defaultDeviceName, "name of device (optional)")
// 	flag.StringVar(&activeKey, "activeKey", "", "active key for device authentication (required)")
// 	flag.StringVar(&platformURL, "platformURL", defaultPlatformURL, "platform url (optional)")
// 	flag.StringVar(&messagingURL, "messagingURL", defaultMessagingURL, "messaging URL (optional)")
// 	flag.StringVar(&logLevel, "logLevel", defaultLogLevel, "The level of logging to use. Available levels are 'debug, 'info', 'warn', 'error', 'fatal' (optional)")
// 	flag.StringVar(&logFilePath, "logFilePath", defaultLogFilePath, "Path for the log file of the adapter (optional - default is stdout, provide a full path to desired log file output location)")
// 	flag.StringVar(&adapterConfigCollID, "adapterConfigCollectionID", "", "The ID of the data collection used to house adapter configuration (required)")
// }

// func usage() {
// 	log.Printf("Usage: socket-listener-adapter [options]\n\n")
// 	flag.PrintDefaults()
// }

// func validateFlags() {
// 	flag.Parse()

// 	if sysKey == "" || sysSec == "" || activeKey == "" || adapterConfigCollID == "" {
// 		log.Println("ERROR - Missing required flags")
// 		flag.Usage()
// 		os.Exit(1)
// 	}

// }

func main() {
	err := adapter_library.ParseArguments(adapterName)
	if err != nil {
		log.Fatalf("[FATAL] Failed to parse arguments: %s\n", err.Error())
	}

	adapterConfig, err = adapter_library.Initialize()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize: %s\n", err.Error())
	}

	adapterSettings = &[]socketAdapterSettings{}
	err = json.Unmarshal([]byte(adapterConfig.AdapterSettings), adapterSettings)
	if err != nil {
		log.Fatalf("[FATAL] Failed to parse Adapter Settings: %s\n", err.Error())
	}

	err = adapter_library.ConnectMQTT("", nil)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect MQTT: %s\n", err.Error())
	}

	go initializeSockets()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			log.Println("[INFO] socket-listener-adapter still listening")
		}
	}

	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	// filter := &logutils.LevelFilter{
	// 	Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"},
	// 	MinLevel: logutils.LogLevel(strings.ToUpper(logLevel)),
	// }
	// if logFilePath == "stdout" {
	// 	log.Println("using stdout for logging")
	// 	filter.Writer = os.Stdout
	// } else {
	// 	log.Printf("using %s for logging\n", logFilePath)
	// 	logfile, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	// 	if err != nil {
	// 		log.Fatalf("Failed to open log file: %s", err.Error())
	// 	}
	// 	defer logfile.Close()
	// 	filter.Writer = logfile
	// }
	// log.SetOutput(filter)

	// initClearBlade()
	// initAdapterConfig()
	// connectClearBlade()

	// log.Println("[DEBUG] main - starting info log ticker")

}

func initializeSockets() {
	log.Println("[INFO] initializeSockets - Creating socket listeners...")
	for _, socketConfig := range *adapterSettings {
		if socketConfig.Protocol == "tcp" {
			go createTCPListener(socketConfig)
		} else {
			go createUDPListener(socketConfig)
		}
	}
}

// func onConnect(client mqtt.Client) {
// 	log.Println("[INFO] onConnect - ClearBlade MQTT successfully connected")
// 	for _, socketConfig := range config.AdapterSettings {
// 		if socketConfig.Protocol == "tcp" {
// 			go createTCPListener(socketConfig)
// 		} else {
// 			go createUDPListener(socketConfig)
// 		}
// 	}
// }

func createTCPListener(socketConfig socketAdapterSettings) {
	log.Println("[INFO] createTCPListener - Creating TCP Listener")

	listener, err := net.Listen(socketConfig.Protocol, ":"+socketConfig.ListenPort)
	if err != nil {
		log.Fatalf("[FATAL] createTCPListener - Error Creating TCP Listener with port %s: %s", socketConfig.ListenPort, err.Error())
	}

	log.Println("[INFO] createTCPListener - Listener opened and ready to accecpt connections")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] createTCPListener - Failed to accept new connection: %s\n", err.Error())
			break
		}
		go handleTCPConnection(conn, socketConfig)
	}
}

func createUDPListener(socketConfig socketAdapterSettings) {
	log.Println("[INFO] createUDPListener - Creating UDP Listener")

	pc, err := net.ListenPacket("udp", ":"+socketConfig.ListenPort)
	if err != nil {
		log.Printf("[ERROR] createUDPListener - Error  creating UDP listener with port %s: %s", socketConfig.ListenPort, err.Error())
		return
	}

	log.Println("[INFO] createUDPListener - Listener opened and ready to accept connections")

	buffer := make([]byte, 65535)
	for {
		n, _, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Printf("[ERROR] createUDPListener - Failed to read from UDP with port %s: %s\n", socketConfig.ListenPort, err.Error())
		}
		if socketConfig.TransformToHex {
			publishMessage(fmt.Sprintf("%x", buffer[0:n-1]), socketConfig)
		} else {
			publishMessage(string(buffer[0:n-1]), socketConfig)
		}
	}
}

func handleTCPConnection(conn net.Conn, socketConfig socketAdapterSettings) {
	log.Println("[INFO] handleConnection - New TCP connection accepted")

	defer func() {
		log.Println("[INFO] handleConnection - closing TCP connection")
		conn.Close()
	}()

	if socketConfig.MessageEndCharacter == "" {
		log.Println("[INFO] handleConnection - Reading all data until TCP connection is closed...")
		bytes, err := ioutil.ReadAll(conn)
		if err != nil {
			log.Printf("[ERROR] handleConnection - Failed to read data from TCP connection: %s\n", err.Error())
			return
		}
		publishMessage(string(bytes[:]), socketConfig)
	} else {
		log.Printf("[INFO] handleConnection - Reading all data until character %s\n", socketConfig.MessageEndCharacter)
		scanner := bufio.NewScanner(bufio.NewReader(conn))
		split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}

			if i := strings.Index(string(data), socketConfig.MessageEndCharacter); i >= 0 {
				return i + 1, data[0:i], nil
			}

			if atEOF {
				return len(data), data, nil
			}

			return
		}
		scanner.Split(split)
		for scanner.Scan() {
			log.Println("[INFO] handleConnection - Read line of data")
			publishMessage(scanner.Text(), socketConfig)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("[ERROR] handleConnection - Failed to read data from TCP connection: %s\n", err.Error())
			return
		}

	}
}

func publishMessage(msg string, socketConfig socketAdapterSettings) {
	log.Println("[INFO] publishMessage - Publishing message")
	topic := adapterConfig.TopicRoot + "/" + socketConfig.Protocol + "/" + socketConfig.ListenPort + "/incoming-data"
	log.Printf("publishing message to topic: %s\n", topic)
	msg = strings.Replace(msg, socketConfig.MessageEndCharacter, "", 1)
	if err := adapter_library.Publish(topic, []byte(msg)); err != nil {
		log.Printf("[ERROR] publishMessage - Failed to publish message: %s\n", err.Error())
	}
	log.Println("[DEBUG] publishMessage - message published!")
}
