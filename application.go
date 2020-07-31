package console

import (
    "errors"
    "fmt"
    "github.com/mix-go/bean"
    "github.com/mix-go/console/argv"
    "reflect"
)

var (
    // 全局APP
    App *Application
    // 版本号
    Version = "1.0.0-alpha";
    // 最后的错误
    LastError interface{}
)

// 上下文
func Context() *bean.ApplicationContext {
    return App.Context
}

// 创建App
func NewApplication(definition ApplicationDefinition) *Application {
    app := &Application{
        ApplicationDefinition: definition,
    }
    app.Init();

    // 保存指针
    App = app

    return app
}

// App 定义
type ApplicationDefinition struct {
    // 应用名称
    AppName string
    // 应用版本
    AppVersion string
    // 应用调试
    AppDebug bool
    // 依赖配置
    Beans []bean.BeanDefinition
    // 命令集合
    Commands []CommandDefinition
}

// App
type Application struct {
    // App 定义
    ApplicationDefinition
    // 基础路径
    BasePath string
    // 应用上下文
    Context *bean.ApplicationContext
    // 是否单命令
    Singleton bool
}

// 命令定义
type CommandDefinition struct {
    // 命令名称
    Name string
    // 使用描述
    Usage string
    // 选项
    Options []OptionDefinition
    // 反射
    Reflect func() reflect.Value
    // 是否单命令
    Singleton bool
}

// 命令选项
type OptionDefinition struct {
    Names []string
    Usage string
}

// 初始化
func (t *Application) Init() {
    t.Context = bean.NewApplicationContext(t.Beans)
    t.BasePath = argv.Program.Dir
}

// 执行
func (t *Application) Run() {
    if t.AppDebug {
        defer func() {
            if err := recover(); err != nil {
                LastError = err
                fmt.Println(err)
            }
        }()
    }



    if len(t.Commands) == 0 {
        panic(errors.New("Command cannot be empty"))
    }

    // 提取命令
    var cmd *CommandDefinition
    cmdName := argv.Command
    if cmdName == "" {
        // 单命令
        for _, c := range t.Commands {
            if c.Singleton {
                cmd = &c
                t.Singleton = true
                break
            }
        }
        if cmd == nil {
            panic(errors.New("Singleton command not found"))
        }
    } else {
        for _, c := range t.Commands {
            if c.Name == cmdName {
                cmd = &c
                break
            }
        }
    }
    if cmd == nil {
        panic(errors.New(fmt.Sprintf("'%s' is not command, see '%s --help'.", cmdName, argv.Program.Path)))
    }

    // 执行命令
    v := cmd.Reflect()
    m := v.MethodByName("Main")
    if !m.IsValid() {
        panic(errors.New(fmt.Sprintf("'%s' Main method not found", fmt.Sprintf("%#v", v))))
    }
    m.Call([]reflect.Value{})
}

func (t *Application) help() {
    program := argv.Program.Path
    flag := ""
    if t.Singleton {
        flag = " [OPTIONS] COMMAND"
    }
    fmt.Println(fmt.Sprintf("Usage: %s%s [opt...]", program, flag))
    t.printGlobalOptions()
    if t.Singleton {
        t.printCommands()
    } else {
        t.printCommandOptions()
    }
    fmt.Println("")
    flag = ""
    if t.Singleton {
        flag = " COMMAND"
    }
    fmt.Println("")
    fmt.Println(fmt.Sprintf("Run '%s%s  --help' for more information on a command.", program, flag))
    fmt.Println("")
    fmt.Println("Developed with Mix Go framework. (openmix.org/mix-go)")
}

func (t *Application) printGlobalOptions() {

}

func (t *Application) printCommands() {

}

func (t *Application) printCommandOptions() {

}
