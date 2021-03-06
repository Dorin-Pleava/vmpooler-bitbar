package commands

import (
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/Dorin-Pleava/vmpooler-bitbar/config"
	"github.com/johnmccabe/go-bitbar"
	"github.com/johnmccabe/go-vmpooler/vm"
	"github.com/spf13/cobra"
)

const logoBase64 = "R0lGODlhIAAgAPQAAP+uGv+uG/+vG/+uHP+vHP+vHf+vHv+vH/+wHv+wH/+wIP+xIf+xIv+xI/+xJf+yJP+zJ/6yKf60LAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACH5BAEAABMALAAAAAAgACAAAAWW4CSOZGmeaKqubIsSQSwTrhoAA0EAQFCngUFEZOj9UAGaKGE8mpIjps9ZuhUOiZuBWuV5AUquqDAoCxIFsDhVDK9LaQKjQWdA3pP0F4DAN/YACX6Agm9/e4Vrh1+JYnpffWsHOAqVlQ6SOIoLnAsPRQNvN3sCeDcEBQVWpmGTU2tQS02wAxJEs2KnBnqvuYC9eMHCwyshADs="

func init() {
	rootCmd.AddCommand(menuCmd)
}

var menuCmd = &cobra.Command{
	Use:    "menu",
	Run:    runMenu,
	Hidden: true,
}

