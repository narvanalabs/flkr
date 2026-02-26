package parser

import (
	"encoding/xml"
	"io/fs"
)

// PomXML represents a Maven pom.xml file.
type PomXML struct {
	XMLName    xml.Name       `xml:"project"`
	GroupID    string         `xml:"groupId"`
	ArtifactID string        `xml:"artifactId"`
	Properties PomProperties  `xml:"properties"`
	Dependencies struct {
		Dependency []PomDependency `xml:"dependency"`
	} `xml:"dependencies"`
	Parent struct {
		GroupID    string `xml:"groupId"`
		ArtifactID string `xml:"artifactId"`
	} `xml:"parent"`
}

// PomProperties holds Maven properties.
type PomProperties struct {
	JavaVersion string `xml:"java.version"`
	MavenCompilerSource string `xml:"maven.compiler.source"`
}

// PomDependency represents a Maven dependency.
type PomDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
}

// HasDep checks if a Maven dependency exists by groupId:artifactId.
func (p *PomXML) HasDep(groupID, artifactID string) bool {
	for _, d := range p.Dependencies.Dependency {
		if d.GroupID == groupID && d.ArtifactID == artifactID {
			return true
		}
	}
	return false
}

// IsSpringBoot checks if this is a Spring Boot project.
func (p *PomXML) IsSpringBoot() bool {
	if p.Parent.GroupID == "org.springframework.boot" {
		return true
	}
	for _, d := range p.Dependencies.Dependency {
		if d.GroupID == "org.springframework.boot" {
			return true
		}
	}
	return false
}

// JavaVersion returns the configured Java version.
func (p *PomXML) JavaVersion() string {
	if p.Properties.JavaVersion != "" {
		return p.Properties.JavaVersion
	}
	return p.Properties.MavenCompilerSource
}

// ParsePomXML reads and parses a pom.xml.
func ParsePomXML(root fs.FS, path string) (*PomXML, error) {
	data, err := fs.ReadFile(root, path)
	if err != nil {
		return nil, err
	}
	var pom PomXML
	if err := xml.Unmarshal(data, &pom); err != nil {
		return nil, err
	}
	return &pom, nil
}
