package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestReactNativeExpo(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__reactnative_expo__")
	require.NoError(t, err)

	t.Log("sample-apps-react-native-expo-ios-and-android")
	{
		sampleAppDir := filepath.Join(tmpDir, "sample-apps-react-native-expo")
		sampleAppURL := "https://github.com/bitrise-samples/sample-apps-react-native-expo.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)

		validateConfigExpectation(t, "sample-apps-react-native-expo-ios-and-android", strings.TrimSpace(sampleAppsReactNativeExpoIosAndAndroidResultYML), strings.TrimSpace(result), sampleAppsReactNativeExpoIosAndAndroidVersions...)
	}
}

var sampleAppsReactNativeExpoIosAndAndroidVersions = []interface{}{
	models.FormatVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.InstallMissingAndroidToolsVersion,
	steps.AndroidBuildVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.XcodeArchiveVersion,
	steps.DeployToBitriseIoVersion,

	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.DeployToBitriseIoVersion,
}

var sampleAppsReactNativeExpoIosAndAndroidResultYML = fmt.Sprintf(`options:
  react-native-expo:
    title: The root directory of an Android project
    env_key: PROJECT_LOCATION
    value_map:
      android:
        title: Module
        env_key: MODULE
        value_map:
          app:
            title: Variant for building
            env_key: BUILD_VARIANT
            value_map:
              "":
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              AndroidTest:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              Debug:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              DebugAndroidTest:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              DebugUnitTest:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              Release:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
              ReleaseUnitTest:
                title: Project (or Workspace) path
                env_key: BITRISE_PROJECT_PATH
                value_map:
                  ios/sampleappsreactnativeexpo.xcodeproj:
                    title: Scheme name
                    env_key: BITRISE_SCHEME
                    value_map:
                      sampleappsreactnativeexpo:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
                      sampleappsreactnativeexpo-tvOS:
                        title: ipa export method
                        env_key: BITRISE_EXPORT_METHOD
                        value_map:
                          ad-hoc:
                            config: react-native-expo-android-ios-test-config
                          app-store:
                            config: react-native-expo-android-ios-test-config
                          development:
                            config: react-native-expo-android-ios-test-config
                          enterprise:
                            config: react-native-expo-android-ios-test-config
configs:
  react-native-expo:
    react-native-expo-android-ios-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: react-native-expo
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        deploy:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - npm@%s:
              inputs:
              - command: run eject
          - install-missing-android-tools@%s:
              inputs:
              - gradlew_path: $PROJECT_LOCATION/gradlew
          - android-build@%s:
              inputs:
              - project_location: $PROJECT_LOCATION
              - module: $MODULE
              - variant: $BUILD_VARIANT
          - certificate-and-profile-installer@%s: {}
          - xcode-archive@%s:
              inputs:
              - project_path: $BITRISE_PROJECT_PATH
              - scheme: $BITRISE_SCHEME
              - export_method: $BITRISE_EXPORT_METHOD
              - configuration: Release
          - deploy-to-bitrise-io@%s: {}
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - npm@%s:
              inputs:
              - command: install
          - npm@%s:
              inputs:
              - command: test
          - deploy-to-bitrise-io@%s: {}
warnings:
  react-native-expo: []
`, sampleAppsReactNativeExpoIosAndAndroidVersions...)
