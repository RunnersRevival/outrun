# Outrun for Revival

**DO *NOT* SELF-HOST THIS BRANCH YET. *Currently, changes to this branch are not put into any database migration scripts, and as such us over at Sonic Runners Revival are not responsible if your database is broken when updating to the latest commit.* In addition, not everything works in this branch yet, and some features may be buggy. Please use the master branch for the stable version, or the self-hostable branch for a self-hostable version.**

### Summary

Outrun for Revival is a fork of Outrun, a custom server for Sonic Runners reverse engineered from the [Sonic Runners Revival](https://www.sonicrunners.com/) project back during the Open Beta. It is intended for use on the Sonic Runners Revival server, but can be used for your own private servers as well.

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
  - Leaderboards
  - Gift Box (WIP)
  - Daily Challenge
  - Login Bonuses

Functional:
  - Android and iOS support
  - High score keeping
  - In game notices
  - Deep configuration options
  - Powerful RPC control functions (WIP)
  - Ticker notices
  - Low CPU usage
  - Analytics support (WIP)

### Building

1. [Download and install Go 1.13](https://golang.org/dl/) (Project tested on Go 1.13.3)
2. [Download and install Git](https://git-scm.com/downloads) (for `go get`)
3. Set your [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) environment variable
4. Open a terminal/command prompt
5. Use `cd` ([Windows,](https://www.digitalcitizen.life/command-prompt-how-use-basic-commands) [Linux/macOS](https://www.macworld.com/article/2042378/master-the-command-line-navigating-files-and-folders.html)) to navigate to a directory of choice
6. Run `go get github.com/Mtbcooler/outrun` and wait until the command line returns
7. Run `go build github.com/Mtbcooler/outrun` and wait until the build is complete
8. Run the produced executable (`outrun.exe` on Windows, `outrun` on Linux/macOS)

Binary releases can be found [in the releases tab.](https://github.com/fluofoxxo/outrun/releases)

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

#### Final steps
Coming soon! This section will talk about how to set up your instance with a MySQL server.

### Misc.

Any pull requests deemed code improvements are strongly encouraged. Refactors may be merged into a different branch.

### Credits

Much thanks to:
  - **YPwn**, whose closest point of online social contact I do not know, for creating and running the Sonic Runners Revival server upon which this project bases much of its code upon.
  - **[@Sazpaimon](https://github.com/Sazpaimon)** for finding the encryption key I so desparately looked for but could not on my own.
  - **nacabaro** (nacabaro#2138 on Discord) for traffic logging and the discovery of **[DaGuAr](https://www.youtube.com/user/Gorila5)**'s asset archive.

#### Additional assistance
  - Story Mode items
    - lukaafx (Discord @Kalu04#3243)
    - [TemmieFlakes](https://twitter.com/pictochat3)
    - SuperSonic893YT
