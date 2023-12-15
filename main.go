package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const prefix = ":alvin"

var weatherCodes = map[int16]string{
	0:    "Unknown",
	1000: "Clear",
	1100: "Mostly Clear",
	1101: "Partly Cloudy",
	1102: "Mostly Cloudy",
	1103: "Slightly Cloudy",
	1001: "Cloudy",
	2000: "Fog",
	2100: "Light Fog",
	4000: "Drizzle",
	4001: "Rain",
	4200: "Light Rain",
	4201: "Heavy Rain",
	5000: "Snow",
	5001: "Flurries",
	5100: "Light Snow",
	5101: "Heavy Snow",
	6000: "Freezing Drizzle",
	6001: "Freezing Rain",
	6200: "Light Freezing Rain",
	6201: "Heavy Freezing Rain",
	7000: "Ice Pellets",
	7101: "Heavy Ice Pellets",
	7102: "Light Ice Pellets",
	8000: "Thunderstorm",
}

type WeatherData struct {
	Data struct {
		Time   string
		Values struct {
			CloudBase                float32
			CloudCeiling             float32
			CloudCover               int16
			DewPoint                 float32
			FreezingRainIntensity    float32
			Humidity                 int16
			PrecipitationProbability int16
			PressureSurfaceLevel     float32
			RainIntensity            float32
			SleetIntensity           float32
			SnowIntensity            float32
			Temperature              float32
			TemperatureApparent      float32
			UVHealthConcern          int16
			UVIndex                  int16
			Visibility               float32
			WeatherCode              int16
			WindDirection            float32
			WindGust                 float32
			WindSpeed                float32
		}
	}
	Location interface{}
}

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

			url := "https://api.tomorrow.io/v4/weather/realtime?location=cincinnati&units=imperial&apikey="
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatal(err)
			}

			req.Header.Add("accept", "application/json")

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			var weather WeatherData

			if err := json.Unmarshal([]byte(body), &weather); err != nil {
				log.Fatal(err)
			}

			weatherMessage := fmt.Sprintf("In your area, the temperature outside is **%.0f F** and it feels like **%.0f F**.\nCurrently you are experiencing **%v**", weather.Data.Values.Temperature, weather.Data.Values.TemperatureApparent, weatherCodes[weather.Data.Values.WeatherCode])

			if weather.Data.Values.WeatherCode <= 1001 {
				weatherMessage = weatherMessage + " conditions."
			}

			if len(args) == 3 {
				if args[2] == "?" {
					session.ChannelMessageSend(message.ChannelID, "To give you an hourly weather forecast I will be utilizing the Tomorrow.io API (https://www.tomorrow.io/). Check out my progress at https://github.com/NickGroveSE/Alvin")
				}
			} else {
				session.ChannelMessageSend(message.ChannelID, weatherMessage)
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
