package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	General struct {
		Domain     string `yaml:"domain"`
		ServerID   string `yaml:"server_id"`
		SSLEnabled bool   `yaml:"ssl_enabled"`
		SSL        struct {
			CertFile string `yaml:"cert_file"`
			KeyFile  string `yaml:"key_file"`
		} `yaml:"ssl"`
	} `yaml:"general"`
	Bot struct {
		Token string `yaml:"token"`
	} `yaml:"bot"`
	Cache struct {
		Expire        string `yaml:"expire"`
		CleanInterval string `yaml:"clean_interval"`
	} `yaml:"cache"`
}

var (
	config *Config
	Cache  *cache.Cache
	logger *log.Logger
)

func main() {
	logger = log.New(os.Stdout, "[Server] ", log.LstdFlags|log.LUTC)

	config = readConfig()
	expire, _ := strconv.Atoi(config.Cache.Expire)
	cleanInterval, _ := strconv.Atoi(config.Cache.CleanInterval)
	expire = expire * 1
	cleanInterval = cleanInterval * 1
	Cache = cache.New(time.Duration(expire)*time.Second, time.Duration(cleanInterval)*time.Second)
	go startBot()
	startHTTPServer()
}

func readConfig() *Config {
	file, err := os.Open("config.yaml")
	if err != nil {
		logger.Fatal(err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Fatal(err)
	}

	return &config
}

func startHTTPServer() {
	http.HandleFunc("/new", handleNewLogin)
	http.HandleFunc("/verify", handleVerifyLogin)
	http.HandleFunc("/login", handleLogin)

	serverURL := config.General.Domain
	var server *http.Server

	if config.General.SSLEnabled {
		server = &http.Server{
			Addr:         serverURL + ":443",
			Handler:      nil, // Use default handler
			TLSConfig:    &tls.Config{},
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
	} else {
		server = &http.Server{
			Addr:    serverURL + ":80",
			Handler: nil,
		}
	}

	if config.General.SSLEnabled {
		logger.Fatal(server.ListenAndServeTLS(config.General.SSL.CertFile, config.General.SSL.KeyFile))
	} else {
		logger.Fatal(server.ListenAndServe())
	}
}

func handleNewLogin(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	logger.Printf("New login request from IP: %s\n", ip)

	username := r.URL.Query().Get("username")
	tag := r.URL.Query().Get("tag")
	if username == "" || tag == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println("Invalid request parameters")
		return
	}
	isExist, _ := UserExist(username, tag)
	if isExist {
		uuidStr := uuid.New().String()
		Cache.Set(uuidStr, username, cache.DefaultExpiration)
		discordSession, err := discordgo.New("Bot " + config.Bot.Token)
		_, userID := UserExist(username, tag)
		userChannel, err := discordSession.UserChannelCreate(userID)
		if err == nil {
			loginURL := generateLoginURL(uuidStr)
			sendMessageToUser(discordSession, userChannel.ID, "Click the following link to log in: "+loginURL)
		}
		response := struct {
			Success bool   `json:"success"`
			UUID    string `json:"uuid"`
		}{
			Success: true,
			UUID:    uuidStr,
		}
		Response(w, http.StatusOK, response)
		logger.Printf("New login request processed successfully, UUID: %s\n", uuidStr)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	logger.Printf("Login request from IP: %s\n", ip)

	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println("Invalid request parameters")
		return
	}

	data, found := Cache.Get(uuidStr)
	if found {
		if data.(string) != "true" {
			Cache.Set(uuidStr, "true", cache.DefaultExpiration)
			fmt.Fprintf(w, "Login success")
		}
	}

	logger.Println("Login request processed successfully")
}

func handleVerifyLogin(w http.ResponseWriter, r *http.Request) {
	ip := strings.Split(r.RemoteAddr, ":")[0]
	logger.Printf("Verify login request from IP: %s\n", ip)

	uuidStr := r.URL.Query().Get("uuid")
	if uuidStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		logger.Println("Invalid request parameters")
		return
	}

	data, found := Cache.Get(uuidStr)
	if found {
		if data.(string) == "true" {
			fmt.Fprintf(w, `{"success": true}`)
			logger.Println("Verify login request processed successfully")
			return
		}
	}

	fmt.Fprintf(w, `{"success": false}`)
	logger.Println("Verify login request processed successfully")
}

func startBot() {
	discordSession, err := discordgo.New("Bot " + config.Bot.Token)
	if err != nil {
		logger.Fatal("Error creating Discord session: ", err)
		return
	}

	discordSession.AddHandler(handleMessageCreate)

	err = discordSession.Open()
	if err != nil {
		logger.Fatal("Error opening Discord session: ", err)
		return
	}

	logger.Println("Bot is now running.")

	<-make(chan struct{})
}

func handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	return
}

func UserExist(username string, tag string) (bool, string) {
	api := fmt.Sprintf("https://discord.com/api/v9/guilds/%s/members?limit=1000", config.General.ServerID)
	req, err := http.NewRequest("GET", api, nil)
	req.Header.Set("Authorization", "Bot "+config.Bot.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var users []struct {
		User struct {
			ID            string `json:"id"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Discriminator string `json:"discriminator"`
			PublicFlags   int    `json:"public_flags"`
			Flags         int    `json:"flags"`
		} `json:"user"`
	}

	err = json.Unmarshal(body, &users)
	if err != nil {
		logger.Fatal(err)
	}

	for _, u := range users {
		if u.User.Username == username && u.User.Discriminator == tag {
			return true, u.User.ID
		}
	}

	return false, ""
}

func generateLoginURL(uuid string) string {
	ProtocolHead := "http://"
	if config.General.SSLEnabled {
		ProtocolHead = "https://"
		return fmt.Sprintf("%s%s/login?uuid=%s", ProtocolHead, config.General.Domain, uuid)
	}
	return fmt.Sprintf("%s%s/login?uuid=%s", ProtocolHead, config.General.Domain, uuid)
}

func sendMessageToUser(s *discordgo.Session, channelID string, message string) {
	_, err := s.ChannelMessageSend(channelID, message)
	if err != nil {
		logger.Printf("Error sending message to user. Channel ID: %s, Error: %s\n", channelID, err)
	}
}

func Response(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
