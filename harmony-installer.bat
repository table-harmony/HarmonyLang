@echo off
setlocal enabledelayedexpansion

echo Installing HarmonyLang...

:: Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go first.
    exit /b 1
)

:: Create harmony directories
set "HARMONY_HOME=%USERPROFILE%\.harmony"
set "HARMONY_BIN=%HARMONY_HOME%\bin"
mkdir "%HARMONY_BIN%" 2>nul

:: Clone/update the repository
set "INSTALL_PATH=%HARMONY_HOME%\HarmonyLang"
if exist "%INSTALL_PATH%" (
    echo Updating HarmonyLang...
    cd "%INSTALL_PATH%"
    git pull
) else (
    echo Cloning HarmonyLang...
    git clone https://github.com/table-harmony/HarmonyLang.git "%INSTALL_PATH%"
)

:: Build the project
echo Building HarmonyLang...
cd "%INSTALL_PATH%"
go mod download
go build -v -o "%HARMONY_BIN%\HarmonyLang.exe" "%INSTALL_PATH%\src\main.go"

if %ERRORLEVEL% NEQ 0 (
    echo Error: Build failed
    exit /b 1
)

:: Create harmony.bat wrapper script
echo Creating wrapper script...
(
echo @echo off
echo IF "%%~1"=="run" ^(
echo     IF NOT "%%~2"=="" ^(
echo         IF EXIST "%%~2" ^(
echo             "%HARMONY_BIN%\HarmonyLang.exe" "%%~2"
echo         ^) ELSE ^(
echo             echo Error: File %%~2 not found
echo             exit /b 1
echo         ^)
echo     ^) ELSE ^(
echo         echo Error: No file specified
echo         exit /b 1
echo     ^)
echo ^) ELSE IF "%%~1"=="repl" ^(
echo     "%HARMONY_BIN%\HarmonyLang.exe"
echo ^) ELSE ^(
echo     echo Usage:
echo     echo   harmony run ^<file.harmony^>  - Run a Harmony source file
echo     echo   harmony repl               - Start Harmony REPL
echo ^)
) > "%HARMONY_BIN%\harmony.bat"

:: Update PATH - both user and current session
echo Updating PATH...

:: Get current user PATH
for /f "tokens=2*" %%a in ('reg query HKCU\Environment /v PATH') do set "USER_PATH=%%b"

:: Check if our bin directory is already in PATH
echo !USER_PATH! | find /i "%HARMONY_BIN%" > nul
if errorlevel 1 (
    :: Add to user PATH if not present
    if defined USER_PATH (
        setx PATH "!USER_PATH!;%HARMONY_BIN%"
    ) else (
        setx PATH "%HARMONY_BIN%"
    )
    
    :: Also add to current session PATH
    set "PATH=%PATH%;%HARMONY_BIN%"
)

echo.
echo Installation complete!
echo.
echo To verify installation, please:
echo 1. Close ALL Command Prompt windows
echo 2. Open a new Command Prompt window
echo 3. Type: harmony repl
echo.
echo If it still doesn't work, you can try running these commands manually:
echo setx PATH "%%PATH%%;%HARMONY_BIN%"
echo.
echo Or add this path manually to your System Environment Variables:
echo %HARMONY_BIN%

pause