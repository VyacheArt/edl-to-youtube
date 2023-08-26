//go:generate fyne bundle --pkg resources --name ResolveManualMedia -o bundled.go "manual/resolve/1. Media.png"
//go:generate fyne bundle --pkg resources --name ResolveManualExport -o bundled.go -append "manual/resolve/2. Export.png"
//go:generate fyne bundle --pkg resources --name ResolveManualSave -o bundled.go -append "manual/resolve/3. Save.png"
package resources
