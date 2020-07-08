package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(infoCmd)
}

var infoCmd = &cobra.Command{
	Use:    "info",
	Run:    runInfo,
	Hidden: true,
}

func runInfo(cmd *cobra.Command, args []string) {

	// cfg, err := config.Read()
	// if err != nil {
	// 	os.Exit(1)
	// }

	// if len(args) < 1 {
	// 	os.Exit(1)
	// }

	// target := args[0]

	// vmclient := vm.NewClient(cfg.Endpoint, cfg.Token)

	// var vms []vm.VM
	// var virtualmachine *vm.VM

	// virtualmachine, err = vmclient.Get(target)
	// vms = []vm.VM{*virtualmachine}

	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)
	// }
	// ssh -t username@host 'top'

	// "ssh://%s@%s", sshUser(vm), vm.Fqdn))
	// ssh -t username@host 'top'

	sshPupVersionCmd := fmt.Sprintf("%s@%s", args[0], args[1])

	puppetVersion, err := exec.Command("ssh", "-t", sshPupVersionCmd, "puppet --version").CombinedOutput()
	// execute(string(puppetVersion), err.Error(), "3333")
	// execute(string(args[0]), args[1], "3333")
	execute(string(puppetVersion), err.Error(), "3333")

}

func execute(puppetVersion string, puppetEnterpriseVersion string, boltVersion string) {
	cmd := fmt.Sprintf("tell application (path to frontmost application as text) to display dialog \"Puppet Version: %s%s\"", puppetVersion, puppetEnterpriseVersion)
	exec.Command("osascript",
		"-e",
		cmd).CombinedOutput()
}

func getPuppetVersion() {

}
