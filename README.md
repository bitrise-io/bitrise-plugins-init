# Init Plugin for [Bitrise CLI](https://github.com/bitrise-io/bitrise)

Initialize bitrise __config (bitrise.yml)__ based on your project.

* For __iOS__ projects detects CocoaPods and scans Xcode project files for valid Xcode command line configurations.

* For __Android__ projects checks for gradle files and lists all the gradle tasks, also checks for gradlew file.

* For __Xamarin__ projects inspects the solution files and lists the configuration options, also checks for NuGet and Xamarin Components packages.

* For __Fastlane__ detects Fastfile and lists the available lanes.

## How to use this Plugin

Can be run directly with the Bitrise CLI, requires version 1.3.0 or newer.

First install the plugin:

```
bitrise plugin install --source https://github.com/bitrise-core/bitrise-plugins-init.git
```

After that, you can use it:

```
bitrise :init
```
