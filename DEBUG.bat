@echo off
echo ==========================================
echo WebRTC Connection Debugger
echo ==========================================
echo.

echo This will help diagnose why connection is stuck.
echo.
echo Please open browser to http://localhost:8080
echo Then press F12 to open Developer Console
echo.
echo In the Console tab, paste this and press Enter:
echo.
echo -----------------------------------------
echo console.log('Checking WebRTC...');
echo navigator.mediaDevices.getUserMedia({audio: true})
echo   .then(() => console.log('MIC OK'))
echo   .catch(e => console.log('MIC ERROR:', e));
echo -----------------------------------------
echo.
echo What does it say?
echo   - If "MIC OK" = microphone works
echo   - If "MIC ERROR" = microphone blocked
echo.
pause
echo.

echo Now check if the /offer endpoint works:
echo.
echo In the same Console tab, paste this:
echo.
echo -----------------------------------------
echo fetch('/offer', {method: 'POST', headers: {'Content-Type': 'application/json'}, body: JSON.stringify({sdp: 'test'})})
echo   .then(r => r.text()).then(console.log).catch(console.error);
echo -----------------------------------------
echo.
echo What does it show?
echo   - If you see JSON with "sdp" = server is responding
echo   - If error = server issue
echo.
pause
echo.

echo ==========================================
echo Common Fixes:
echo ==========================================
echo.
echo 1. Browser Permissions
echo    - Click the lock icon in address bar
echo    - Make sure Microphone is "Allow"
echo.
echo 2. Try Different Browser
echo    - Chrome or Edge work best
echo    - Firefox also good
echo.
echo 3. Restart Gateway
echo    - Press Ctrl+C in the gateway window
echo    - Run SIMPLE_START.bat again
echo.
echo 4. Check Windows Firewall
echo    - Allow when prompted
echo.
pause
