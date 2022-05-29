package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type playerStats struct {
	Icon         string `json:"icon"`
	Name         string `json:"name"`
	Level        int    `json:"level"`
	LevelIcon    string `json:"levelIcon"`
	Prestige     int    `json:"prestige"`
	PrestigeIcon string `json:"prestigeIcon"`
	Rating       string `json:"rating"`
	RatingIcon   string `json:"ratingIcon"`
	GamesWon     int    `json:"gamesWon"`
	Ratings      []struct {
		Level int    `json:"level"`
		Role  string `json:"role"`
	}
	QuickPlayStats struct {
		EliminationsAvg   float64 `json:"eliminationsAvg"`
		DamageDoneAvg     int     `json:"damageDoneAvg"`
		DeathsAvg         float64 `json:"deathsAvg"`
		FinalBlowsAvg     float64 `json:"finalBlowsAvg"`
		HealingDoneAvg    int     `json:"healingDoneAvg"`
		ObjectiveKillsAvg float64 `json:"objectiveKillsAvg"`
		ObjectiveTimeAvg  string  `json:"objectiveTimeAvg"`
		SoloKillsAvg      float64 `json:"soloKillsAvg"`
		Games             struct {
			Played int `json:"played"`
			Won    int `json:"won"`
		} `json:"games"`
		Awards struct {
			Cards        int `json:"cards"`
			Medals       int `json:"medals"`
			MedalsBronze int `json:"medalsBronze"`
			MedalsSilver int `json:"medalsSilver"`
			MedalsGold   int `json:"medalsGold"`
		} `json:"awards"`
	} `json:"quickPlayStats"`
	CompetitiveStats struct {
		EliminationsAvg   float64 `json:"eliminationsAvg"`
		DamageDoneAvg     int     `json:"damageDoneAvg"`
		DeathsAvg         float64 `json:"deathsAvg"`
		FinalBlowsAvg     float64 `json:"finalBlowsAvg"`
		HealingDoneAvg    int     `json:"healingDoneAvg"`
		ObjectiveKillsAvg float64 `json:"objectiveKillsAvg"`
		ObjectiveTimeAvg  string  `json:"objectiveTimeAvg"`
		SoloKillsAvg      float64 `json:"soloKillsAvg"`
		Games             struct {
			Played int `json:"played"`
			Won    int `json:"won"`
		} `json:"games"`
		Awards struct {
			Cards        int `json:"cards"`
			Medals       int `json:"medals"`
			MedalsBronze int `json:"medalsBronze"`
			MedalsSilver int `json:"medalsSilver"`
			MedalsGold   int `json:"medalsGold"`
		} `json:"awards"`
	} `json:"competitiveStats"`
}

func main() {
	fmt.Println("Starting bot...")
	fmt.Println(os.Getenv("DISCORD_OW_TOKEN"))

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_OW_TOKEN"))
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	discord.AddHandler(messageCreate)
	discord.AddHandler(ready)
	err = discord.Open()

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}
	<-make(chan struct{})
	defer discord.Close()
	fmt.Println("Closing bot...")

}
func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "Active")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!stats ") {

		var btag = strings.SplitAfter(m.Content, "!stats ")[1]
		fmt.Println(btag)
		var response, err = http.Get("https://ow-api.com/v1/stats/pc/us/" + strings.Replace(btag, "#", "-", -1) + "/profile")
		fmt.Println(strings.Replace(btag, "#", "-", -1))
		if err != nil {
			fmt.Println(err)
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(body))
		bytes := []byte(string(body))
		var stats playerStats
		json.Unmarshal(bytes, &stats)
		fmt.Println(stats)
		for _, v := range stats.Ratings {
			fmt.Printf("%s: %d\n", v.Role, v.Level)
		}
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0x00ff00, // Green
			Description: btag,
			Fields: []*discordgo.MessageEmbedField{
				&discordgo.MessageEmbedField{
					Name:   strings.Title(stats.Ratings[0].Role),
					Value:  fmt.Sprintf("%d", stats.Ratings[0].Level),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   strings.Title(stats.Ratings[1].Role),
					Value:  fmt.Sprintf("%d", stats.Ratings[1].Level),
					Inline: true,
				},
				&discordgo.MessageEmbedField{
					Name:   strings.Title(stats.Ratings[2].Role),
					Value:  fmt.Sprintf("%d", stats.Ratings[2].Level),
					Inline: true,
				},
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: stats.Icon,
			},
			Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
			Title:     btag,
		}
		fmt.Println("Sending embed")
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		s.ChannelMessageSend(m.ChannelID, "https://playoverwatch.com/en-us/career/pc/us/"+strings.Replace(btag, "#", "-", -1))
		fmt.Println("sent embed")

	}

}

func commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}
	data := i.ApplicationCommandData()
	switch data.Options[0].Name {
	case "register":
		s.ChannelMessageSend(i.ChannelID, "Registering...")
	}
}
