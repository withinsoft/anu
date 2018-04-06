package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joeshaw/envdecode"
)

type config struct {
	Prefix       string `env:"COMMAND_PREFIX,default=;"`
	DiscordToken string `env:"DISCORD_TOKEN,required"`
	Home         string `env:"ANU,default=."`
}

func main() {
	var cfg config
	err := envdecode.StrictDecode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	dg, err := discordgo.New("Bot " + cfg.DiscordToken)
	if err != nil {
		log.Fatalf("can't create discord session: %v", err)
	}

	dg.AddHandler(messageCreate(cfg.Prefix, cfg.Home))

	err = dg.Open()
	if err != nil {
		log.Fatalf("can't open discord websocket: %v", err)
	}
	defer dg.Close()

	log.Println("Anu is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(prefix, home string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if m.Author.ID == s.State.User.ID {
			return
		}

		// if the message doesn't start with the command prefix, ignore it.
		if !strings.HasPrefix(m.ContentWithMentionsReplaced(), prefix) {
			return
		}

		ch, err := s.State.Channel(m.ChannelID)
		if err != nil {
			log.Printf("error in finding channel %s: %v", m.ChannelID, err)
			return
		}

		gu, err := s.State.Guild(ch.GuildID)
		if err != nil {
			log.Printf("error in finding guild %s with channel %s:%s: %v", ch.GuildID, ch.ID, ch.Name, err)
			return
		}

		msg, err := m.ContentWithMoreMentionsReplaced(s)
		if err != nil {
			log.Printf("can't create message corpus %s %s: %v", ch.ID, m.Author, err)
			return
		}

		userCmd := strings.Fields(msg)[0][1:]

		us, err := s.State.Member(gu.ID, m.Author.ID)
		if err != nil {
			log.Printf("can't look up guild %s user %s: %v", gu.ID, m.Author.ID, err)
			return
		}

		nick := us.Nick
		if nick == "" {
			nick = m.Author.Username
		}

		env := []string{}
		env = append(env, "DISCORD_CHANNEL_ID="+m.ChannelID)
		env = append(env, "DISCORD_CHANNEL_NAME="+ch.Name)
		env = append(env, "DISCORD_GUILD_ID="+ch.GuildID)
		env = append(env, "DISCORD_GUILD_NAME="+gu.Name)
		env = append(env, "DISCORD_MESSAGE_ID="+m.ID)
		env = append(env, "DISCORD_MESSAGE_AUTHOR_ID="+m.Author.ID)
		env = append(env, "DISCORD_MESSAGE_AUTHOR_USERNAME="+m.Author.Username+"#"+m.Author.Discriminator)
		env = append(env, "DISCORD_MESSAGE_AUTHOR_NICK="+nick)
		env = append(env, "I_AM_LISTENING=for a sound beyond sound")
		env = append(env, "VERB="+userCmd)
		env = append(env, "PWD="+home)
		env = append(env, "HOME="+home)
		env = append(env, "USER="+s.State.User.Username)

		buf := bytes.NewBuffer([]byte(msg))
		out := bytes.NewBuffer(nil)

		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, filepath.Join(home, userCmd))
		cmd.Stdin = buf
		cmd.Stdout = out
		cmd.Stderr = os.Stderr
		cmd.Env = env
		cmd.Dir = home

		var reply string
		err = cmd.Run()
		if err != nil {
			reply = fmt.Sprintf("oopsie whoopsie!\nuwu\nwe made a fucky wucky %s!!1 a wittle fucko boingo! the code monkies at our headquarters are working VEWY HAWD to fix this!: %v", userCmd, err)
			log.Printf("error in %s:%v: %v", userCmd, env, err)
		} else {
			reply = out.String()
		}

		_, err = s.ChannelMessageSend(ch.ID, reply)
		if err != nil {
			log.Printf("would have sent: %q", reply)
			log.Printf("can't send message: %v", err)
		}
	}
}
