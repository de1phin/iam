package cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "run shell on host via bastion",

	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		validateSshKeysPresent()

		host := args[0]

		user, err := user.Current()
		if err != nil {
			log.Fatal("failed to get currect user: ", err)
		}

		shell, err := NewShell(host, user.Username)
		if err != nil {
			log.Fatal(err)
		}
		defer shell.Close()

		err = shell.Run()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func validateSshKeysPresent() {
	if _, err := os.Stat(IamBastionSshKeyFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal(IamBastionSshKeyFile, "is not present. Please use `iamssh load-key` (requires configured YC CLI)")
	}
	if _, err := os.Stat(IamBastionSshPubKeyFile); errors.Is(err, os.ErrNotExist) {
		log.Fatal(IamBastionSshKeyFile, "is not present. Please use `iamssh load-key` (requires configured YC CLI)")
	}
}

type Shell struct {
	localClient   *ssh.Client
	bastionClient *ssh.Client
	username      string
	host          string
	bastionKey    string
	userKey       string
	userPubKey    string
}

func NewShell(host, username string) (*Shell, error) {
	shell := &Shell{
		host:     host,
		username: username,
	}

	var err error
	shell.localClient, err = createLocalClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create local->bastion ssh client: %w", err)
	}

	err = shell.initUser(username)
	if err != nil {
		return nil, fmt.Errorf("failed to init user: %w", err)
	}

	shell.bastionClient, err = shell.connectToHost()
	if err != nil {
		return nil, fmt.Errorf("failed to create bastion->host ssh client: %w", err)
	}

	return shell, nil
}

func (sh *Shell) connectToHost() (*ssh.Client, error) {
	client, err := sh.createBastionClient("ubuntu", sh.bastionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create root bastion client: %w", err)
	}

	err = sh.initTargetHostUser(client)
	if err != nil {
		return nil, fmt.Errorf("failed to init user %s on target host: %w", sh.username, err)
	}
	client.Close()

	return sh.createBastionClient(sh.username, sh.userKey)
}

func (sh *Shell) initTargetHostUser(client *ssh.Client) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to establish local->bastion ssh session")
	}
	defer session.Close()

	w, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh session stdin: %w", err)
	}
	defer w.Close()
	r, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh session stdout: %w", err)
	}

	err = session.Start("/bin/bash")
	if err != nil {
		return fmt.Errorf("failed to connect ssh: %w", err)
	}

	err = write(w, "cat /etc/passwd; echo '^'\n")
	if err != nil {
		return fmt.Errorf("failed to check if user %s exists: %w", sh.username, err)
	}

	output, err := read(r)
	if err != nil {
		return fmt.Errorf("failed to check if user %s exists: %w", sh.username, err)
	}

	if !strings.Contains(output, sh.username+":") {
		err = write(w, fmt.Sprintf("sudo useradd -s /bin/bash -d /home/%s "+
			"-m -G sudo %s 2>&1; echo '^'\n", sh.username, sh.username))
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", sh.username, err)
		}
		read(r)
	}

	err = write(w, fmt.Sprintf("sudo mkdir /home/%s/.ssh/; echo '%s' | sudo tee 2>&1 /home/%s/.ssh/authorized_keys;\n"+
		"echo '^'\n", sh.username, strings.Trim(sh.userPubKey, "\n"), sh.username))
	read(r)
	if err != nil {
		return fmt.Errorf("failed to write authorized ssh key: %w", err)
	}
	return nil
}

func (sh *Shell) createBastionClient(user, key string) (*ssh.Client, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}
	conn, err := sh.localClient.Dial("tcp", sh.host+":22")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", sh.host, err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(conn, sh.host+":22", config)
	if err != nil {
		return nil, fmt.Errorf("failed to establish ssh conn to %s: %w", sh.host, err)
	}

	return ssh.NewClient(ncc, chans, reqs), nil
}

