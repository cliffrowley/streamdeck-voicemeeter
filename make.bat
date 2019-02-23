@ECHO off

call env.bat
call clean.bat

go build -o %__plugin_dir%/streamdeck-voicemeeter.exe
if not ERRORLEVEL 0 (
  echo Build failed
  exit /b 1
)

DistributionTool com.github.cliffrowley.streamdeck-voicemeeter.sdPlugin release
if not ERRORLEVEL 0 (
  echo Release failed
  exit /b 1
)
