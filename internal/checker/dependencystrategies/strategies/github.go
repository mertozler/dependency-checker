package strategies

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/mcuadros/go-version"
	"github.com/mertozler/internal/models"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"unicode"
)

type Github struct {
	GitAPI      string
	RegistryURL string
	Dependency  []string
}

func NewGithub(gitAPI string, registryURL string, Dependency []string) *Github {
	return &Github{
		GitAPI:      gitAPI,
		RegistryURL: registryURL,
		Dependency:  Dependency,
	}
}

func (g *Github) Check(url string) (*models.OutDatedResponse, error) {
	logrus.Infof("Checking %s url with Github Strategy", url)
	repositoryURL := g.createRequestURL(url)
	githubResponse := new(models.GitResponse)

	err := getJSON(repositoryURL, githubResponse)
	if err != nil {
		logrus.Error("Error while getting github response: ", err)
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(githubResponse.Content)
	if err != nil {
		logrus.Error("Error while decoding git content: ", err)
		return nil, err
	}
	githubContent := new(models.GitContentResponse)
	err = json.Unmarshal(data, &githubContent)
	if err != nil {
		logrus.Error("Error while unmarshalling json ", err)
		return nil, err
	}
	outDatedDependencies, err := g.checkDependency(githubContent)
	if err != nil {
		return nil, err
	}
	outDatedDependency := models.OutDatedResponse{RepoURL: url,
		OutdatedDependency: outDatedDependencies}
	return &outDatedDependency, nil
}

func (g *Github) createRequestURL(url string) string {
	logrus.Infof("Generating request url for %s", url)
	var splitedURL = strings.Split(url, "/")
	owner := splitedURL[1]
	repo := splitedURL[2]
	return g.GitAPI + owner + "/" + repo + "/contents/" + g.Dependency[0]
}

func (g *Github) createRegistryRequestURL(dependencyName string) string {
	return g.RegistryURL + dependencyName
}

func getJSON(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}

	if r.StatusCode == 404 {
		return errors.New("This repository was not found, please make sure you entered the correct url")
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func (g *Github) checkDependency(gitContent *models.GitContentResponse) ([]models.OutdatedDependency, error) {
	var outDatedDependenies []models.OutdatedDependency
	for dependencyName, dependencyValue := range gitContent.Dependencies {
		logrus.Infof("Checking dependency for %s", dependencyName)
		url := g.createRegistryRequestURL(dependencyName)
		registryResponse := new(models.RegistiryResponse)
		err := getJSON(url, registryResponse)
		if err != nil {
			return nil, err
		}

		repoDependency, _ := dependencyValue.(string)
		registryDependency := registryResponse.DistTags.Latest

		isInteger := isInt(string(repoDependency[0]))

		if !isInteger {
			replacer := strings.NewReplacer("^", " ", "~", " ")
			repoDependency = replacer.Replace(repoDependency)
		}

		outDatedDependency := compareDependencies(dependencyName, repoDependency, registryDependency)
		outDatedDependenies = append(outDatedDependenies, *outDatedDependency)
	}
	return outDatedDependenies, nil
}

func compareDependencies(dependencyName string, repoDependency string, registryDependency string) *models.OutdatedDependency {
	compareStatus := version.CompareSimple(repoDependency, registryDependency)
	if compareStatus == -1 {
		logrus.Infof("%s is outdated. RepoDependency %s, RegistryDependency %s", dependencyName, repoDependency, registryDependency)
		return &models.OutdatedDependency{
			DependencyName:              dependencyName,
			RegistiryDependencyVersion:  registryDependency,
			RepositoryDependencyVersion: repoDependency,
		}
	}
	return nil
}

func isInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}
