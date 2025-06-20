package db

import (
	"fmt"
	"strings"

	"github.com/vrld/ansicht/internal/model"
	"github.com/zenhack/go.notmuch"
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

	for _, char := range filename[flagsStartIndex + 2:] {
		switch char {
		case 'D':
			flags.Draft = true;
		case 'F':
			flags.Flagged = true;
		case 'P':
			flags.Passed = true;
		case 'R':
			flags.Replied = true;
		case 'S':
			flags.Seen = true;
		case 'T':
			flags.Trashed = true;
		case 'Z':  // stop reading after the last possible flag
			break
		}
	}

	return flags
}
