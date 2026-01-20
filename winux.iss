; ============================================================================
; WINUX Installer Script
; Inno Setup 6.x
; ============================================================================
; 
; Build command:
;   "C:\Program Files (x86)\Inno Setup 6\ISCC.exe" winux.iss
;
; After build, sign with:
;   signtool sign /f cert.pfx /p PASSWORD /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 installer\winux-VERSION-setup.exe
;
; ============================================================================

; ----------------------------------------------------------------------------
; Version Configuration (UPDATE THIS FOR NEW RELEASES)
; ----------------------------------------------------------------------------
#define MyAppName "WINUX"
#ifndef MyAppVersion
  #define MyAppVersion "0.3.10"
#endif

#define MyAppPublisher "CRTYPUBG"
#define MyAppURL "https://github.com/CRTYPUBG/winux"
#define MyAppExeName "winux.exe"

#define MyAppVerName MyAppName + " " + MyAppVersion

; ----------------------------------------------------------------------------
; Signing Configuration (Local only - comment out for GitHub Actions)
; ----------------------------------------------------------------------------
; #define SignToolPath "C:\Program Files (x86)\Windows Kits\10\bin\10.0.22621.0\x64\signtool.exe"
; #define CertFile "C:\Users\LenovoPC\cert.pfx"
; #define CertPass "ueo586_crty555"
; #define TimestampURL "http://timestamp.digicert.com"
; #define SignCommand '"' + SignToolPath + '" sign /f "' + CertFile + '" /p "' + CertPass + '" /fd SHA256 /tr ' + TimestampURL + ' /td SHA256 /v $f'


[Setup]
; Application identity
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppVerName}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}/releases

; Installation paths
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes

; License
LicenseFile=LICENSE

; Output
OutputDir=installer
OutputBaseFilename=winux-{#MyAppVersion}-setup

; Icons
SetupIconFile=assets\winux.ico
UninstallDisplayIcon={app}\winux.ico
UninstallDisplayName={#MyAppName}

; Compression (maximum)
Compression=lzma2/ultra64
SolidCompression=yes

; Modern Wizard Style (Inno Setup 6+)
WizardStyle=modern
WizardResizable=no
; WizardSizePercent=100
WizardImageFile=none
WizardSmallImageFile=none

; Permissions
PrivilegesRequired=admin
ArchitecturesInstallIn64BitMode=x64
ArchitecturesAllowed=x64

; Environment
ChangesEnvironment=yes

; Note: Sign the output installer manually after build:
; signtool sign /f cert.pfx /p PASSWORD /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 installer\winux-VERSION-setup.exe

; Misc
DisableWelcomePage=no
DisableProgramGroupPage=yes
ShowLanguageDialog=auto
CloseApplications=yes
RestartApplications=no

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
; Name: "turkish"; MessagesFile: "compiler:Languages\Turkish.isl"

[Tasks]
Name: "addtopath"; Description: "Add WINUX to system PATH (recommended)"; GroupDescription: "System Integration:"; Flags: checkedonce
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
; Main binaries (pre-signed)
Source: "winux.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "update.exe"; DestDir: "{app}"; Flags: ignoreversion

; Resources
Source: "assets\winux.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "LICENSE"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\winux.ico"; Comment: "Linux-like CLI for Windows"
Name: "{group}\WINUX Update"; Filename: "{app}\update.exe"; IconFilename: "{app}\winux.ico"; Parameters: "--check"; Comment: "Check for updates"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\winux.ico"; Tasks: desktopicon

[Registry]
; Add to system PATH
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addtopath; Check: NeedsAddPath('{app}')

; App registration
Root: HKLM; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\winux.exe"; ValueType: string; ValueName: ""; ValueData: "{app}\winux.exe"; Flags: uninsdeletekey
Root: HKLM; Subkey: "SOFTWARE\Microsoft\Windows\CurrentVersion\App Paths\winux.exe"; ValueType: string; ValueName: "Path"; ValueData: "{app}"

[Code]
// ============================================================================
// Pascal Script Functions
// ============================================================================

// Check if path needs to be added to system PATH
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKEY_LOCAL_MACHINE,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  // Check if path already exists
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

// Remove from PATH on uninstall
procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
var
  OrigPath: string;
  NewPath: string;
  AppPath: string;
begin
  if CurUninstallStep = usPostUninstall then
  begin
    AppPath := ExpandConstant('{app}');
    if RegQueryStringValue(HKEY_LOCAL_MACHINE,
      'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
      'Path', OrigPath) then
    begin
      NewPath := OrigPath;
      // Remove all variations of the path
      StringChangeEx(NewPath, ';' + AppPath, '', True);
      StringChangeEx(NewPath, AppPath + ';', '', True);
      StringChangeEx(NewPath, AppPath, '', True);
      // Write back if changed
      if NewPath <> OrigPath then
      begin
        RegWriteExpandStringValue(HKEY_LOCAL_MACHINE,
          'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
          'Path', NewPath);
      end;
    end;
  end;
end;

// Custom initialization
function InitializeSetup(): Boolean;
begin
  Result := True;
  // Additional checks can be added here
end;

// Post-installation actions
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall then
  begin
    // Notify shell about environment changes
    // This helps update PATH without requiring reboot
  end;
end;

[UninstallDelete]
; Clean up everything in install directory (but not root/parent folders)
Type: filesandordirs; Name: "{app}\*"
Type: dirifempty; Name: "{app}"

[Messages]
; Custom messages
WelcomeLabel1=Welcome to WINUX Setup
WelcomeLabel2=This will install [name/ver] on your computer.%n%nWINUX provides native Linux-like command-line utilities for Windows.%n%nNo WSL. No emulation. Just native performance.

[Run]
; Post-installation launch options
Filename: "{app}\{#MyAppExeName}"; Parameters: "--version"; Description: "Show installed version"; Flags: postinstall skipifsilent nowait runhidden
Filename: "cmd"; Parameters: "/k cd /d ""{app}"" && echo WINUX {#MyAppVersion} installed successfully! && winux --help"; Description: "Open WINUX in Command Prompt"; Flags: postinstall skipifsilent unchecked nowait
