package certgen

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
	"text/template"
        "path/filepath"
        "os"
        "fmt"
        "io"

	"github.com/Azure/acs-engine/pkg/certgen/templates"
	"github.com/Azure/acs-engine/pkg/filesystem"
)

// PrepareMasterFiles creates the shared authentication and encryption secrets
func (c *Config) PrepareMasterFiles() error {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return err
	}
	c.AuthSecret = base64.StdEncoding.EncodeToString(b)

	_, err = rand.Read(b)
	if err != nil {
		return err
	}
	c.EncSecret = base64.StdEncoding.EncodeToString(b)

	return nil
}

// WriteMasterFiles writes the templated master config
func (c *Config) WriteMasterFiles(fs filesystem.Filesystem) error {
	for _, name := range templates.AssetNames() {
		if !strings.HasPrefix(name, "master/") {
			continue
		}
		tb := templates.MustAsset(name)

		t, err := template.New("template").Funcs(template.FuncMap{
			"QuoteMeta": regexp.QuoteMeta,
		}).Parse(string(tb))
		if err != nil {
			return err
		}

		b := &bytes.Buffer{}
		err = t.Execute(b, c)
		if err != nil {
			return err
		}

		err = fs.WriteFile(strings.TrimPrefix(name, "master/"), b.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

	return nil
}


func FilePathWalkDir(root string) ([]string, error) {
    var files []string
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })
    return files, err
}

// WriteDynamicMasterFiles writes the run-time created master config
// Assumed in _output/tmp
func (c *Config) WriteDynamicMasterFiles(fs filesystem.Filesystem) error {
        root := "_output/tmp"
        files, err := FilePathWalkDir(root)
	for _, name := range files {
                fmt.Printf("%v\n",name)
		if !strings.HasPrefix(name, root + "/" +  "master/") {
			continue
		}

                fmt.Printf("Reading %v\n",name);
                b := bytes.NewBuffer(nil);
                f, _ := os.Open(name)
                io.Copy(b,f)
                f.Close()
                short_name := strings.TrimPrefix(name,root + "/" + "master/")
                fmt.Printf("%v\n",short_name);
		err = fs.WriteFile(short_name, b.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

        fmt.Printf("Success\n")
	return nil
}

// WriteNodeFiles writes the templated node config
func (c *Config) WriteNodeFiles(fs filesystem.Filesystem) error {
	for _, name := range templates.AssetNames() {
		if !strings.HasPrefix(name, "node/") {
			continue
		}

		tb := templates.MustAsset(name)

		t, err := template.New("template").Funcs(template.FuncMap{
			"QuoteMeta": regexp.QuoteMeta,
		}).Parse(string(tb))
		if err != nil {
			return err
		}

		b := &bytes.Buffer{}
		err = t.Execute(b, struct{}{})
		if err != nil {
			return err
		}

		err = fs.WriteFile(strings.TrimPrefix(name, "node/"), b.Bytes(), 0666)
		if err != nil {
			return err
		}
	}

	return nil
}
