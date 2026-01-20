package app

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type App struct {
	Name    string
	Version string
	ModTime time.Time

	// 项目根目录
	RootDir     string
	TemplateDir string

	// 启动时间
	LaunchTime time.Time
	Uptime     time.Duration

	Year int

	Domain string

	Build struct {
		GitCommitLog string
		BuildTime    string
		GitRelease   string
		GoVersion    string
		GRPCVersion  string
	}

	PID int

	Env string

	locker sync.Mutex
}

func (a *App) SetUptime() {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.Uptime = time.Since(a.LaunchTime)
}

func (a *App) IsDev() bool {
	return a.Env == "dev"
}

func (a *App) IsPro() bool {
	return a.Env == "pro"
}

func (a *App) FillBuildInfo(gitCommitLog, buildTime, gitRelease, version string) {
	a.Build.GitCommitLog = gitCommitLog
	a.Build.BuildTime = buildTime
	a.Version = version

	pos := strings.Index(gitRelease, "/")
	if pos >= -1 {
		a.Build.GitRelease = gitRelease[pos+1:]
	}
}

func (a *App) String() string {
	return "Build Info:" +
		"\n\tGit Commit Log: " + a.Build.GitCommitLog +
		"\n\tGit Release Info: " + a.Build.GitRelease +
		"\n\tBuild Time: " + a.Build.BuildTime +
		"\n\tGo Version: " + a.Build.GoVersion +
		"\n\tgRPC Version: " + a.Build.GRPCVersion +
		"\n\tLaunch Time: " + a.LaunchTime.String() +
		"\n\tApp Version: " + a.Version +
		"\n\tEnviroment: " + a.Env +
		"\n\tPID: " + strconv.Itoa(a.PID)
}
