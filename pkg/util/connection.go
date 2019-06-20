package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/openshift-online/uhc-cli/pkg/config"
	"github.com/openshift-online/uhc-cli/pkg/dump"
	"github.com/openshift-online/uhc-sdk-go/pkg/client"
)

func NewConnection(debug bool) (*client.Connection, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("Can't load config file: %v\n", err)
	}
	if cfg == nil {
		fmt.Fprintf(os.Stderr, "Not logged in, run the 'login' command\n")
		os.Exit(1)
	}

	// Check that the configuration has credentials or tokens that don't have expired:
	armed, err := config.Armed(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't check if tokens have expired: %v\n", err)
		os.Exit(1)
	}
	if !armed {
		fmt.Fprintf(os.Stderr, "Tokens have expired, run the 'login' command\n")
		os.Exit(1)
	}

	// Create the connection:
	logger, err := NewLogger(debug)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't create logger: %v\n", err)
		os.Exit(1)
	}
	connection, err := client.NewConnectionBuilder().
		Logger(logger).
		TokenURL(cfg.TokenURL).
		Client(cfg.ClientID, cfg.ClientSecret).
		Scopes(cfg.Scopes...).
		URL(cfg.URL).
		User(cfg.User, cfg.Password).
		Tokens(cfg.AccessToken, cfg.RefreshToken).
		Insecure(cfg.Insecure).
		Build()

	// Save the configuration:
	cfg.AccessToken, cfg.RefreshToken, err = connection.Tokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't get tokens: %v\n", err)
		os.Exit(1)
	}
	err = config.Save(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't save config file: %v\n", err)
		os.Exit(1)
	}

	return connection, nil
}

func AddBody(request *client.Request, body string) {
	var bytes []byte
	var err error
	if body != "" {
		bytes, err = ioutil.ReadFile(body)
	} else {
		bytes, err = ioutil.ReadAll(os.Stdin)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read body: %v\n", err)
		os.Exit(1)
	}
	request.Bytes(bytes)
}

func AddParamsAndHeaders(request *client.Request, parameters, headers []string) {
	for _, parameter := range parameters {
		var name string
		var value string
		position := strings.Index(parameter, "=")
		if position != -1 {
			name = parameter[:position]
			value = parameter[position+1:]
		} else {
			name = parameter
			value = ""
		}
		request.Parameter(name, value)
	}
	for _, header := range headers {
		var name string
		var value string
		position := strings.Index(header, "=")
		if position != -1 {
			name = header[:position]
			value = header[position+1:]
		} else {
			name = header
			value = ""
		}
		request.Header(name, value)
	}
}

func DoHTTP(request *client.Request) {
	// Send the request:
	response, err := request.Send()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't send request: %v\n", err)
		os.Exit(1)
	}
	status := response.Status()
	body := response.Bytes()
	if status < 400 {
		err = dump.Pretty(os.Stdout, body)
	} else {
		err = dump.Pretty(os.Stderr, body)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't print body: %v\n", err)
		os.Exit(1)
	}

	// Bye:
	if status < 400 {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}
