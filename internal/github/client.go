package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	baseURL = "https://api.github.com/repos/godotengine/godot-builds/releases"
	perPage = 50
)

type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Size        int    `json:"size"`
}

type Release struct {
	TagName    string  `json:"tag_name"`
	Name       string  `json:"name"`
	Draft      bool    `json:"draft"`
	Prerelease bool    `json:"prerelease"`
	Assets     []Asset `json:"assets"`
}

func (r Release) HasAssetWithSuffix(suffix string) bool {
	for _, asset := range r.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			return true
		}
	}
	return false
}

func (r Release) FindAssetBySuffix(suffix string) (Asset, bool) {
	for _, asset := range r.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			return asset, true
		}
	}
	return Asset{}, false
}

func FetchReleases(includePrerelease bool) ([]Release, error) {
	var allReleases []Release
	page := 1

	for {
		url := fmt.Sprintf("%s?per_page=%d&page=%d", baseURL, perPage, page)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to reach Github API: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusForbidden {
			return nil, fmt.Errorf("Github API rate limit exceeded, please wait a moment and try again")
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("Github API returned unexpected status: %s", resp.Status)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read Github API response: %v", err)
		}

		var releases []Release
		err = json.Unmarshal(body, &releases)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Github API response: %v", err)
		}

		if len(releases) == 0 {
			break
		}

		for _, r := range releases {
			if r.Draft {
				continue
			}

			if r.Prerelease && !includePrerelease {
				continue
			}
			allReleases = append(allReleases, r)
		}

		if len(releases) < perPage {
			break
		}

		page++
	}
	return allReleases, nil
}

func FilterByAssetSuffix(releases []Release, suffix string) []Release {
	var filtered []Release
	for _, r := range releases {
		if r.HasAssetWithSuffix(suffix) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
