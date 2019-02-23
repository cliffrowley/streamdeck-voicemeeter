@ECHO off

call env.bat

IF EXIST %__plugin_bin% del %__plugin_bin%
IF EXIST %__release_bin% del %__release_bin%
