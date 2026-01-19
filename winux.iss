; WINUX Installer Script
; Inno Setup 6.x

#define MyAppName "WINUX"
#define MyAppVersion "0.1.0"
#define MyAppPublisher "CRTYPUBG"
#define MyAppURL "https://github.com/CRTYPUBG/winux"
#define MyAppExeName "winux.exe"

[Setup]
AppId={{A1B2C3D4-E5F6-7890-ABCD-EF1234567890}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}/releases
DefaultDirName={autopf}\{#MyAppName}
DefaultGroupName={#MyAppName}
AllowNoIcons=yes
LicenseFile=LICENSE
OutputDir=installer
OutputBaseFilename=winux-{#MyAppVersion}-setup
SetupIconFile=assets\winux.ico
Compression=lzma2/ultra64
SolidCompression=yes
WizardStyle=modern
PrivilegesRequired=admin
ArchitecturesInstallIn64BitMode=x64
ChangesEnvironment=yes
UninstallDisplayIcon={app}\{#MyAppExeName}
UninstallDisplayName={#MyAppName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "turkish"; MessagesFile: "compiler:Languages\Turkish.isl"

[Tasks]
Name: "addtopath"; Description: "Add WINUX to system PATH"; GroupDescription: "Additional options:"; Flags: checkedonce
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "winux.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets\winux.ico"; DestDir: "{app}"; Flags: ignoreversion
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "LICENSE"; DestDir: "{app}"; Flags: ignoreversion skipifsourcedoesntexist

[Icons]
Name: "{group}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\winux.ico"
Name: "{group}\{cm:UninstallProgram,{#MyAppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; IconFilename: "{app}\winux.ico"; Tasks: desktopicon

[Registry]
; Add to PATH
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; Tasks: addtopath; Check: NeedsAddPath('{app}')

[Code]
// Check if path needs to be added
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
      StringChangeEx(NewPath, ';' + AppPath, '', True);
      StringChangeEx(NewPath, AppPath + ';', '', True);
      StringChangeEx(NewPath, AppPath, '', True);
      if NewPath <> OrigPath then
      begin
        RegWriteExpandStringValue(HKEY_LOCAL_MACHINE,
          'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
          'Path', NewPath);
      end;
    end;
  end;
end;

[UninstallDelete]
; Clean up everything in install directory (but not root/parent folders)
Type: filesandordirs; Name: "{app}\*"
Type: dirifempty; Name: "{app}"

[Run]
Filename: "{cmd}"; Parameters: "/C echo WINUX installed successfully! && pause"; Description: "Show installation complete message"; Flags: postinstall skipifsilent runhidden
