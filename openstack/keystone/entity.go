package keystone

import (
	"time"
)

type token struct {
	Value      string
	ExperiesAt time.Time
	ProjectID  string
}

// type auth struct {
// 	Identity identity `json:"identity"`
// }

type authScoped struct {
	Identity identity `json:"identity"`
	Scope    scope    `json:"scope,omitempty"`
}

type identity struct {
	Methods  []string `json:"methods"`
	Password password `json:"password"`
}

type password struct {
	User user `json:"user"`
}

type user struct {
	Name     string `json:"name"`
	Domain   domain `json:"domain"`
	Password string `json:"password"`
}

type scope struct {
	Project project `json:"project"`
	// Domain domain `json:"domain"`
}

type domain struct {
	ID string `json:"id,omitempty"`
}

type project struct {
	Domain domain `json:"domain"`
	Name   string `json:"name"`
}

type resp struct {
	Token struct {
		ExpiresAt time.Time `json:"expires_at"`
		Project   struct {
			Domain struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"domain"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"project"`
		Catalog []struct {
			Endpoints []struct {
				URL       string `json:"url"`
				Interface string `json:"interface"`
				Region    string `json:"region"`
				RegionID  string `json:"region_id"`
				ID        string `json:"id"`
			} `json:"endpoints"`
			Type string `json:"type"`
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"catalog"`
	} `json:"token"`
}

// ProjectResponse list all project type
type ProjectResponse struct {
	Projects []projectData `json:"projects"`
}

// projectData project metadata
type projectData struct {
	IsDomain    bool   `json:"is_domain"`
	Description string `json:"description"`
	Links       struct {
		Self string `json:"self"`
	} `json:"links"`
	Tags     []interface{} `json:"tags"`
	Enabled  bool          `json:"enabled"`
	ID       string        `json:"id"`
	ParentID string        `json:"parent_id"`
	DomainID string        `json:"domain_id"`
	Name     string        `json:"name"`
}
