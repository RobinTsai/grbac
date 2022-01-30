package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	attrRouter    = "@router"
	attrAuthRoles = "@authroles"
	attrForbRoles = "@forbiddenroles"
)

type Parser struct {
	packages    map[string]*PackageDefinition
	permissions []*Permission
	json        *strings.Builder
}

type PackageDefinition struct {
	Name  string               // pkgName
	Files map[string]*ast.File // pkgName + relative path => astFile
}

// NewPkgBuild
// @Router       /admin/rent/{rentId}/app-bind-skill [post]
// @AuthRoles	 SuperAdministrator,Administrator

func (p *Parser) Skip(path string, f os.FileInfo) bool {
	if f.IsDir() {
		// 默认跳过文件（夹），不用于 parser 解析
		if f.Name() == "vendor" || f.Name() == "docs" ||
			len(f.Name()) > 1 && f.Name()[0] == '.' {
			return true
		}

		// TODO: add customized skip files
	}
	return false
}

func (p *Parser) CollectAstFile(pkgKey, path string, astFile *ast.File) error {
	if p.packages == nil {
		p.packages = make(map[string]*PackageDefinition)
	}
	if pkgKey == "" {
		return nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	pd, ok := p.packages[pkgKey]
	if ok {
		_, exists := pd.Files[path]
		if exists {
			return nil
		}
		pd.Files[path] = astFile
	} else {
		p.packages[pkgKey] = &PackageDefinition{
			Name:  astFile.Name.Name,
			Files: map[string]*ast.File{path: astFile},
		}
	}

	return nil
}

func (p *Parser) ParseFile(pkgKey, path string, src interface{}) error {
	if strings.HasSuffix(strings.ToLower(path), "_test.go") || filepath.Ext(path) != ".go" {
		return nil
	}

	astFile, err := goparser.ParseFile(token.NewFileSet(), path, src, goparser.ParseComments)
	if err != nil {
		return err
	}

	if err = p.CollectAstFile(pkgKey, path, astFile); err != nil {
		return err
	}

	return nil
}

func New(options ...func(*Parser)) *Parser { // 这种 options 传入的设计方法
	p := &Parser{
		packages:    make(map[string]*PackageDefinition), // NewPackagesDefinitions(),
		permissions: make([]*Permission, 0, 30),
		json:        &strings.Builder{},
	}

	for _, o := range options {
		o(p) // 传入的 option 对 Parser 资源进行处理
	}
	return p
}

func (p *Parser) ParseAllGoFiles(pkgName, searchDir string) error {
	return filepath.Walk(searchDir, func(path string, info fs.FileInfo, err error) error {
		if skipped := p.Skip(pkgName, info); skipped {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(searchDir, path) // 为什么要获取 relative path 呢？是为了
		if err != nil {
			return err
		}

		// 统一 path 的 key
		s1 := filepath.Join(pkgName, relPath) // pkgName 相当于根包路径，再结合相对路径，拼出来目的包的路径
		s2 := filepath.Clean(s1)              // 尽量简化目录
		s3 := filepath.Dir(s2)                // 获取目录
		pathKey := filepath.ToSlash(s3)       // 将不同系统的 Separator 替换成 /

		if err = p.ParseFile(pathKey, path, nil); err != nil {
			return err
		}

		return nil
	})
}

func (p *Parser) GenPermissions() error {
	if p.packages == nil {
		return fmt.Errorf("nil parser.packages")
	}

	for pkgRootPath, definitions := range p.packages {
		for fileRootPath, astFile := range definitions.Files {
			for _, commentGroup := range astFile.Comments {
				comments := strings.Split(commentGroup.Text(), "\n")

				permission := &Permission{
					PermissionDoc: new(PermissionDoc),
					Pkg:           pkgRootPath,
					Filepath:      fileRootPath,
				}
				for _, comment := range comments {
					comment = strings.TrimSpace(comment)
					arr := strings.Split(comment, " ")

					attr := strings.TrimSpace(strings.ToLower(arr[0]))
					info := strings.TrimSpace(strings.Join(arr[1:], " "))
					switch attr {
					case attrRouter:
						permission.RawRouterLine = info
					case attrAuthRoles:
						permission.RawAuthRolesLine = info
					case attrForbRoles:
						permission.RawForbiddenRolesLine = info
					}
				}

				if err := permission.Parse(); err != nil {
					continue
				}

				p.permissions = append(p.permissions, permission)
			}
		}
	}
	return nil
}

func (p *Parser) GetPermissions() (*bytes.Buffer, error) {
	for i, permission := range p.permissions {
		permission.Id = i
		if permission.Host == "" {
			permission.Host = "*"
		}
		if permission.Path == "" {
			permission.Path = "*"
		}
		permission.Method = fmt.Sprintf("{%s}", permission.Method)
	}

	writer := &bytes.Buffer{}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(p.permissions); err != nil {
		return nil, err
	}

	return writer, nil
}
