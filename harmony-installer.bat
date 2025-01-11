@echo off
echo Installing HarmonyLang...

:: Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo Error: Go is not installed. Please install Go first.
    exit /b 1
)

:: Set up Go workspace if it doesn't exist
if not defined GOPATH (
    echo GOPATH is not set. Using default...
    set "GOPATH=%USERPROFILE%\go"
)

:: Create harmony directory
mkdir "%USERPROFILE%\.harmony\bin" 2>nul

:: Clone/update the repository
set "INSTALL_PATH=%USERPROFILE%\.harmony\HarmonyLang"
if exist "%INSTALL_PATH%" (
    echo Updating HarmonyLang...
    cd "%INSTALL_PATH%"
    git pull
) else (
    echo Cloning HarmonyLang...
    git clone https://github.com/table-harmony/HarmonyLang.git "%INSTALL_PATH%"
    cd "%INSTALL_PATH%"
)

:: Build the project
echo Building HarmonyLang...
cd "%INSTALL_PATH%"
go build -o "%USERPROFILE%\.harmony\bin\HarmonyLang.exe"

:: Create harmony.bat wrapper script
echo Creating wrapper script...
(
echo @echo off
echo IF "%%~1"=="run" ^(
echo     IF NOT "%%~2"=="" ^(
echo         IF EXIST "%%~2" ^(
echo             "%USERPROFILE%\.harmony\bin\HarmonyLang.exe" "%%~2"
echo         ^) ELSE ^(
echo             echo Error: File %%~2 not found
echo             exit /b 1
echo         ^)
echo     ^) ELSE ^(
echo         echo Error: No file specified
echo         exit /b 1
echo     ^)
echo ^) ELSE IF "%%~1"=="repl" ^(
echo     "%USERPROFILE%\.harmony\bin\HarmonyLang.exe"
echo ^) ELSE ^(
echo     echo Usage:
echo     echo   harmony run ^<file.harmony^>  - Run a Harmony source file
echo     echo   harmony repl               - Start Harmony REPL
echo ^)
) > "%USERPROFILE%\.harmony\bin\harmony.bat"

:: Add to PATH if not already there
echo Adding to PATH...
setx PATH "%PATH%;%USERPROFILE%\.harmony\bin"

echo.
echo Installation complete!
echo Please restart your Command Prompt to use the 'harmony' command
echo.
echo Usage:
echo   harmony run ^<file.harmony^>  - Run a Harmony source file
echo   harmony repl               - Start Harmony REPL

pause