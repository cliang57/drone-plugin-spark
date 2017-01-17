package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	urlSparkMessage = "https://api.ciscospark.com/v1/messages"
	urlSparkRoom    = "https://api.ciscospark.com/v1/rooms"
	urlSparkPeople  = "https://api.ciscospark.com/v1/people"
)

type (
	Repo struct {
		Owner    string
		Name     string
		FullName string
	}

	Build struct {
		Tag        string
		Event      string
		Number     int
		Commit     string
		Ref        string
		Branch     string
		Author     string
		Email      string
		Status     string
		Link       string
		CommitLink string
		Message    string
		DroneLink  string
		Started    int64
		Created    int64
	}

	Config struct {
		Message   string
		AuthToken string
		RoomName  string
		RoomID    string
	}

	Job struct {
		Started int64
	}

	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

// SparkRoomListRsp - Response Struct for Spark Room List API
type SparkRoomListRsp struct {
	Rooms []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Type  string `json:"type"` // "direct" or "group"
	} `json:"items"`
}

func (p Plugin) Exec() error {

	var err error

	roomID := ""

	// Get Room ID
	if p.Config.RoomID == "" && p.Config.RoomName == "" {
		msg := fmt.Sprintf("Must Specify roomId or roomName.")
		return errors.New(msg)
	} else if p.Config.RoomID != "" {
		roomID = p.Config.RoomID
	} else {
		roomID, err = p.roomNameToRoomID(p.Config.RoomName)
		if err != nil {
			return err
		}
	}

	// POST Build Message to Spark
	msg := message(p.Repo, p.Build)
	err = p.sendMessage(roomID, msg)
	if err != nil {
		return err
	}

	// POST Additional Message to Spark
	if p.Config.Message != "" {
		err = p.sendMessage(roomID, p.Config.Message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Plugin) sendMessage(roomID string, msg string) error {

	payloadRaw := map[string]string{}
	payloadRaw["roomId"] = roomID
	payloadRaw["markdown"] = msg
	payloadJSONBytes, _ := json.Marshal(payloadRaw)

	req, err := http.NewRequest("POST", urlSparkMessage, bytes.NewBuffer(payloadJSONBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.Config.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return nil
}

func (p Plugin) roomNameToRoomID(roomName string) (string, error) {

	roomID := ""

	// Query Rooms from Spark (with Spark User Token)
	req, err := http.NewRequest("GET", urlSparkRoom, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+p.Config.AuthToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	// Parse Room Query Result
	m := SparkRoomListRsp{}
	err = json.Unmarshal(content, &m)
	if err != nil {
		msg := fmt.Sprintf("Failed to Query Rooms from Spark ('%s'). Bad JSON Response ('%s').", urlSparkRoom, string(content))
		return "", errors.New(msg)
	}

	// March Room Name with Room ID
	for _, room := range m.Rooms {
		if room.Title == roomName {
			roomID = room.ID
		}
	}
	if roomID == "" {
		msg := fmt.Sprintf("No Room Name Match for Spark Token User.")
		return "", errors.New(msg)
	}

	return roomID, nil
}

func message(repo Repo, build Build) string {

	msg := ""

	if build.Status == "success" {
		msg = msg + fmt.Sprintf("##Build for %s is Successful \n", repo.FullName)
		msg = msg + fmt.Sprintf("**Build author:** [%s](%s) \n", build.Author, build.Email)
	} else {
		msg = msg + fmt.Sprintf("#Build for %s is FAILED!!! \n", repo.FullName)
		msg = msg + fmt.Sprintf("**Drone blames build author:** [%s](%s) \n", build.Author, build.Email)
	}

	msg = msg + "###Build Details \n"
	msg = msg + fmt.Sprintf("* [Build Log](%s)\n", build.Link)
	msg = msg + fmt.Sprintf("* [Commit Log](%s)\n", build.CommitLink)
	msg = msg + fmt.Sprintf("* **Branch:** %s\n", build.Branch)
	msg = msg + fmt.Sprintf("* **Event:** %s\n", build.Event)
	msg = msg + fmt.Sprintf("* **Commit Message:** %s\n", build.Message)

	/* Plan : Need to beautify multi-line commit message. */

	return msg
}
