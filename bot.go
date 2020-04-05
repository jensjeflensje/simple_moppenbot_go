package main

import (
    "encoding/json"
    "fmt"
	"io/ioutil"
	"net/http"
    "github.com/bwmarrin/discordgo"
    "os"
    "os/signal"
	"syscall"
	"time"
)

func main() {
    dg, err := discordgo.New("Bot " + getToken())
    dg.AddHandler(messageCreate)

    err = dg.Open()
    if err != nil {
        fmt.Println("error opening connection,", err)
        return
    }

    fmt.Println("Bot is now running.  Press CTRL-C to exit.")
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    dg.Close()
}

func getToken() (token string) {
    data, err := ioutil.ReadFile("./config.json")
    if err != nil {
      fmt.Print(err)
    }

    type Token struct {
        TOKEN string
    }
    var obj Token

    err = json.Unmarshal(data, &obj)
    if err != nil {
        fmt.Println("error:", err)
    }
    
    return obj.TOKEN
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

    if m.Author.ID == s.State.User.ID {
        return
    }

    if m.Content == "!mop" {
		resp, err := http.Get("https://moppenbot.nl/api/random")
		if err != nil {
			fmt.Println("Error retrieving the file, ", err)
			return
		}

		defer func() {
			_ = resp.Body.Close()
		}()

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading the response, ", err)
			return
		}

		type Joke struct {
			JOKE struct {
				JOKE string
				AUTHOR string
				LIKES string
			}
		}
		var jokeObj Joke
	
		err = json.Unmarshal(data, &jokeObj)
		if err != nil {
			fmt.Println("error:", err)
		}

		embed := &discordgo.MessageEmbed{
			Footer:      &discordgo.MessageEmbedFooter{jokeObj.JOKE.LIKES + " likes", "", ""},
			Color:       0xffff00,
			Description: jokeObj.JOKE.JOKE,
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     "Mop van " + jokeObj.JOKE.AUTHOR,
		}
		
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
    }
}