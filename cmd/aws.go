package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/fatih/color"
	"github.com/go-ini/ini"
	homedir "github.com/mitchellh/go-homedir"
)

func login(region string, token string, serialNumber string, profile string, assumedRoleName string, account Account) {
	svc := sts.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile),
	}),
	)

	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(3600),
		RoleArn:         aws.String("arn:aws:iam::" + account.ID + ":role/" + assumedRoleName),
		RoleSessionName: aws.String("gasy" + user.Username),
		SerialNumber:    aws.String(serialNumber),
		TokenCode:       aws.String(token),
	}

	if account.SID != "" {
		input.SetExternalId(account.SID)
	}

	result, err := svc.AssumeRole(input)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	setProfile(result, account)
	url, err := getAWSConsoleURL(result)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	boldGreen := color.New(color.FgGreen, color.Bold)
	boldGreen.Println("Credentials written to profile!")
	fmt.Println()
	boldGreen.Println("export AWS_PROFILE=" + account.Name)
	fmt.Println()
	boldGreen.Println("URL: " + url)
}

func setProfile(credentials *sts.AssumeRoleOutput, account Account) {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg, err := ini.Load(home + "/.aws/credentials")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	// write the temporary credentials to a profile

	// TODO: I think we can do this with reflection too
	cfg.Section(account.Name).Key("aws_access_key_id").SetValue(*credentials.Credentials.AccessKeyId)
	cfg.Section(account.Name).Key("aws_secret_access_key").SetValue(*credentials.Credentials.SecretAccessKey)
	cfg.Section(account.Name).Key("aws_session_token").SetValue(*credentials.Credentials.SessionToken)
	cfg.Section(account.Name).Key("expiration").SetValue(credentials.Credentials.Expiration.String())

	err = cfg.SaveTo(home + "/.aws/credentials")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

func getAWSConsoleURL(credentials *sts.AssumeRoleOutput) (string, error) {
	session := map[string]string{
		"sessionId":    *credentials.Credentials.AccessKeyId,
		"sessionKey":   *credentials.Credentials.SecretAccessKey,
		"sessionToken": *credentials.Credentials.SessionToken,
	}
	sessionString, err := json.Marshal(session)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	federationValues := url.Values{}
	federationValues.Add("Action", "getSigninToken")
	federationValues.Add("Session", string(sessionString))
	federationURL := "https://signin.aws.amazon.com/federation?" +
		federationValues.Encode()

	federationResponse, err := http.Get(federationURL)
	if err != nil {
		return "", fmt.Errorf("fetching federated signin URL: %s", err)
	}
	tokenDocument := struct{ SigninToken string }{}
	err = json.NewDecoder(federationResponse.Body).Decode(&tokenDocument)
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Add("Action", "login")
	values.Add("Destination",
		"https://console.aws.amazon.com/")
	values.Add("SigninToken", tokenDocument.SigninToken)

	return "https://signin.aws.amazon.com/federation?" + values.Encode(), nil
}
