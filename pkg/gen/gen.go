package gen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"grbac-gen/pkg/parser"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

type Gen struct {
	json func(data interface{}) ([]byte, error)
}

type Config struct {
	SearchDir  string
	OutputDir  string
	ParseDepth int
	MainFile   string
	OutputVar  string
}

func New() *Gen {
	return &Gen{
		json: func(data interface{}) ([]byte, error) {
			return json.MarshalIndent(data, "", "    ")
		},
	}
}

func (g *Gen) Build(config *Config) error {
	// check dir exists
	if _, err := os.Stat(config.SearchDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s does not exist", config.SearchDir)
	}
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s does not exist", config.OutputDir)
	}

	// get pkg name
	pkgName, err := getPkgName(config.SearchDir)
	if err != nil {
		log.Fatalln("Failed get pkg name: " + err.Error())
	}
	log.Println("[DEBUG] get pkg name is", pkgName)

	p := parser.New()
	log.Println("Parsing all go files...")
	if err = p.ParseAllGoFiles(pkgName, config.SearchDir); err != nil {
		return err
	}

	log.Println("Gen permissions...")
	if err = p.GenPermissions(); err != nil {
		return err
	}

	log.Println("Output json...")
	content := new(bytes.Buffer)
	if content, err = p.GetPermissions(); err != nil {
		return err
	}

	arr := strings.Split(config.OutputDir, "/")
	fileConfig := &FileConfig{
		PkgName:   arr[len(arr)-1],
		OutputVar: config.OutputVar,
		Content:   content,
	}
	if err = g.Output(fileConfig); err != nil {
		return err
	}

	return nil
}

func (g *Gen) Output(fileConf *FileConfig) error {
	if g.json == nil {
		return fmt.Errorf("nil Parser.json")
	}

	arr := strings.Split(fileConf.PkgName, "/")
	fileConf.PkgName = arr[len(arr)-1]

	f, err := os.Create("permission")
	if err != nil {
		return err
	}

	t, err := template.New("permission_info").Parse(tempStr)
	if err != nil {
		return err
	}

	return t.ExecuteTemplate(f, "permission_info", fileConf)
}

func getPkgName(searchDir string) (string, error) {
	cmd := exec.Command("go", "list", "-f={{.ImportPath}}")
	cmd.Dir = searchDir
	var stdOut, stdErr strings.Builder
	cmd.Stdout = &stdOut // 用自定义的 io.Writer 顶替 stdOut 和 stdErr 接收数据
	cmd.Stderr = &stdErr //

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("execute go list err. err: %s, stdOut: %s, stdErr: %s",
			err.Error(), stdOut.String(), stdErr.String())
	}
	outStr := stdOut.String()

	res := strings.Split(outStr, "\n")
	return res[0], nil
}
