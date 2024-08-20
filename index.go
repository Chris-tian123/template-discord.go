package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// Token from environment variable
const token = "YOUR_BOT_TOKEN"
const prefix = "YOUR_PREFIX"

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the message handler
	dg.AddHandler(messageCreate)

	// Open a connection to Discord
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Set the bot's activity after successfully opening the connection
	err = dg.UpdateGameStatus(4, "Made with love <3 by Pella.app")
	if err != nil {
		fmt.Println("error updating status,", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")
	select {}
}


func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Content) > 0 && m.Content[0] == prefix {
		command := m.Content[1:]

		switch command {
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		case "help":
			helpMessage := `**Available commands:**
- **ping**: Replies with Pong!
- **help**: Lists all commands.
- **joke**: Tells a random joke.
- **quote**: Sends a random inspirational quote.`
			s.ChannelMessageSend(m.ChannelID, helpMessage)
		case "joke":
			joke, err := fetchJoke()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't fetch a joke right now.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, joke)
		case "quote":
			quote, err := fetchQuote()
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Sorry, I couldn't fetch a quote right now.")
				return
			}
			s.ChannelMessageSend(m.ChannelID, quote)
		}
	}
}

func fetchJoke() (string, error) {
	resp, err := http.Get("https://official-joke-api.appspot.com/random_joke")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var joke map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &joke); err != nil {
		return "", err
	}

	setup := joke["setup"].(string)
	punchline := joke["punchline"].(string)
	return setup + " - " + punchline, nil
}

func fetchQuote() (string, error) {
	resp, err := http.Get("https://api.quotable.io/random")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var quote map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &quote); err != nil {
		return "", err
	}

	content := quote["content"].(string)
	author := quote["author"].(string)
	return "\"" + content + "\" - " + author, nil
}
