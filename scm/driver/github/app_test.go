// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jenkins-x/go-scm/scm"
	"gopkg.in/h2non/gock.v1"
)

func TestAppRepositoryInstallation(t *testing.T) {
	defer gock.Off()

	gock.New("https://api.github.com").
		Get("/repos/octocat/hello-world/installation").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/app_repo_install.json")

	client := NewDefault()
	got, _, err := client.Apps.GetRepositoryInstallation(
		context.Background(),
		"octocat/hello-world",
	)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Installation)
	raw, _ := os.ReadFile("testdata/app_repo_install.json.golden")
	err = json.Unmarshal(raw, want)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)

		data, err := json.Marshal(got)
		if err != nil {
			t.Logf("got error marshalling response to JSON: %s", err.Error())
		} else {
			t.Logf("got JSON: %s", string(data))
		}
	}
}
