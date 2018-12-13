package main

import (
	"encoding/json"
	"io"
	"os"
)

func readJSON(rc io.ReadCloser, data interface{}) error {
	defer rc.Close()
	return json.NewDecoder(rc).Decode(data)
}

type EnvPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Dockerrun struct {
	AWSEBDockerrunVersion int `json:"AWSEBDockerrunVersion"`
	ContainerDefinitions  []*struct {
		Command           []string  `json:"command,omitempty"`
		Essential         bool      `json:"essential"`
		Image             string    `json:"image"`
		MemoryReservation int       `json:"memoryReservation"`
		Name              string    `json:"name"`
		Environment       []EnvPair `json:"environment,omitempty"`
		MountPoints       []struct {
			SourceVolume  string `json:"sourceVolume"`
			ContainerPath string `json:"containerPath"`
			ReadOnly      bool   `json:"readOnly"`
		} `json:"mountPoints"`
		Links        []string `json:"links,omitempty"`
		PortMappings []struct {
			ContainerPort int `json:"containerPort"`
			HostPort      int `json:"hostPort"`
		} `json:"portMappings,omitempty"`
	} `json:"containerDefinitions"`
	Family  string `json:"family"`
	Volumes []struct {
		Name string `json:"name"`
		Host struct {
			SourcePath string `json:"sourcePath"`
		} `json:"host"`
	} `json:"volumes"`
}

func writeJSON(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
