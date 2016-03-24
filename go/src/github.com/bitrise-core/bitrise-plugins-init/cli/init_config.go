package cli

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-core/bitrise-plugins-init/templates"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/goinp/goinp"
	"github.com/codegangsta/cli"
)

// ConfigModel ...
type ConfigModel struct {
	FormatVersion string
	AppTitle      string
	DevBranch     string
}

const (
	bitriseConfigFileName  = "bitrise.yml"
	bitriseSecretsFileName = ".bitrise.secrets.yml"
)

func generateBitriseYMLContent(config ConfigModel) (string, error) {
	bitriseConfigTemplate := template.New("bitrise_config")
	bitriseConfigTemplate, err := bitriseConfigTemplate.Parse(templates.BitriseConfigTemplate)
	if err != nil {
		log.Fatalf("failed to parse bitrise config template, error: %#v", err)
	}

	var bitriseConfigBytes bytes.Buffer
	err = bitriseConfigTemplate.Execute(&bitriseConfigBytes, config)
	if err != nil {
		return "", err
	}

	return bitriseConfigBytes.String(), nil
}

func saveSecretsToFile(pth, secretsStr string) (bool, error) {
	if exists, err := pathutil.IsPathExists(pth); err != nil {
		return false, err
	} else if exists {
		ask := fmt.Sprintf("A secrets file already exists at %s - do you want to overwrite it?", pth)
		if val, err := goinp.AskForBool(ask); err != nil {
			return false, err
		} else if !val {
			log.Infoln("Init canceled, existing file (" + pth + ") won't be overwritten.")
			return false, nil
		}
	}

	if err := fileutil.WriteStringToFile(pth, secretsStr); err != nil {
		return false, err
	}
	return true, nil
}

func addToGitignore(ignorePattern string) error {
	return fileutil.AppendStringToFile(".gitignore", "\n"+ignorePattern+"\n")
}

func initConfig(c *cli.Context) {
	bitriseSecretsFileRelPath := "./" + bitriseSecretsFileName
	bitriseConfigFileRelPath := "./" + bitriseConfigFileName

	if exists, err := pathutil.IsPathExists(bitriseConfigFileRelPath); err != nil {
		log.Fatalln("Error:", err)
	} else if exists {
		ask := fmt.Sprintf("A config file already exists at %s - do you want to overwrite it?", bitriseConfigFileRelPath)
		if val, err := goinp.AskForBool(ask); err != nil {
			log.Fatalln("Error:", err)
		} else if !val {
			log.Infoln("Init canceled, existing file won't be overwritten.")
			os.Exit(0)
		}
	}

	config := ConfigModel{
		FormatVersion: "1.2.0",
	}

	if val, err := goinp.AskForString("What's the BITRISE_APP_TITLE?"); err != nil {
		log.Fatalln(err)
	} else {
		config.AppTitle = val
	}
	if val, err := goinp.AskForString("What's your development branch's name?"); err != nil {
		log.Fatalln(err)
	} else {
		config.DevBranch = val
	}

	bitriseConfContent, err := generateBitriseYMLContent(config)
	if err != nil {
		log.Fatalf("Invalid Bitrise YML: %s", err)
	}

	if err := fileutil.WriteStringToFile(bitriseConfigFileRelPath, bitriseConfContent); err != nil {
		log.Fatalln("Failed to init the bitrise config file:", err)
	} else {
		fmt.Println()
		fmt.Println("# NOTES about the " + bitriseConfigFileName + " config file:")
		fmt.Println()
		fmt.Println("We initialized a " + bitriseConfigFileName + " config file for you.")
		fmt.Println("If you're in this folder you can use this config file")
		fmt.Println(" with bitrise automatically, you don't have to")
		fmt.Println(" specify it's path.")
		fmt.Println()
	}

	bitriseSecretsContent := templates.BitriseSecretsTemplate
	if initialized, err := saveSecretsToFile(bitriseSecretsFileRelPath, bitriseSecretsContent); err != nil {
		log.Fatalln("Failed to init the secrets file:", err)
	} else if initialized {
		fmt.Println()
		fmt.Println("# NOTES about the " + bitriseSecretsFileName + " secrets file:")
		fmt.Println()
		fmt.Println("We also created a " + bitriseSecretsFileName + " file")
		fmt.Println(" in this directory, to keep your passwords, absolute path configurations")
		fmt.Println(" and other secrets separate from your")
		fmt.Println(" main configuration file.")
		fmt.Println("This way you can safely commit and share your configuration file")
		fmt.Println(" and ignore this secrets file, so nobody else will")
		fmt.Println(" know about your secrets.")
		fmt.Println(colorstring.Yellow("You should NEVER commit this secrets file into your repository!!"))
		fmt.Println()
	}

	// add the general .bitrise* item
	//  which will include both secret files like .bitrise.secrets.yml
	//  and the .bitrise work temp dir
	if err := addToGitignore(".bitrise*"); err != nil {
		log.Fatalln("Failed to add .gitignore pattern. Error: ", err)
	}
	fmt.Println(colorstring.Green("For your convenience we added the pattern '.bitrise*' to your .gitignore file"))
	fmt.Println(" to make it sure that no secrets or temporary work directories will be")
	fmt.Println(" committed into your repository.")

	fmt.Println()
	fmt.Println("Hurray, you're good to go!")
	fmt.Println("You can simply run:")
	fmt.Println("-> bitrise run test")
	fmt.Println("to test the sample configuration (which contains")
	fmt.Println("an example workflow called 'test').")
	fmt.Println()
	fmt.Println("Once you tested this sample setup you can")
	fmt.Println(" open the " + bitriseConfigFileName + " config file,")
	fmt.Println(" modify it and then run a workflow with:")
	fmt.Println("-> bitrise run YOUR-WORKFLOW-NAME")
	fmt.Println(" or trigger a build with a pattern:")
	fmt.Println("-> bitrise trigger YOUR/PATTERN")
}
