set GME_VER=%1
set DESTDIR=%2

set URL="http://blargg.8bitalley.com/parodius/libs/Game_Music_Emu-%GME_VER%.zip"
set DEST=.\%DESTDIR%\gme.zip

rmdir /S /Q %DESTDIR%
mkdir %DESTDIR%
powershell -c wget %URL% -OutFile "%DEST%"
powershell -c Expand-Archive -Path %DEST%

