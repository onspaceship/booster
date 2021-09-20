package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onspaceship/booster/pkg/client"
	"github.com/onspaceship/booster/pkg/config"
	"github.com/onspaceship/booster/pkg/controller"

	"github.com/apex/log"
	kpackapi "github.com/pivotal/kpack/pkg/apis/build/v1alpha1"
	kpackclient "github.com/pivotal/kpack/pkg/client/clientset/versioned/typed/build/v1alpha1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

type buildPayload struct {
	AppId        string `json:"app_id"`
	BuildId      string `json:"build_id"`
	TeamHandle   string `json:"team_handle"`
	AppHandle    string `json:"app_handle"`
	GitRef       string `json:"git_ref"`
	GitSHA       string `json:"git_sha"`
	BuilderImage string `json:"builder_image"`
}

func handleBuild(jsonPayload []byte, options *config.SocketOptions) {
	var payload buildPayload
	err := json.Unmarshal(jsonPayload, &payload)
	if err != nil {
		log.WithError(err).Info("Payload is invalid")
		return
	}

	githubCreds, err := client.NewClient().CoreAppGithubCredentials(payload.AppId)
	if err != nil {
		log.WithError(err).Info("Couldn't fetch GitHub credentials")
		return
	}

	log.WithField("payload", payload).Info("Handling build")

	name := fmt.Sprintf("%s-%s-%s", payload.TeamHandle, payload.AppHandle, payload.BuildId)
	tags := []string{fmt.Sprintf("%s/%s/%s:%s", options.BuildRegistryURL, payload.TeamHandle, payload.AppHandle, payload.GitSHA)}

	if strings.HasPrefix(payload.GitRef, "refs/") {
		refParts := strings.Split(payload.GitRef, "/")
		tags = append(tags, fmt.Sprintf("%s/%s/%s:%s", options.BuildRegistryURL, payload.TeamHandle, payload.AppHandle, refParts[len(refParts)-1]))
	}

	gitURL := fmt.Sprintf("https://x-access-token:%s@github.com/%s.git", githubCreds.Token, githubCreds.RepoPath)

	client := kpackclient.NewForConfigOrDie(ctrl.GetConfigOrDie())

	_, err = client.Builds(options.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		log.WithError(err).Info("Couldn't fetch Build from Kubernetes")
		return
	}

	if err == nil {
		log.WithError(err).WithField("app-id", payload.AppId).WithField("build-id", payload.BuildId).Info("Build already exists")
		return
	}

	build := &kpackapi.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				controller.AppIdAnnotation:   payload.AppId,
				controller.BuildIdAnnotation: payload.BuildId,
				controller.GitRefAnnotation:  payload.GitRef,
			},
		},
		Spec: kpackapi.BuildSpec{
			Tags: tags,
			Source: kpackapi.SourceConfig{
				Git: &kpackapi.Git{
					URL:      gitURL,
					Revision: payload.GitSHA,
				},
			},
			Builder: kpackapi.BuildBuilderSpec{
				Image: payload.BuilderImage,
			},
			ServiceAccount: options.BuildServiceAccount,
		},
	}

	_, err = client.Builds(options.Namespace).Create(context.Background(), build, metav1.CreateOptions{})
	if err != nil {
		log.WithError(err).WithField("app-id", payload.AppId).WithField("build-id", payload.BuildId).Infof("Error creating Build: %v", err)
		return
	}
}

func init() {
	handlers["build"] = handleBuild
}
