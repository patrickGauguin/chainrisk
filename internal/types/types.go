package types

import "time"

type RepoInfo struct {
	Owner       string
	Name        string
	FullName    string
	Description string
	Stars       int
	Forks       int
	LastPushed  time.Time
	Archived    bool
	Language    string
}

type Dependency struct {
	Name      string
	Version   string
	Ecosystem string
	IsDev     bool
}
