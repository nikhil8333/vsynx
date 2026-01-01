# Debug Guide - Vsynx Manager

## Recent Changes

### 1. Fixed Import Path Issue
- **Problem**: Frontend was importing from `../wailsjs/go/main/App` (incorrect parent directory)
- **Solution**: Changed to `./wailsjs/go/main/App` (correct relative path)
- **File**: `frontend/src/App.tsx` line 10

### 2. Fixed Wails Configuration
- **Problem**: `wailsjsdir` was set to `./frontend/src/wailsjs`, causing double-nested directory structure
- **Solution**: Changed to `./frontend/src` to generate bindings correctly
- **File**: `wails.json` line 20

### 3. Added Comprehensive Logging
All backend components now have detailed logging:
- **App methods** (`app.go`): Logs when methods are called from frontend
- **Validator** (`internal/validation/validator.go`): Logs validation steps
- **Scanner** (`internal/validation/scanner.go`): Logs extension scanning
- **Marketplace Client** (`internal/marketplace/client.go`): Logs API calls to Microsoft Marketplace
- **OpenVSX Client** (`internal/openvsx/client.go`): Logs API calls to OpenVSX

### 4. Enhanced Frontend Error Handling
- Added error state and error banner UI
- Added console logging for all API calls
- Better user feedback for failures

## How to Debug API Calls

### Check Backend Logs
When you run `wails dev`, watch the terminal for log messages:

```
[App] GetDefaultExtensionsPath called
[App] Default extensions path: /home/user/.vscode/extensions
[Scanner] Scanning extensions at: /home/user/.vscode/extensions
[Scanner] Found 42 entries in directory
[Scanner] Successfully scanned 10 extensions
[App] Found 10 extensions
```

### Check Frontend Console
Open the browser DevTools (F12) and look for frontend logs:

```
[Frontend] Loading default extensions path...
[Frontend] Default path: /home/user/.vscode/extensions
[Frontend] Loading extensions from: /home/user/.vscode/extensions
[Frontend] Loaded extensions: 10
```

### Validation Logs
When you click "Validate" on an extension:

```
[App] ValidateExtension called for: ms-python.python
[Validator] Starting validation for extension: ms-python.python
[Marketplace] Fetching metadata for extension: ms-python.python
[Marketplace] Response status for ms-python.python: 200
[Marketplace] Successfully fetched metadata for ms-python.python (version 2024.0.1)
[OpenVSX] Fetching metadata for extension: ms-python.python
[OpenVSX] Response status for ms-python.python: 200
[OpenVSX] Successfully fetched metadata for ms-python.python (version 2024.0.1)
[Validator] Validation complete for ms-python.python: Legitimate
```

## Common Issues & Solutions

### 1. No Extensions Found
**Symptom**: App loads but shows 0 extensions

**Debug Steps**:
1. Check terminal logs for `[Scanner] Scanning extensions at:`
2. Verify the path is correct
3. Check `[Scanner] Found X entries in directory`
4. Look for errors in scanning

**Solution**: Use "Change Path" button to select correct VS Code extensions directory

### 2. API Calls Not Working
**Symptom**: Validation shows no results or errors

**Debug Steps**:
1. Check for `[Marketplace] Request failed` or `[OpenVSX] Request failed` in logs
2. Verify internet connection
3. Check if APIs are accessible:
   - Marketplace: `https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery`
   - OpenVSX: `https://open-vsx.org/api`

**Solution**: Check network connectivity, firewall settings, or proxy configuration

### 3. Extension Not Found in Marketplace
**Symptom**: `extension not found in marketplace` error

**Possible Causes**:
- Extension ID format is incorrect (should be `publisher.name`)
- Extension is only available on OpenVSX
- Extension has been delisted

### 4. Time.Time Binding Warnings
**Symptom**: `Not found: time.Time` warnings during `wails dev`

**Note**: These are warnings from Wails binding generation and don't affect functionality. The `time.Time` fields are properly serialized as JSON strings.

## Testing the Application

### Test 1: Load Extensions
1. Start app with `wails dev`
2. Check terminal for path detection
3. Verify extensions list populates
4. Look for frontend console logs

### Test 2: Validate Single Extension
1. Select an extension from the list
2. Click "Validate" button
3. Watch terminal for API calls
4. Verify result shows trust level and recommendation

### Test 3: Full Audit
1. Click "Audit" button in header
2. Watch terminal for progress
3. Verify summary cards show counts
4. Check detailed results list

### Test 4: Download Official Extension
1. Select an extension
2. Click "Download Official"
3. Check terminal for download progress
4. Verify success/error message

## API Endpoints Used

### Microsoft Marketplace API
- **URL**: `https://marketplace.visualstudio.com/_apis/public/gallery/extensionquery`
- **Method**: POST
- **API Version**: 7.0 (required in Accept header)
- **Purpose**: Fetch official extension metadata
- **Headers Required**:
  - `Content-Type: application/json`
  - `Accept: application/json; api-version=7.0`
  - `User-Agent: Vsynx/1.0`
- **Request Format**: JSON with filter criteria
- **Response**: Extension details, version info, download URLs

### OpenVSX Registry API
- **URL**: `https://open-vsx.org/api/{namespace}/{extension}`
- **Method**: GET
- **Purpose**: Fetch open-source registry metadata
- **Response Format**: JSON with extension details

## Next Steps for Testing

1. **Test with real extensions**: Use your actual VS Code installation
2. **Verify API responses**: Check that Microsoft and OpenVSX data matches
3. **Test error cases**: Try with invalid extension IDs
4. **Performance**: Test with large number of extensions (100+)
5. **Network failures**: Test behavior when APIs are unreachable

## Troubleshooting Commands

### Rebuild bindings
```bash
cd d:\code\secureopenvsx
wails dev
```

### Check Go dependencies
```bash
go mod tidy
go mod verify
```

### Check frontend dependencies
```bash
cd frontend
npm install
npm run build
```

### View Wails logs
The terminal running `wails dev` shows all backend logs in real-time.

### Access browser console
1. App running at `http://localhost:34115` (dev mode)
2. Press F12 to open DevTools
3. Check Console tab for frontend logs
4. Check Network tab for failed requests
