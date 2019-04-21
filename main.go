package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	com "github.com/unswpcsoc/pcsocbot/commands"
)

var (
	dg     *discordgo.Session
	router Router
	prefix = "!"
)

// Router Routes a command string to a handler
type Router struct {
	routes map[string]com.Command
}

// AddCommand Adds command-string mapping
func (r *Router) AddCommand(command Command, name ...string) {
	// TODO: Tree
	// AddCommand(PingTagPeople, ask)
	// AddCommand(PingTagPeople, tags, ping)
}

// Route Routes to handler from string
func (r *Router) Route(cmd string) (handlers.Handler, bool) {
	// TODO: Tree
	if routes == nil {
		return nil, false
	}
	return routes, true
}

// init for discordgo things
func init() {
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("Missing Discord API Key: TOKEN")
	}

	var err error
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error:", err)
	}

	err = dg.Open()
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
}

// command initialisation
func init() {

}

func main() {
	var err error
	dg.UpdateListeningStatus("you")

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		message := strings.TrimSpace(m.Content)

		if !strings.HasPrefix(message, prefix) {
			return
		}

		s.ChannelTyping(m.ChannelID)

		command, found := router.Route(message)
		if !found {
			s.ChannelMessageSend(m.ChannelID"Error: Unknown command")
		}

		// Call handler
		str, err = command.MsgHandle()
		if err != nil {
			log.Println(err)
			s.ChannelMessageSend(m.ChannelID, err.String())
			return
		}
	})

	// Don't close the connection, wait for a kill signal
	log.Println("Logged in as: ", dg.State.User.ID)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	sig := <-sc
	log.Println("Received Signal: " + sig.String())
	log.Println("Bye!")
	dg.Close()
}
