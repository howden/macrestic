package main

import (
	"fmt"
	"github.com/keybase/go-keychain"
	"os"
	"os/exec"
	"syscall"
)

var (
	ResticCommand            = "restic"
	SecretAccountId          = "backup-restic-b2-account-id"
	SecretAccountKey         = "backup-restic-b2-account-key"
	SecretResticRepo         = "backup-restic-repo"
	SecretResticRepoPassword = "backup-restic-repo-password"
)

func main() {
	// To add the secrets to the keychain:
	// ==> `security add-generic-password -s backup-restic-b2-account-id -a restic_backup -w`

	accountId := getSecret(SecretAccountId, "restic_backup")
	accountKey := getSecret(SecretAccountKey, "restic_backup")
	repo := getSecret(SecretResticRepo, "restic_backup")
	repoPassword := getSecret(SecretResticRepoPassword, "restic_backup")

	setEnvVar("B2_ACCOUNT_ID", accountId)
	setEnvVar("B2_ACCOUNT_KEY", accountKey)
	setEnvVar("RESTIC_REPOSITORY", repo)
	setEnvVar("RESTIC_PASSWORD", repoPassword)

	bin, err := exec.LookPath(ResticCommand)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to find restic - is it installed?: %s", err)
		os.Exit(1)
	}

	env := os.Environ()
	args := append([]string{ResticCommand}, os.Args[1:]...)

	err = syscall.Exec(bin, args, env)
	if err != nil {
		panic(err)
	}
}

func setEnvVar(key string, value string) {
	if err := os.Setenv(key, value); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to set environment variable %s: %s", key, err)
		os.Exit(1)
	}
}

func getSecret(service string, account string) string {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(service)
	query.SetAccount(account)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "unable to get secret %s: %s", service, err)
		os.Exit(1)
		return ""
	} else if len(results) != 1 {
		_, _ = fmt.Fprintf(os.Stderr, "secret value '%s' not found", service)
		os.Exit(1)
		return ""
	} else {
		return string(results[0].Data)
	}
}
