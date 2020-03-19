package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/commio"
	"github.com/levinholsety/common-go/util"
)

// JDK provides functions of JDK.
type JDK struct {
	javaHome    string
	javaPath    string
	javacPath   string
	jarPath     string
	javadocPath string
}

const (
	java    = "java"
	javac   = "javac"
	jar     = "jar"
	javadoc = "javadoc"
)

// NewJDK creates a JDK instance.
func NewJDK(javaHome string) *JDK {
	jdk := JDK{javaHome: javaHome}
	jdk.javaPath = jdk.getBinaryFilePath(java)
	jdk.javacPath = jdk.getBinaryFilePath(javac)
	jdk.javadocPath = jdk.getBinaryFilePath(javadoc)
	jdk.jarPath = jdk.getBinaryFilePath(jar)
	return &jdk
}

func (p *JDK) getBinaryFilePath(name string) string {
	if comm.IsWindows {
		name += ".exe"
	}
	if len(p.javaHome) > 0 {
		name = filepath.Join(p.javaHome, "bin", name)
	}
	return name
}

// JavaVersion prints java version.
func (p *JDK) JavaVersion() error {
	return util.NewCommand(p.javaPath).Execute("-version")
}

// JavacVersion prints javac version.
func (p *JDK) JavacVersion() error {
	return util.NewCommand(p.javacPath).Execute("-version")
}

// JavaC invokes javac command.
func (p *JDK) JavaC(sourcePath string, classPaths []string, binDir string, javaVersion string) error {
	var args []string
	args = append(args, "-J-Dfile.encoding=UTF-8")
	args = append(args, "-J-Duser.language=en_US")
	args = append(args, "-Xlint:deprecation")
	args = append(args, "-Xlint:unchecked")
	args = append(args, "-encoding")
	args = append(args, "UTF-8")
	if len(classPaths) > 0 {
		args = append(args, "-cp")
		args = append(args, strings.Join(classPaths, string(os.PathListSeparator)))
	}
	args = append(args, "-sourcepath")
	args = append(args, sourcePath)
	args = append(args, "-d")
	args = append(args, binDir)
	args = append(args, "-source")
	args = append(args, javaVersion)
	args = append(args, "-target")
	args = append(args, javaVersion)
	return p.javac(sourcePath, binDir, args)
}

func (p *JDK) javac(dir string, binDir string, args []string) error {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	compile := false
	for _, info := range infos {
		name := info.Name()
		srcPath := filepath.Join(dir, name)
		destPath := filepath.Join(binDir, name)
		if info.IsDir() {
			err = p.javac(srcPath, destPath, args)
			if err != nil {
				return err
			}
		} else {
			switch filepath.Ext(name) {
			case ".java":
				{
					compile = true
				}
			case ".class":
				{
					continue
				}
			default:
				{
					fmt.Printf("Copying %s into %s\n", name, binDir)
					err = os.MkdirAll(binDir, os.ModePerm)
					if err != nil {
						return err
					}
					_, err = commio.CopyFile(srcPath, binDir)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	if compile {
		path := filepath.Join(dir, "*.java")
		fmt.Printf("Compiling %s\n", path)
		return util.NewCommand(p.javacPath).Execute(append(args, path)...)
	}
	return nil
}

// Jar invokes jar command.
func (p *JDK) Jar(dir string, jarFile string, overwrite bool) error {
	d, _ := filepath.Split(jarFile)
	os.MkdirAll(d, os.ModePerm)
	var args []string
	args = append(args, "-J-Dfile.encoding=UTF-8")
	args = append(args, "-J-Duser.language=en_US")
	if overwrite {
		args = append(args, "cvfM")
	} else {
		args = append(args, "uvfM")
	}
	args = append(args, jarFile)
	args = append(args, "-C")
	args = append(args, dir)
	args = append(args, ".")
	return util.NewCommand(p.jarPath).Execute(args...)
}

// JavaDoc invokes javadoc command.
func (p *JDK) JavaDoc(sourcePaths []string, classPaths []string, packageName string, version string, docPath string, windowTitle string) error {
	var args []string
	args = append(args, "-J-Dfile.encoding=UTF-8")
	args = append(args, "-J-Duser.language=en_US")
	args = append(args, "-locale")
	args = append(args, "en_US")
	args = append(args, "-protected")
	args = append(args, "-sourcepath")
	args = append(args, strings.Join(sourcePaths, string(os.PathListSeparator)))
	if len(classPaths) > 0 {
		args = append(args, "-classpath")
		args = append(args, strings.Join(classPaths, string(os.PathListSeparator)))
	}
	args = append(args, "-subpackages")
	args = append(args, packageName)
	args = append(args, "-source")
	args = append(args, version)
	args = append(args, "-encoding")
	args = append(args, "UTF-8")
	args = append(args, "-d")
	args = append(args, docPath)
	args = append(args, "-version")
	args = append(args, "-author")
	if len(windowTitle) > 0 {
		args = append(args, "-windowtitle")
		args = append(args, windowTitle)
	}
	args = append(args, "-charset")
	args = append(args, "UTF-8")
	args = append(args, "-docencoding")
	args = append(args, "UTF-8")
	return util.NewCommand(p.javadocPath).Execute(args...)
}
