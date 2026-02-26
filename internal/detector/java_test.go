package detector

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/narvanalabs/flkr/pkg/flkr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJavaDetector_SpringBoot(t *testing.T) {
	fsys := fstest.MapFS{
		"pom.xml": &fstest.MapFile{
			Data: []byte(`<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <parent>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-parent</artifactId>
    </parent>
    <properties>
        <java.version>21</java.version>
    </properties>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
    </dependencies>
</project>`),
		},
	}

	d := &JavaDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.LangJava, profile.Language)
	assert.Equal(t, flkr.PkgMaven, profile.PackageManager)
	assert.Equal(t, flkr.FrameworkSpring, profile.Framework)
	assert.Equal(t, "21", profile.Version)
}

func TestJavaDetector_Gradle(t *testing.T) {
	fsys := fstest.MapFS{
		"build.gradle": &fstest.MapFile{
			Data: []byte(`plugins { id 'java' }`),
		},
	}

	d := &JavaDetector{}
	profile, matched, err := d.Detect(context.Background(), fsys)
	require.NoError(t, err)
	assert.True(t, matched)
	assert.Equal(t, flkr.PkgGradle, profile.PackageManager)
}
