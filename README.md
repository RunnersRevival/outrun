# Outrun for Revival

### Summary

Outrun for Revival is a fork of Outrun, a custom server for Sonic Runners reverse engineered from the [Sonic Runners Revival](https://sonic.runner.es/) project back during the Open Beta. It is intended for use on the Sonic Runners Revival server, but can be used for your own private servers as well.

Unlike the master and mysql_testing branches, this version can still be used with 2.0.3. Code changes we do to Outrun for Revival may be backported to this branch occasionally, so you won't miss out on the latest bug fixes. However, since this version of Outrun relies on 2.0.3, no Revival-exclusive content can be used.

### Current functionality

Notable:
  - Timed Mode
  - Story Mode
  - Ring/Red Star Ring keeping
  - Functional shop
  - Character/Chao equipping
  - Character leveling and progression
  - Item/Chao roulette functionality
  - Events
  - Campaigns
  - Basic ranking
  - Login Bonuses
  - Daily Challenge

Functional:
  - Android and iOS support
  - High score keeping
  - In game notices
  - Deep configuration options
  - Powerful RPC control functions
  - Ticker notices
  - Low CPU usage
  - Analytics support
  - Revive Token keeping

### Building

1. [Download and install Go](https://golang.org/dl/) (latest version should work, but the availability of the legacy `go get` syntax )
2. [Download and install Git](https://git-scm.com/downloads) (for `go get`)
3. Set your [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) environment variable (on Windows, the default is `C:\Users\UsernameHere\go`, while on Linux the default is `/home/usernamehere/go` - only change this if you want to put your work in another directory)
4. Open a terminal/command prompt
5. Use `cd` ([Windows,](https://www.digitalcitizen.life/command-prompt-how-use-basic-commands) [Linux/macOS](https://www.macworld.com/article/2042378/master-the-command-line-navigating-files-and-folders.html)) to navigate to a directory of choice
6. Run `go get github.com/RunnersRevival/outrun@self-hostable` and wait until the command line returns
7. Run `go build github.com/RunnersRevival/outrun` and wait until the build is complete
8. Run the produced executable (`outrun.exe` on Windows, `outrun` on Linux/macOS)

#### Modifying an APK to connect to your instance (from Windows)

1. Install [dnSpy](https://github.com/0xd4d/dnSpy/releases) (dnSpy-netcore-win64.zip)
2. Install [7-Zip](https://www.7-zip.org/download.html)
3. Install [ZipSigner](https://www.apkmirror.com/apk/ken-ellinwood/zipsigner/zipsigner-3-4-release/zipsigner-3-4-android-apk-download/) on an Android device or emulator
4. Open a Sonic Runners v2.0.3 APK file with 7-Zip
5. Navigate to assets/bin/Data/Managed and extract all the DLL files to their own folder
6. Open Assembly-CSharp.dll in dnSpy
7. Open the class `NetBaseUtil`, and find the variable `mActionServerUrlTable `
8. Edit every string in the `mActionServerUrlTable` array to `http://<IP>:<PORT>/` where `<IP>` is replaced by the IP for your instance and `<PORT>` is replaced by the port for your instance (Default: 9001)
9. Repeat step 7 for `mSecureActionServerUrlTable`
10. If you have an assets server, use its IP and port to replace the values in `mAssetURLTable` and `mInformationURLTable` to `http://<IP>:<PORT>/assets/` and `http://<IP>:<PORT>/information/` respectively
11. Click File -> Save Module... and save the DLL file
12. Drag the newly saved Assembly-CSharp.dll back into assets/bin/Data/Managed in 7-Zip, confirming to overwrite if asked
13. Transfer the APK to an Android device and use ZipSigner to sign it
14. Install the APK

#### Modifying an IPA to connect to your instance (from Windows)

*Disclaimer: Make sure your server URLs (including the port, information, and asset extensions) are below 26 characters so it can fit within the hex*
1. Install [7-Zip](https://www.7-zip.org/download.html)
2. Install a hex editor of your choice
3. Open a Sonic Runners v2.0.3 IPA file with 7-Zip
4. Navigate to Payload/sonicrunners.app/Data/Managed/Metadata and extract global-metdata.dat
5. Open global-metadata.dat in your hex editor
6. Jump to offset `2CEE0` and find the table of server URLs
7. Edit every string in the table to match your server URLs (If your URL is below the character limit, you can use slashes to pad out the empty space until you reach the next listing)
8. Save global-metadata.dat
9. Insert the modified global-metadata.dat into Payload/sonicrunners.app/Data/Managed/Metadata in 7-Zip, confirming to overwrite if asked
10. Transfer the IPA to an iOS device
11. Install the IPA using AltStore or AppSync Unified

### Misc.

Any pull requests deemed code improvements are strongly encouraged. Refactors may be merged into a different branch.

### Credits

Much thanks to:
  - **YPwn**, whose closest point of online social contact I do not know, for creating and running the Sonic Runners Revival server upon which this project based much of its code upon.
  - **[@Sazpaimon](https://github.com/Sazpaimon)** for finding the encryption key I so desparately looked for but could not on my own.
  - **nacabaro** (nacabaro#2138 on Discord) for traffic logging and the discovery of **[DaGuAr](https://www.youtube.com/user/Gorila5)**'s asset archive.

#### Additional assistance
  - Story Mode items
    - lukaafx (Discord @Kalu04#3243)
    - [TemmieFlakes](https://twitter.com/pictochat3)
    - SuperSonic893YT
