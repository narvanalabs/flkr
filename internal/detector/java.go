package detector

import (
	"context"
	"io/fs"

	"github.com/narvanalabs/flkr/internal/parser"
	"github.com/narvanalabs/flkr/pkg/flkr"
)

// JavaDetector detects Java applications.
type JavaDetector struct{}

func (d *JavaDetector) Name() string  { return "java" }
func (d *JavaDetector) Priority() int { return 80 }

func (d *JavaDetector) Detect(ctx context.Context, root fs.FS) (*flkr.AppProfile, bool, error) {
	hasPom := fileExists(root, "pom.xml")
	hasGradle := fileExists(root, "build.gradle") || fileExists(root, "build.gradle.kts")

	if !hasPom && !hasGradle {
		return nil, false, nil
	}

	profile := &flkr.AppProfile{
		Language:   flkr.LangJava,
		Confidence: 0.7,
		DetectedBy: d.Name(),
		Port:       8080,
	}

	if hasGradle {
		profile.PackageManager = flkr.PkgGradle
		profile.BuildCommand = "./gradlew build"
		profile.StartCommand = "java -jar build/libs/*.jar"
	}

	if hasPom {
		profile.PackageManager = flkr.PkgMaven
		profile.BuildCommand = "mvn package -DskipTests"
		profile.StartCommand = "java -jar target/*.jar"

		pom, err := parser.ParsePomXML(root, "pom.xml")
		if err == nil {
			if v := pom.JavaVersion(); v != "" {
				profile.Version = v
			}
			if pom.IsSpringBoot() {
				profile.Framework = flkr.FrameworkSpring
				profile.Confidence = 0.9
			}
		}
	}

	return profile, true, nil
}
