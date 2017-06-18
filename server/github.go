package main

import (
	"context"
	"strings"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/github"
)

// handlePushEvent handles GitHub push events
func handlePushEvent(payload interface{}, header webhooks.Header) {

	log.Debug("Handling Push Request")

	// get playload for push event
	pl := payload.(github.PushPayload)
	// Do whatever you want from here...
	log.Debugf("%+v", pl)

	// get branch name
	ref := strings.TrimPrefix(pl.Ref, "refs/heads/")
	// get clone URL
	cloneURL := pl.Repository.CloneURL

	// do build
	ctx, cancel := context.WithCancel(context.Background())
	gCancelCommands = append(gCancelCommands, cancel)
	go gClient.BuildPushImage(ctx, cloneURL, ref, pl.Repository.Name, pl.Repository.FullName, pl.HeadCommit.ID, gNotify)
}

// handleCreateEvent handles GitHub create events
func handleCreateEvent(payload interface{}, header webhooks.Header) {

	log.Debug("Handling Create Request")

	// get playload for push event
	pl := payload.(github.CreatePayload)
	// Do whatever you want from here...
	log.Debugf("%+v", pl)

	// get ref type: branch or tag and build
	if pl.RefType == "branch" || pl.RefType == "tag" {
		ref := pl.Ref
		cloneURL := pl.Repository.CloneURL
		// build
		ctx, cancel := context.WithCancel(context.Background())
		gCancelCommands = append(gCancelCommands, cancel)
		go gClient.BuildPushImage(ctx, cloneURL, ref, pl.Repository.Name, pl.Repository.FullName, ref, gNotify)
	}
}
