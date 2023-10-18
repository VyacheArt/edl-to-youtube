/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

//go:generate fyne bundle --pkg resources --name ResolveManualMedia -o bundled.go "manual/resolve/1. Media.png"
//go:generate fyne bundle --pkg resources --name ResolveManualExport -o bundled.go -append "manual/resolve/2. Export.png"
//go:generate fyne bundle --pkg resources --name ResolveManualSave -o bundled.go -append "manual/resolve/3. Save.png"
package resources
