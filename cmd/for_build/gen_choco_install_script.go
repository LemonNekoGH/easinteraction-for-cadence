package main

import (
	"fmt"
	"os"
)

// ONLY used to generate package.json
// usage: go run gen_choco_install_script <version> <checksum>
func main() {
	version := os.Args[1]
	checksum := os.Args[2]

	template := `$ErrorActionPreference = 'Stop';

$url64 = 'https://github.com/LemonNekoGH/easinteraction-for-cadence/releases/download/%s/easi-gen-windows-amd64.exe'
checksum64 = '%s'
$checksumType64 = 'sha256'
$File = Join-Path (Join-Path $env:ChocolateyInstall (Join-Path 'lib' $env:ChocolateyPackageName)) 'easi-gen.exe'
$file = Get-ChocolateyWebFile -PackageName $env:ChocolateyPackageName -FileFullPath $File -Url64bit $url64 -CheckSum64 $checksum64 -CheckSumType64 $checksumType64
`
	fmt.Printf(template, version, checksum)
}
