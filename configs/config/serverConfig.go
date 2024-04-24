package config

const (
	API_Server_Port = "8080"
	Server_IP = "127.0.0.1"
)

func APIServerUrl() string {
	return "http://" + Server_IP + ":" + API_Server_Port
}