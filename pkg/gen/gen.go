package gen

import (
	"encoding/json"
	"fmt"
	"grbac-gen/pkg/parser"
	"grbac-gen/pkg/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type Gen struct {
	json func(data interface{}) ([]byte, error)
	yaml func(data interface{}) ([]byte, error)
}

type Config struct {
	SearchDir    string
	OutputDir    string
	ParseDepth   int
	ExcludeFiles []string
	OutputFile   string
	Format       string
	Tag          string
	SsRole       string
	Tidy         bool
}

func New() *Gen {
	return &Gen{
		json: func(data interface{}) ([]byte, error) {
			return json.MarshalIndent(data, "", "    ")
		},
		yaml: func(data interface{}) ([]byte, error) {
			return nil, nil
			// return yml.Unmarshal()
		},
	}
}

func (g *Gen) Build(config *Config) error {
	if _, err := os.Stat(config.SearchDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s does not exist", config.SearchDir)
	}
	if _, err := os.Stat(config.OutputDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s does not exist", config.OutputDir)
	}

	rootPkgName, err := getPkgName(config.SearchDir)
	if err != nil {
		log.Fatalln("Failed get pkg name: " + err.Error())
	}

	p := parser.New(
		parser.SetExcludeFiles(config.ExcludeFiles),
		parser.SetTag(config.Tag),
		parser.SetSsrole(config.SsRole),
	)

	log.Println("Parse all go files...")
	if err = p.ParseAllGoFiles(rootPkgName, config.SearchDir); err != nil {
		return err
	}

	log.Println("Gen permissions...")
	if err = p.GenPermissions(); err != nil {
		return err
	}

	permissions := p.GetPermissions()

	sort.Sort(parser.PS(permissions))

	log.Println("Output file...")
	fullfile, err := g.Output(config, permissions)
	if err != nil {
		return err
	}
	log.Println("file: " + fullfile)

	return nil
}

func (g *Gen) Output(config *Config, data []*parser.Permission) (string, error) {
	var (
		byts []byte
		err  error
	)
	if config.SsRole != "" { // TODO: 验证是否可以通过一个 permission 添加超级管理员（所有的维度都是星号）
		for _, p := range data {
			if p.PermissionDoc == nil {
				continue
			}
			if !utils.Contains(p.PermissionDoc.AuthorizedRoles, config.SsRole) {
				p.PermissionDoc.AuthorizedRoles = append(p.PermissionDoc.AuthorizedRoles, config.SsRole)
			}
			// 如果设置了 ssRole，不再使用星号
			p.PermissionDoc.AuthorizedRoles = utils.Remove(p.PermissionDoc.AuthorizedRoles, "*")
		}
	}
	if config.Tidy {

	}
	filename := ""
	switch config.Format {
	case "json":
		filename = config.OutputFile + ".json"
		byts, err = g.json(data)
	case "yaml":
		// TODO: ...
		fmt.Println("yaml TODO: ...")
		os.Exit(1)
	}

	dir, err := filepath.Abs(config.OutputDir)
	if err != nil {
		fmt.Println("------------ err output dir")
		os.Exit(1)
	}
	fullDir := filepath.Join(dir, filename)
	f, err := os.Create(fullDir)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.Write(byts)
	return fullDir, err
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