func runMenu(cmd *cobra.Command, args []string) {
	ex, _ := os.Executable()

	cfg, err := config.Read()
	if err != nil {
		plugin := bitbar.New()
		plugin.StatusLine(" ❓").Font("Avenir").Size(16)
		menu := plugin.NewSubMenu()
		if err.Error() == "Config file not found" {
			menu.Line("Initialise config").Bash(ex).Params([]string{"config"}).Terminal(true).Refresh(true)
		} else {
			menu.Line(fmt.Sprintf("Error: %v", err))
		}
		fmt.Print(plugin.Render())
		os.Exit(1)
	}

	vmclient := vm.NewClient(cfg.Endpoint, cfg.Token)

	templates, err := vmclient.ListTemplates()
	if err != nil {
		errorMenu(err)
	}

	virtualmachines, err := vmclient.GetAll()
	if err != nil {
		errorMenu(err)
	}

	plugin := bitbar.New()
	plugin.StatusLine(fmt.Sprintf("VMs: %d", len(virtualmachines))).Color("green")
	menu := plugin.NewSubMenu()
	menu.HR()
	menu.Line("vmpooler").Size(22).Font("Arial Bold").TemplateImage(logoBase64)
	menu.HR()

	if len(virtualmachines) == 0 {
		menu.Line("No running VMs found")
	}

	progressBarStates := []string{"███▏", "██▊▏", "██▋▏", "██▌▏", "██▍▏", "██▎▏", "██▏▏", "██ ▏", "█▉ ▏", "█▊ ▏", "█▋ ▏", "█▌ ▏", "█▍ ▏", "█▎ ▏", "█▏ ▏", "█  ▏", "▉  ▏", "▊  ▏", "▋  ▏", "▌  ▏", "▍  ▏", "▎  ▏", "▏  ▏"}

	for _, vm := range virtualmachines {
		timebar := progressBarStates[int(vm.Running*float64((len(progressBarStates)-1))/float64(vm.Lifetime))]
		vmcolour := "green"
		if (float64(vm.Lifetime) - vm.Running) <= float64(cfg.LifetimeWarning) {
			vmcolour = "red"
		}
		menu.Line(fmt.Sprintf("%s %s (%s)", timebar, vm.Hostname, vm.Template.Id)).
			Color(vmcolour).Font("Menlo-Regular").
			Size(14).
			CopyToClipboard(vm.Fqdn)

		vmmenu := menu.NewSubMenu()

		vmmenu.Line("Action").Font("Arial Bold").Size(14)

		// cmddd := bitbar.Cmd{Bash: "/Users/dorin.pleava/.rvm/gems/ruby-2.5.1/bin/terminal-notifier", Params: []string{"-message", "Hello, this is my message", "-title", "Message Title"}}

		// Puppet Main Menu
		vmmenu.Line("Puppet").
			// Bash(ex).
			Href(fmt.Sprintf("ssh %s@%s -t 'puppet --version'", sshUser(vm), vm.Fqdn)).
			// Params([]string{"info", sshUser(vm), vm.Fqdn}).
			Terminal(true).
			Size(12)

		// exec.Command("osascript -e 'tell application (path to frontmost application as text) to display dialog \"Hello from osxdaily.com\" buttons {\"OK\"} with icon stop'", "")

		puppetmenu := vmmenu.NewSubMenu()

		// Puppet Agent Menu
		puppetmenu.Line("Install Puppet-Agent Version...").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		puppetVersionsMenu := puppetmenu.NewSubMenu()

		puppetVersionsMenu.Line("Latest").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		puppetVersionsMenu.Line("5.5.10").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		// Puppet Enterprise menu
		puppetmenu.Line("Install Puppet-Enterprise Version...").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		puppetEnterpriseVersionsMenu := puppetmenu.NewSubMenu()

		puppetEnterpriseVersionsMenu.Line("Latest").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		puppetEnterpriseVersionsMenu.Line("5.5.10").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

			// Bolt Menu
		puppetmenu.Line("Install Bolt Version...").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		boltVersionsMenu := puppetmenu.NewSubMenu()

		boltVersionsMenu.Line("Latest").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		boltVersionsMenu.Line("3.3.1").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

			// Something new Menu
		puppetmenu.Line("Install Something Version...").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		somethingVersionsMenu := puppetmenu.NewSubMenu()

		somethingVersionsMenu.Line("Latest").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		somethingVersionsMenu.Line("x.y.z").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		vmmenu.Line("SSH to VM").
			Href(fmt.Sprintf("ssh://%s@%s", sshUser(vm), vm.Fqdn)).
			Terminal(true).
			Size(12)

		vmmenu.Line("Delete VM").
			Bash(ex).
			Params([]string{"delete", vm.Hostname}).
			Terminal(false).
			Refresh(true).
			Size(12)

		vmmenu.HR()

		vmmenu.Line("Extend Lifetime (+2h)").
			Bash(ex).
			Params([]string{"extend", vm.Hostname, "TwoHours"}).
			Terminal(false).
			Refresh(true).
			Size(12)

		vmmenu.Line("Extend Lifetime for one day").
			Bash(ex).
			Params([]string{"extend", vm.Hostname, "OneDay"}).
			Terminal(false).
			Refresh(true).
			Size(12)

		vmmenu.Line("Extend Lifetime for three days").
			Bash(ex).
			Params([]string{"extend", vm.Hostname, "ThreeDays"}).
			Terminal(false).
			Refresh(true).
			Size(12)

		vmmenu.HR()

		vmmenu.Line("Status").Font("Arial Bold").Size(14)

		timeText := fmt.Sprintf("%.2f/%.2f hours", vm.Running, float64(vm.Lifetime))
		vmmenu.Line(timeText).
			Color(vmcolour).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(timeText)

		vmmenu.Line(fmt.Sprintf("IP: %s", vm.Ip)).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(vm.Ip)

		vmmenu.HR()

		vmmenu.Line("Template").Font("Arial Bold").Size(14)

		vmmenu.Line(vm.Template.Id).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(vm.Template.Id)

		vmmenu.Line(fmt.Sprintf("OS: %s", vm.Template.Os)).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(vm.Template.Os)

		vmmenu.Line(fmt.Sprintf("Ver: %s", vm.Template.Osver)).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(vm.Template.Osver)

		vmmenu.Line(fmt.Sprintf("Arch: %s", vm.Template.Arch)).
			Font("Menlo-Regular").
			Size(12).
			CopyToClipboard(vm.Template.Arch)

		// TODO get puppet version here or some other info
		// ssh vmpoo.er../v1/./.. puppet --version

	}

	menu.HR()

	menu.Line("Bulk Actions!!!")

	bulkmenu := menu.NewSubMenu()

	bulkmenu.Line("Delete").
		Bash(ex).
		Params([]string{"delete", "all"}).
		Terminal(false).
		Refresh(true).
		Size(12)

	bulkmenu.Line("Extend Lifetime (+2h)").
		Bash(ex).
		Params([]string{"extend", "all", "TwoHours"}).
		Terminal(false).
		Refresh(true).
		Size(12)

	bulkmenu.Line("Extend Lifetime for one day").
		Bash(ex).
		Params([]string{"extend", "all", "OneDay"}).
		Terminal(false).
		Refresh(true).
		Size(12)

	bulkmenu.Line("Extend Lifetime for three days").
		Bash(ex).
		Params([]string{"extend", "all", "ThreeDays"}).
		Terminal(false).
		Refresh(true).
		Size(12)

	menu.HR()

	menu.Line("New VM")

	newVM := menu.NewSubMenu()

	newVMMap := createNewVMMap(templates)
	var oskeys []string
	for k := range newVMMap {
		oskeys = append(oskeys, k)
	}
	sort.Strings(oskeys)
	for _, oskey := range oskeys {
		newVM.Line(oskey)
		templatemenu := newVM.NewSubMenu()

		osTemplates := newVMMap[oskey]
		sort.Strings(osTemplates)
		for _, template := range osTemplates {
			templatemenu.Line(template).
				Bash(ex).
				Params([]string{"newvm", template}).
				Terminal(false).
				Refresh(true).
				Size(12)
		}
	}
	menu.HR()
	menu.Line("Refresh..").Refresh(true)

	fmt.Print(plugin.Render())
}

func sshUser(vm vm.VM) string {
	var user string
	switch vm.Template.Os {
	case "win":
		user = "Administrator"
	default:
		user = "root"
	}
	return user
}

func createNewVMMap(templates []string) map[string][]string {
	result := map[string][]string{}
	for _, template := range templates {
		os := getTemplateOS(template)
		if _, ok := result[os]; ok {
			result[os] = append(result[os], template)
		} else {
			result[os] = []string{template}
		}
	}
	return result
}

func getTemplateOS(template string) string {
	parts := strings.Split(template, "-")
	return parts[0]
}

func errorMenu(err error) {
	var errMsg string

	switch err.(type) {
	case *url.Error:
		errMsg = "Unable to connect to VMPooler"
	default:
		errMsg = fmt.Sprintf("%s ...", err.Error()[:12])
	}
	plugin := bitbar.New()
	plugin.StatusLine("VMs: ⛔️").Color("red")
	menu := plugin.NewSubMenu()
	menu.Line(errMsg).CopyToClipboard(err.Error())
	// fmt.Print()
	plugin.Render()
	os.Exit(1)
}
