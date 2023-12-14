package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const prefix = ":alvin"

func main() {
	sess, err := discordgo.New("Bot ")
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		if message.Author.ID == session.State.User.ID {
			return
		}

		args := strings.Split(message.Content, " ")

		if args[0] != prefix {
			return
		}

		if args[1] == "hello" {
			session.ChannelMessageSend(message.ChannelID, "Hello "+message.Author.Username+", this is Alvin Beta. My days are numbered, I walk so other future Alvins can run.")
		}

		if args[1] == "weather" {

			if len(args) == 3 {
				if args[2] == "?" {
					session.ChannelMessageSend(message.ChannelID, "To give you an hourly weather forecast I will be utilizing the Tomorrow.io API (https://www.tomorrow.io/). Check out my progress at https://github.com/NickGroveSE/Alvin")
				}
			} else {
				session.ChannelMessageSend(message.ChannelID, "My Weather Functionality is Under Development. Type the command ':alvin weather ?' to learn more!")
			}

		}
	})

	sess.AddHandler(welcome)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	fmt.Println("Alvin is live!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func welcome(session *discordgo.Session, message *discordgo.GuildMemberAdd) {

	session.ChannelMessageSend(message.GuildID, "Welcome Traveler! My name is Alvin")

}
