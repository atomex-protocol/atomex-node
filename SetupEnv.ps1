$envFilePath = "./.env"
$isEnvExist = Test-Path -Path $envFilePath

if (-not($isEnvExist)) {
	Write-Host "Should be initialized .env file with requred variables first."
	Write-Host "Docs: https://github.com/atomex-protocol/atomex-node#environment-variables"

	return
}

$variablesPattern = '(AP_ENV|ETHEREUM_PRIVATE|TEZOS_PRIVATE)=\S'
$wihoutErrors = $true

Get-Content $envFilePath | ?{ $_ -match $variablesPattern } | % {
	$line = $_

	try {
		$pair = $line -Split "="
		$variable, $value = $pair[0], $pair[1]

		[Environment]::SetEnvironmentVariable($variable, $value, [EnvironmentVariableTarget]::Process)
	} catch {
		"Error with handling .env line: " + $line + [Environment]::Newline + $_
		$wihoutErrors = $false
	}
}

if ($wihoutErrors) {
	Write-Host "Atomex Node environment setup successfully done."
}
