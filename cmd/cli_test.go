package main

import (
	"github.com/rendon/testcli"
	"github.com/statping/statping/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
	"time"
)

var (
	dir string
)

func init() {
	dir = utils.Directory
	//core.SampleHits = 480
}

func TestStartServerCommand(t *testing.T) {
	t.SkipNow()
	os.Setenv("DB_CONN", "sqlite")
	cmd := helperCommand(nil, "")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(60*time.Second), got)
	os.Unsetenv("DB_CONN")
	gg, _ := <-got
	assert.Contains(t, gg, "DB_CONN environment variable was found")
	assert.Contains(t, gg, "Core database does not exist, creating now!")
	assert.Contains(t, gg, "Starting monitoring process for 5 Services")
}

func TestVersionCommand(t *testing.T) {
	c := testcli.Command("statping", "version")
	c.Run()
	assert.True(t, c.StdoutContains(VERSION))
}

func TestHelpCommand(t *testing.T) {
	c := testcli.Command("statping", "help")
	c.Run()
	t.Log(c.Stdout())
	assert.True(t, c.StdoutContains("statping help               - Shows the user basic information about Statping"))
}

func TestStaticCommand(t *testing.T) {
	t.SkipNow()
	cmd := helperCommand(nil, "static")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(10*time.Second), got)
	gg, _ := <-got
	t.Log(gg)
	assert.Contains(t, gg, "Exporting Static 'index.html' page...")
	assert.Contains(t, gg, "Exported Statping index page: 'index.html'")
	assert.FileExists(t, dir+"/index.html")
}

func TestExportCommand(t *testing.T) {
	t.SkipNow()
	cmd := helperCommand(nil, "export")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(10*time.Second), got)
	gg, _ := <-got
	t.Log(gg)
	assert.FileExists(t, dir+"/statping-export.json")
}

func TestUpdateCommand(t *testing.T) {
	t.SkipNow()
	cmd := helperCommand(nil, "version")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(15*time.Second), got)
	gg, _ := <-got
	t.Log(gg)
	assert.Contains(t, gg, VERSION)
}

func TestAssetsCommand(t *testing.T) {
	t.SkipNow()
	c := testcli.Command("statping", "assets")
	c.Run()
	t.Log(c.Stdout())
	t.Log("Directory for Assets: ", dir)
	time.Sleep(1 * time.Second)
	err := utils.DeleteDirectory(dir + "/assets")
	require.Nil(t, err)
	assert.FileExists(t, dir+"/assets/robots.txt")
	assert.FileExists(t, dir+"/assets/scss/base.scss")
	assert.FileExists(t, dir+"/assets/scss/main.scss")
	assert.FileExists(t, dir+"/assets/scss/variables.scss")
	assert.FileExists(t, dir+"/assets/css/main.css")
	assert.FileExists(t, dir+"/assets/css/vendor.css")
	assert.FileExists(t, dir+"/assets/css/style.css")
	err = utils.DeleteDirectory(dir + "/assets")
	require.Nil(t, err)
}

func TestRunCommand(t *testing.T) {
	t.SkipNow()
	cmd := helperCommand(nil, "run")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(15*time.Second), got)
	gg, _ := <-got
	t.Log(gg)
	assert.Contains(t, gg, "Running 1 time and saving to database...")
	assert.Contains(t, gg, "Check is complete.")
}

func TestEnvironmentVarsCommand(t *testing.T) {
	c := testcli.Command("statping", "env")
	c.Run()
	assert.True(t, c.StdoutContains("Statping Environment Variable"))
}

func TestVersionCLI(t *testing.T) {
	run := catchCLI([]string{"version"})
	assert.EqualError(t, run, "end")
}

func TestAssetsCLI(t *testing.T) {
	catchCLI([]string{"assets"})
	//assert.EqualError(t, run, "end")
	assert.FileExists(t, dir+"/assets/css/main.css")
	assert.FileExists(t, dir+"/assets/css/style.css")
	assert.FileExists(t, dir+"/assets/css/vendor.css")
	assert.FileExists(t, dir+"/assets/scss/base.scss")
	assert.FileExists(t, dir+"/assets/scss/mobile.scss")
	assert.FileExists(t, dir+"/assets/scss/variables.scss")
}

func TestSassCLI(t *testing.T) {
	catchCLI([]string{"sass"})
	assert.FileExists(t, dir+"/assets/css/main.css")
	assert.FileExists(t, dir+"/assets/css/style.css")
	assert.FileExists(t, dir+"/assets/css/vendor.css")
}

func TestUpdateCLI(t *testing.T) {
	t.SkipNow()
	cmd := helperCommand(nil, "update")
	var got = make(chan string)
	commandAndSleep(cmd, time.Duration(15*time.Second), got)
	gg, _ := <-got
	t.Log(gg)
	assert.Contains(t, gg, "version")
}

func TestHelpCLI(t *testing.T) {
	run := catchCLI([]string{"help"})
	assert.EqualError(t, run, "end")
}

func TestRunOnceCLI(t *testing.T) {
	t.SkipNow()
	run := catchCLI([]string{"run"})
	assert.EqualError(t, run, "end")
}

func TestEnvCLI(t *testing.T) {
	run := catchCLI([]string{"env"})
	assert.Error(t, run)
}

func commandAndSleep(cmd *exec.Cmd, duration time.Duration, out chan<- string) {
	go func(out chan<- string) {
		runCommand(cmd, out)
	}(out)
	time.Sleep(duration)
	cmd.Process.Kill()
}

func helperCommand(envs []string, s ...string) *exec.Cmd {
	cmd := exec.Command("statping", s...)
	return cmd
}

func runCommand(c *exec.Cmd, out chan<- string) {
	bout, _ := c.CombinedOutput()
	out <- string(bout)
}