func createLocalClient() (*ssh.Client, error) {
	keyRaw, err := os.ReadFile(IamBastionSshKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", IamBastionSshKeyFile, err)
	}

	signer, err := ssh.ParsePrivateKey(keyRaw)
	if sshErr, ok := err.(*ssh.PassphraseMissingError); ok && sshErr != nil {
		fmt.Println("Private Key is passphrase protected.")
		fmt.Print("passphrase: ")
		passphrase := ""
		fmt.Scanln(&passphrase)
		signer, err = ssh.ParsePrivateKeyWithPassphrase(keyRaw, []byte(passphrase))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", IamBastionSshKeyFile, err)
	}

	config := &ssh.ClientConfig{
		User:            BastionUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return ssh.Dial("tcp", BastionHost+":22", config)
}

func (sh *Shell) initUser(username string) error {
	session, err := sh.localClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to establish local->bastion ssh session")
	}
	defer session.Close()

	w, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh session stdin: %w", err)
	}
	defer w.Close()
	r, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get ssh session stdout: %w", err)
	}

	err = session.Start("/bin/bash")
	if err != nil {
		return fmt.Errorf("failed to connect ssh: %w", err)
	}

	err = write(w, "ls | grep keys; echo '^'")
	if err != nil {
		return fmt.Errorf("failed to verify key directory exists: %w", err)
	}

	output, err := read(r)
	if err != nil {
		return fmt.Errorf("failed to read stdout: %w", err)
	}

	if !strings.Contains(output, "keys") {
		err = write(w, "mkdir keys\n")
		if err != nil {
			return fmt.Errorf("failed to create key directory: %w", err)
		}
	}

	err = write(w, "ls ./keys; echo '^'")
	if err != nil {
		return fmt.Errorf("failed to verify key directory exists: %w", err)
	}

	output, err = read(r)
	if err != nil {
		return fmt.Errorf("failed to read stdout: %w", err)
	}

	if !strings.Contains(output, username) {
		err = write(w, fmt.Sprintf("ssh-keygen -t rsa -N '' -f ./keys/%s\n", username))
		if err != nil {
			return fmt.Errorf("failed to create key directory: %w", err)
		}
	}

	err = write(w, fmt.Sprintf("cat ./keys/%s.pub; echo '^'\n", username))
	if err != nil {
		return fmt.Errorf("failed to read user pubkey: %w", err)
	}
	sh.userPubKey, err = read(r)
	if err != nil {
		return fmt.Errorf("failed to read user pubkey: %w", err)
	}

	err = write(w, fmt.Sprintf("cat ./keys/%s; echo '^'\n", username))
	if err != nil {
		return fmt.Errorf("failed to read user key: %w", err)
	}
	sh.userKey, err = read(r)
	if err != nil {
		return fmt.Errorf("failed to read user key: %w", err)
	}

	err = write(w, "sudo cat /etc/ssh/ssh_host_rsa_key; echo '^'\n")
	if err != nil {
		return fmt.Errorf("failed to read bastion key: %w", err)
	}
	sh.bastionKey, err = read(r)
	if err != nil {
		return fmt.Errorf("failed to read bastion key: %w", err)
	}

	return nil
}

func write(w io.WriteCloser, command string) error {
	_, err := w.Write([]byte(command + "\n"))
	return err
}

func read(r io.Reader) (string, error) {
	var buf [64 * 1024]byte
	var t int
	for {
		n, err := r.Read(buf[t:])
		if err != nil {
			return "", err
		}
		t += n
		for i := t - n; i <= t; i++ {
			if buf[i] == '^' {
				str := string(buf[:i])
				return str, nil
			}
		}
	}
}

func (sh *Shell) Run() error {
	session, err := sh.bastionClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to establish new ssh session: %w", err)
	}
	defer session.Close()

	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	err = session.RequestPty("xterm", 80, 40, ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	})
	if err != nil {
		return fmt.Errorf("failed to initialize ssh session: %w", err)
	}

	err = session.Start("/bin/bash")
	if err != nil {
		return fmt.Errorf("failed to start ssh session: %w", err)
	}

	return session.Wait()
}

func (sh *Shell) Close() {
	sh.bastionClient.Close()
	sh.localClient.Close()
}
