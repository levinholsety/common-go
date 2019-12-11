package commutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/commio"
)

// JDKTool provides functions of JDK.
type JDKTool struct {
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

// NewJDKTool creates a JDKTool instance.
func NewJDKTool(javaHome string) *JDKTool {
	jdk := JDKTool{javaHome: javaHome}
	jdk.javaPath = jdk.getBinaryFilePath(java)
	jdk.javacPath = jdk.getBinaryFilePath(javac)
	jdk.javadocPath = jdk.getBinaryFilePath(javadoc)
	jdk.jarPath = jdk.getBinaryFilePath(jar)
	return &jdk
}

func (jdk *JDKTool) getBinaryFilePath(name string) string {
	if comm.IsWindows {
		name += ".exe"
	}
	if len(jdk.javaHome) > 0 {
		name = filepath.Join(jdk.javaHome, "bin", name)
	}
	return name
}

// JavaVersion prints java version.
func (jdk *JDKTool) JavaVersion() error {
	return comm.NewCommand(jdk.javaPath).Execute("-version")
}

// JavacVersion prints javac version.
func (jdk *JDKTool) JavacVersion() error {
	return comm.NewCommand(jdk.javacPath).Execute("-version")
}

// JavaC invokes javac command.
func (jdk *JDKTool) JavaC(sourcePath string, classPaths []string, binDir string, javaVersion string) error {
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
	return jdk.javac(sourcePath, binDir, args)
}

func (jdk *JDKTool) javac(dir string, binDir string, args []string) error {
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
			err = jdk.javac(srcPath, destPath, args)
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
		return comm.NewCommand(jdk.javacPath).Execute(append(args, path)...)
	}
	return nil
}

// Jar invokes jar command.
func (jdk *JDKTool) Jar(dir string, jarFile string, overwrite bool) error {
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
	return comm.NewCommand(jdk.jarPath).Execute(args...)
}

// JavaDoc invokes javadoc command.
func (jdk *JDKTool) JavaDoc(sourcePaths []string, classPaths []string, packageName string, version string, docPath string, windowTitle string) error {
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
	return comm.NewCommand(jdk.javadocPath).Execute(args...)
}
