package cmd

import (
	"context"
	"fmt"
	"github.com/murphysecurity/murphysec/display"
	"github.com/murphysecurity/murphysec/env"
	"github.com/murphysecurity/murphysec/inspector"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/murphysecurity/murphysec/utils/must"
	"github.com/spf13/cobra"
	"path/filepath"
)

var CliJsonOutput bool

var DeepScan bool
var ProjectId string

func scanCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "scan DIR",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			ctx := context.TODO()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败")
					SetGlobalExitCode(1)
					return
				}
			}
			tt := model.TaskTypeCli
			if CliJsonOutput {
				tt = model.TaskTypeJenkins
			}
			if !utils.IsDir(projectDir) {
				tt.UI().Display(display.MsgInfo, "正在为您检测该文件所在的目录")
				projectDir = filepath.Dir(projectDir)
			}
			task := model.CreateScanTask(projectDir, model.TaskKindNormal, tt)
			task.ProjectId = ProjectId
			if env.SpecificProjectName != "" {
				task.ProjectName = env.SpecificProjectName
			}
			task.EnableDeepScan = DeepScan
			ctx = model.WithScanTask(ctx, task)

			if e := inspector.Scan(ctx); e != nil {
				if tt == model.TaskTypeJenkins {
					fmt.Println(model.GenerateIdeaErrorOutput(e))
				}
				SetGlobalExitCode(-1)
			} else {
				if tt == model.TaskTypeJenkins {
					fmt.Println(model.GenerateIdeaOutput(ctx))
				}
			}
		},
		Short: "Scan the source code of the specified project, currently supporting java, javascript, go, and python",
	}
	c.Flags().BoolVar(&CliJsonOutput, "json", false, "json output")
	if env.AllowDeepScan {
		c.Flags().BoolVar(&DeepScan, "deep", false, "deep scan, will upload the source code")
	}
	c.Flags().StringVar(&ProjectId, "project-id", "", "team id")
	must.Must(c.Flags().MarkHidden("project-id"))
	c.Flags().StringVar(&env.SpecificProjectName, "project-name", "", "force specific project name")
	c.Flags().BoolVar(&env.DisableGit, "skip-git", false, "force ignore git info")
	c.Flags().StringVar(&env.Scope, "scope", "", "specify the scope type (only for maven)\ndefault \"compile,runtime\"\nto specify all scopes, use \"all\"\ncan be multiple, but need to be separated by commas")
	c.Flags().StringVar(&env.GradleProjects, "gradle-projects", "", "specify the gradle projects, split by ',', like \"model,:app,:app:guide\"")
	c.Args = cobra.ExactArgs(1)
	return c
}

func binScanCmd() *cobra.Command {
	var jsonOutput bool
	c := &cobra.Command{
		Use: "binscan DIR",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			ctx := context.TODO()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败")
					SetGlobalExitCode(1)
					return
				}
			}
			taskType := model.TaskTypeCli
			if jsonOutput {
				taskType = model.TaskTypeJenkins
			}
			task := model.CreateScanTask(projectDir, model.TaskKindBinary, taskType)
			ctx = model.WithScanTask(ctx, task)
			if e := inspector.BinScan(ctx); e != nil {
				SetGlobalExitCode(1)
			} else {
				if jsonOutput {
					fmt.Println(model.GenerateIdeaOutput(ctx))
				}
			}
		},
		Short: "Scan specified binary files and software artifacts, currently supporting .jar, .war, and common binary file formats (The file will be uploaded to the server for analysis.)",
	}
	c.Flags().BoolVar(&jsonOutput, "json", false, "json output")
	c.Args = cobra.ExactArgs(1)
	return c
}

func iotScanCmd() *cobra.Command {
	var jsonOutput bool
	c := &cobra.Command{
		Use:   "iotscan DIR",
		Short: "Scan the specified IoT device firmware, currently supporting .bin or other formats (The file will be uploaded to the server for analysis.)",
		Run: func(cmd *cobra.Command, args []string) {
			initConsoleLoggerOrExit()
			ctx := context.TODO()
			projectDir := args[0]
			var e error
			if !filepath.IsAbs(projectDir) {
				projectDir, e = filepath.Abs(projectDir)
				if e != nil || !utils.IsPathExist(projectDir) {
					fmt.Println("读取路径失败")
					SetGlobalExitCode(1)
					return
				}
			}
			taskType := model.TaskTypeCli
			if jsonOutput {
				taskType = model.TaskTypeJenkins
			}
			task := model.CreateScanTask(projectDir, model.TaskKindIotScan, taskType)
			ctx = model.WithScanTask(ctx, task)
			if e := inspector.BinScan(ctx); e != nil {
				SetGlobalExitCode(1)
			} else {
				if jsonOutput {
					fmt.Println(model.GenerateIdeaOutput(ctx))
				}
			}
		},
	}
	c.Args = cobra.ExactArgs(1)
	c.Flags().BoolVar(&jsonOutput, "json", false, "json output")
	return c
}
