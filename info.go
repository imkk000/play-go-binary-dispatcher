package main

import (
	"encoding/json"
	"os/exec"
)

func fetchCompletion(binPath string) (CmdInfo, error) {
	cmd := exec.Command(binPath, "completion", "json")
	out, err := cmd.Output()
	if err != nil {
		return CmdInfo{}, err
	}

	var info CmdInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return CmdInfo{}, err
	}

	return info, nil
}

type CmdInfo struct {
	Name     string    `json:"name"`
	Usage    string    `json:"usage"`
	Commands []CmdInfo `json:"commands,omitempty"`
	Flags    []string  `json:"flags,omitempty"`
}
