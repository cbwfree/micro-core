// 编译版本信息
package compile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/micro/go-micro/v2"
	"os"
	"runtime"
	"strings"
	"text/template"
)

// set by build LD_FLAGS
var (
	name      = "ServiceName" // 默认服务名称
	version   = "latest"      // 默认版本号
	revision  = ""            // Revision
	buildTime = ""            // 默认编译时间
)

func SetName(n string) {
	name = n
}

func Name() string {
	return name
}

func Version() string {
	return version
}

func Revision() string {
	return revision
}

func FullVersion() string {
	return fmt.Sprintf("%s.%s", version, revision)
}

func BuildTime() string {
	return buildTime
}

var (
	borderWidth = 80
	sideWidth   = 6
	split       = "#"
	border      = strings.Repeat(split, borderWidth)
	side        = strings.Repeat(split, sideWidth)
)

var verTemplate = ` SrvName:        {{.SrvName}}
 Version:        {{.Version}}
 Revision:       {{.Revision}}
 BuildTime:      {{.BuildTime}}
 GoVersion:      {{.GoVersion}}
 OS/Arch:        {{.Os}}/{{.Arch}}
`
var infoTemplate = ` NodeId:         {{.NodeId}}
 Registry:       {{.Registry}}
 Broker:         {{.Broker}}
 Server:         {{.Server}}
 Client:         {{.Client}}
`

type VerInfo struct {
	SrvName   string `json:"srv_name"`   // 服务名称
	Version   string `json:"version"`    // 版本号
	Revision  string `json:"revision"`   // SVN Revision
	BuildTime string `json:"build_time"` // 构建时间
	GoVersion string `json:"go_version"` // Go版本
	Os        string `json:"os"`         // 平台
	Arch      string `json:"arch"`       // 架构
}

type SrvInfo struct {
	NodeId   string `json:"node_id"` // 节点ID
	Registry string `json:"registry"`
	Broker   string `json:"broker"`
	Server   string `json:"server"`
	Client   string `json:"client"`
}

// GetSrvInfo 获取服务信息
func GetSrvInfo(srv micro.Service) []string {
	var buf = new(bytes.Buffer)
	data := &SrvInfo{
		NodeId:   srv.Server().Options().Id,
		Registry: srv.Options().Registry.String(),
		Broker:   srv.Options().Broker.String(),
		Server:   srv.Options().Server.String(),
		Client:   srv.Options().Client.String(),
	}
	_ = template.Must(template.New("srvInfo").Parse(infoTemplate)).Execute(buf, data)
	return strings.Split(buf.String(), "\n")
}

// GetVersion 获取编译信息
func GetVerInfo() []string {
	var buf = new(bytes.Buffer)
	data := &VerInfo{
		SrvName:   name,
		Version:   version,
		Revision:  revision,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
	_ = template.Must(template.New("verInfo").Parse(verTemplate)).Execute(buf, data)
	return strings.Split(buf.String(), "\n")
}

// BuildLineInfo 生成行消息
func BuildLineInfo(lines []string) []byte {
	var buf = new(bytes.Buffer)
	for _, str := range lines {
		if str == "" {
			continue
		}
		space := borderWidth - len(str) - 4 - (2 * sideWidth)
		buf.WriteString(side)
		buf.WriteString("    ")
		buf.WriteString(str)
		buf.WriteString(strings.Repeat(" ", space))
		buf.WriteString(side)
		buf.WriteString("\n")
	}
	return buf.Bytes()
}

// GetFormatInfo 获取格式化的信息
func GetFormatInfo(info ...[]byte) []byte {
	var buf = new(bytes.Buffer)
	buf.WriteString(border)
	buf.WriteString("\n")
	for _, b := range info {
		buf.Write(b)
	}
	buf.WriteString(border)
	buf.WriteString("\n")
	return buf.Bytes()
}

// EchoVersion 输出版本信息
func EchoVersion(srv micro.Service) {
	var info [][]byte
	info = append(info, BuildLineInfo(GetVerInfo()))
	if srv != nil {
		info = append(info, BuildLineInfo(GetSrvInfo(srv)))
	}
	_, _ = os.Stdout.Write(GetFormatInfo(info...))
}

// 输出JSON版本信息
func EchoVersionJson() {
	data := &VerInfo{
		SrvName:   name,
		Version:   version,
		Revision:  revision,
		BuildTime: buildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
	}
	bt, _ := json.Marshal(data)
	_, _ = os.Stdout.Write(bt)
	_, _ = os.Stdout.Write([]byte("\n"))
}

func init() {
	cmd := os.Args[1:]
	if len(cmd) > 0 {
		switch cmd[0] {
		case "version": // 格式化输出版本号
			EchoVersion(nil)
			os.Exit(0)
		case "jversion": // JSON格式输出版本信息
			EchoVersionJson()
			os.Exit(0)
		}
	}
}
