package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/vrld/ansicht/internal/model"
	notmuch "github.com/zenhack/go.notmuch"
	ini "gopkg.in/ini.v1"
)

func FindThreads(query *model.SearchQuery) (model.SearchResult, error) {
	db, err := notmuch.OpenWithConfig(nil, nil, nil, notmuch.DBReadOnly)
	if err != nil {
		return model.SearchResult{}, fmt.Errorf("cannot open notmuch database: %v", err)
	}
	defer db.Close()

	notmuchQuery := db.NewQuery(query.Query)
	if notmuchQuery == nil {
		return model.SearchResult{}, fmt.Errorf("cannot create query: %v", query.Query)
	}

	threads, err := notmuchQuery.Threads()
	if err != nil {
		return model.SearchResult{}, fmt.Errorf("cannot get threads: %v", err)
	}

	var notmuchThread *notmuch.Thread
	result := model.SearchResult{Query: query}
	for threads.Next(&notmuchThread) {
		if notmuchThread == nil {
			panic("unexpected nil in threads.Next()")
		}

		result.Threads = append(result.Threads, ThreadFromNotmuch(notmuchThread))
	}

	return result, nil
}

func ThreadFromNotmuch(nmThread *notmuch.Thread) model.Thread {
	matchedAuthors, authors := nmThread.Authors()

	// read messages
	var messages []model.Message
	nmMessages := nmThread.Messages()
	var nmMessage *notmuch.Message
	for nmMessages.Next(&nmMessage) {
		if nmMessage == nil {
			panic("unexpected nil in messages.Next()")
		}
		messages = append(messages, MessageFromNotmuch(nmMessage))
	}

	// build thread model
	return model.Thread{
		ID:                   nmThread.ID(),
		Authors:              append(matchedAuthors, authors...),
		Subject:              nmThread.Subject(),
		Tags:                 ReadTags(nmThread.Tags()),
		NewestDate:           nmThread.NewestDate(),
		OldestDate:           nmThread.OldestDate(),
		CountMatchedMessages: nmThread.CountMatched(),
		Messages:             messages,
	}
}

func MessageFromNotmuch(nmMessage *notmuch.Message) model.Message {
	return model.Message{
		ID:       nmMessage.ID(),
		ThreadID: nmMessage.ThreadID(),
		Date:     nmMessage.Date(),
		Filename: nmMessage.Filename(),
		Tags:     ReadTags(nmMessage.Tags()),
		From:     nmMessage.Header("from"),
		To:       nmMessage.Header("to"),
		Subject:  nmMessage.Header("subject"),
		Flags:    MessageFlagsFromFilename(nmMessage.Filename()),
	}
}

func ReadTags(nmTags *notmuch.Tags) []string {
	var tags []string
	var nmTag *notmuch.Tag
	for nmTags.Next(&nmTag) {
		if nmTag == nil {
			panic("unexpected nil in tags.Next()")
		}
		tags = append(tags, nmTag.Value)
	}
	return tags
}

func MessageFlagsFromFilename(filename string) model.MessageFlags {
	flags := model.MessageFlags{}

	flagsStartIndex := strings.LastIndex(filename, "2,")
	if flagsStartIndex == -1 {
		return flags
	}

	for _, char := range filename[flagsStartIndex+2:] {
		switch char {
		case 'D':
			flags.Draft = true
		case 'F':
			flags.Flagged = true
		case 'P':
			flags.Passed = true
		case 'R':
			flags.Replied = true
		case 'S':
			flags.Seen = true
		case 'T':
			flags.Trashed = true
		case 'Z': // stop reading after the last possible flag
			break
		}
	}

	return flags
}

func GetSavedQueries() ([]model.SearchQuery, error) {
	// config lists did not work out for some reason, i.e., this does *not* work:
	//     configList, err := db.GetConfigLilst("query")
  //     for configList.Next(...) {...}
	// Next() always returned false
	//
	// So we just parse the config file ourselves.
	configPath, err := NotmuchConfigLocation()
	if err != nil {
		return nil, fmt.Errorf("cannot find config file: %v", err)
	}

	config, err := ini.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file from %s: %v", configPath, err)
	}

	querySection, err := config.GetSection("query")
	if err != nil {
		// no query section
		return []model.SearchQuery{}, nil
	}

	var queries []model.SearchQuery
	for _, key := range querySection.Keys() {
		queries = append(queries, model.SearchQuery{
				Name:  key.Name(),
				Query: key.Value(),
		})
	}

	return queries, nil
}

func NotmuchConfigLocation() (string, error) {
	// Search order as specified in man notmuch-config
	if notmuchConfig := os.Getenv("NOTMUCH_CONFIG"); notmuchConfig != "" {
		if _, err := os.Stat(notmuchConfig); err == nil {
			return notmuchConfig, nil
		}
		return "", fmt.Errorf("config file specified by NOTMUCH_CONFIG does not exist: %s", notmuchConfig)
	}

	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			xdgConfigHome = filepath.Join(homeDir, ".config")
		}
	}

	profile := os.Getenv("NOTMUCH_PROFILE")
	if profile == "" {
		profile = "default"
	}
	xdgConfigPath := filepath.Join(xdgConfigHome, "notmuch", profile, "config")
	if _, err := os.Stat(xdgConfigPath); err == nil {
		return xdgConfigPath, nil
	}

	if homeDir, err := os.UserHomeDir(); err != nil {
		return "", fmt.Errorf("unable to determine home directory: %v", err)
	} else {
		homeConfigPath := filepath.Join(homeDir, ".notmuch-config."+profile)
		if _, err := os.Stat(homeConfigPath); err == nil {
			return homeConfigPath, nil
		}

		homeConfigPath = filepath.Join(homeDir, ".notmuch-config")
		if _, err := os.Stat(homeConfigPath); err == nil {
			return homeConfigPath, nil
		}
	}

	return "", fmt.Errorf("no config file found in any of the expected locations")
}
