package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	sess, err := discordgo.New("Bot MTE4NDUyMjM3Njc4MDcxNDEyNA.GgLv1P.oM1D_jRTDi1Zul0r_TCBMePanLviQ7cd_5mP7I")
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(hello)
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

// func boot

func hello(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}

	if message.Content == "Hi Alvin" {
		session.ChannelMessageSend(message.ChannelID, "Hello "+message.Author.Username+", this is Alvin Beta. My days are numbered, I walk so other future Alvins can run. Do not let the dark and evil entities of this world consume you.")
	}
}

func welcome(session *discordgo.Session, message *discordgo.GuildMemberAdd) {

	session.ChannelMessageSend(message.GuildID, "Welcome Traveler! My name is Alvin")

}
