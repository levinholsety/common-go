package utils

import (
	"os"
	"path/filepath"

	"github.com/levinholsety/common-go/comm"
)

// ProGuard provides functions of ProGuard
type ProGuard struct {
	home         string
	proguardPath string
}

// NewProGuard creates a ProGuard instance.
func NewProGuard(home string) *ProGuard {
	tool := ProGuard{home: home, proguardPath: proguard}
	if comm.IsWindows() {
		tool.proguardPath += ".bat"
	} else {
		tool.proguardPath += ".sh"
	}
	if len(home) > 0 {
		tool.proguardPath = filepath.Join(home, "bin", tool.proguardPath)
	}
	return &tool
}

const (
	proguard = "proguard"
)

// Pack packs contents.
func (p *ProGuard) Pack(inJars []string, libJars []string, outJar string) (err error) {
	filePath := outJar + ".pro"
	func() {
		file, err := os.Create(filePath)
		if err != nil {
			return
		}
		defer file.Close()
		tw := NewTextWriter(file)
		for _, inJar := range inJars {
			tw.WriteLine("-injars " + inJar)
		}
		tw.WriteLine("-outjars " + outJar)
		for _, libJar := range libJars {
			tw.WriteLine("-libraryjars " + libJar)
		}
		tw.WriteLine("-dontoptimize")
		tw.WriteLine("-dontnote")
		tw.WriteLine("-keepparameternames")
		tw.WriteLine("-renamesourcefileattribute SourceFile")
		tw.WriteLine("-keepattributes Exceptions,InnerClasses,Signature,Deprecated,SourceFile,LineNumberTable,*Annotation*,EnclosingMethod")
		tw.WriteLine("-keep public class * {")
		tw.WriteLine("    public protected *;")
		tw.WriteLine("}")
		tw.WriteLine("-keepclassmembernames class * {")
		tw.WriteLine("    java.lang.Class class$(java.lang.String);")
		tw.WriteLine("    java.lang.Class class$(java.lang.String, boolean);")
		tw.WriteLine("}")
		tw.WriteLine("-keepclasseswithmembernames,includedescriptorclasses class * {")
		tw.WriteLine("    native <methods>;")
		tw.WriteLine("}")
		tw.WriteLine("-keepclassmembers,allowoptimization enum * {")
		tw.WriteLine("    public static **[] values(); public static ** valueOf(java.lang.String);")
		tw.WriteLine("}")
		tw.WriteLine("-keepclassmembers class * implements java.io.Serializable {")
		tw.WriteLine("    static final long serialVersionUID;")
		tw.WriteLine("    private static final java.io.ObjectStreamField[] serialPersistentFields;")
		tw.WriteLine("    private void writeObject(java.io.ObjectOutputStream);")
		tw.WriteLine("    private void readObject(java.io.ObjectInputStream);")
		tw.WriteLine("    java.lang.Object writeReplace();")
		tw.WriteLine("    java.lang.Object readResolve();")
		tw.WriteLine("}")
	}()
	err = NewCommand(p.proguardPath).Execute("@" + filePath)
	if err != nil {
		return
	}
	err = os.Remove(filePath)
	if err != nil {
		return
	}
	return
}
