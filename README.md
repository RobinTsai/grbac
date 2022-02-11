# README

参考 [gin-swagger](https://github.com/swaggo/gin-swagger) 练手项目。 

基于代码注释生成 [grbac](https://github.com/storyicon/grbac) 的配置 json。

## 需求

- [x] 注释添加 tags 字段进行生成，生成时传入 tag 字段，寻找 tags 包含的注释进行生成 
- [ ] 添加选项，支持所有配置中添加一个超级超级管理员角色 (输入字段 SSRole: super-super-role)
- [ ] 根据 url 进行整合，相同路径的用 **
- [ ] 配置扫描目录入口 `dir`
- [ ] 配置排除文件 `excludeFiles`，支持排除多个目录或文件
- [ ] 配置扫描目录层数 `parseDepth`
- [ ] 配置输出目录 `outDir`
- [ ] 配置输出文件名 `output/o`
- [ ] 输出文件格式 `format`，支持 json 或 yaml，默认 json
- [ ] 注释命令格式 `@Attr Arg1 Arg2 ...`（继承自 gin-swagger）
- [ ] 注释属性 ID 为 `@Router /PATH/TO/RESOURCE/URI [METHOD,METHOD]`
- [ ] 可用的注释属性
  - `@AuthRoles` 后接一个参数，类型为字符串，参数内以逗号分隔表示多项
  - `@ForbiddenRoles` 后接一个参数，类型为字符串，参数内以逗号分隔表示多项
  - `@AllowAnyone` 后接一个参数，参数类型为 bool
