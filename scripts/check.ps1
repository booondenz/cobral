# Self-check for the cidr fill-in-the-blank exercise.
#
#   powershell -ExecutionPolicy Bypass -File scripts\check.ps1
#
# It builds the project, runs every command and compares the real output with
# the output of the original project (github.com/bschaatsbergen/cidr), so it
# tells you exactly which blanks are still missing.
#
# This file is intentionally plain ASCII: keep it like that or Windows
# PowerShell 5.1 may misread it.

[CmdletBinding()]
param(
    [string]$RepoRoot = ''
)

if (-not $RepoRoot) {
    $RepoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
}

# 'Continue' on purpose: this script does its own error reporting and native
# tools (go build, the CLI itself) write to stderr, which PowerShell 5.1 turns
# into ErrorRecords. With 'Stop' those would abort the run half way.
$ErrorActionPreference = 'Continue'

# fatih/color prints ANSI escapes only on a real terminal; NO_COLOR keeps the
# output stable no matter where we run this.
$env:NO_COLOR = '1'

function Get-NormalizedOutput {
    param([object[]]$Lines)

    $parts = @()
    foreach ($line in $Lines) {
        if ($null -eq $line) { continue }
        # collapse every run of whitespace (tabs included) into one space
        $text = ("$line" -replace '\s+', ' ').Trim()
        # noise injected by PowerShell when a native tool writes to stderr
        if ($text -eq 'System.Management.Automation.RemoteException') { continue }
        if ($text -ne '') { $parts += $text }
    }
    return ($parts -join "`n")
}

# 51 expected lines: 25 subnets, the ellipsis, the last 25 subnets.
$divideTruncated = @()
for ($i = 0; $i -lt 25; $i++) { $divideTruncated += "10.0.$($i * 4).0/22" }
$divideTruncated += '......'
for ($i = 35; $i -lt 60; $i++) { $divideTruncated += "10.0.$($i * 4).0/22" }

$cases = @(
    @{ name = 'L1  root help lists every subcommand'; cmd = @('--help')
       containsAll = @('Available Commands:', 'cidr', 'count', 'contains', 'overlaps', 'explain', 'divide') }

    @{ name = 'L3  count 10.0.0.0/16';               cmd = @('count', '10.0.0.0/16')
       want = @('65536') }
    @{ name = 'L3  count 2001:db8:1234:1a00::/106';  cmd = @('count', '2001:db8:1234:1a00::/106')
       want = @('4194304') }
    @{ name = 'L3  count 172.16.18.0/31';            cmd = @('count', '172.16.18.0/31')
       want = @('2') }
    @{ name = 'L3  count 172.16.18.0/32';            cmd = @('count', '172.16.18.0/32')
       want = @('1') }

    @{ name = 'L3  contains 10.0.0.0/16 10.0.14.5';  cmd = @('contains', '10.0.0.0/16', '10.0.14.5')
       want = @('true') }
    @{ name = 'L3  contains 10.0.0.0/16 10.1.14.5';  cmd = @('contains', '10.0.0.0/16', '10.1.14.5')
       want = @('false') }

    @{ name = 'L5  overlaps 10.0.0.0/16 10.0.14.0/22'; cmd = @('overlaps', '10.0.0.0/16', '10.0.14.0/22')
       want = @('true') }
    @{ name = 'L5  overlaps 10.0.0.0/16 10.1.0.0/16';  cmd = @('overlaps', '10.0.0.0/16', '10.1.0.0/16')
       want = @('false') }

    @{ name = 'L4  explain 10.1.0.0/16'; cmd = @('explain', '10.1.0.0/16')
       want = @(
           'Base Address: 10.1.0.0',
           'Usable Address Range: 10.1.0.1 to 10.1.255.254 (65,534)',
           'Broadcast Address: 10.1.255.255',
           'Addresses: 65,536',
           'Netmask: 255.255.0.0 (/16 bits)'
       ) }
    @{ name = 'L4  explain 2001:db8:1234:1a00::/110'; cmd = @('explain', '2001:db8:1234:1a00::/110')
       want = @(
           'Base Address: 2001:db8:1234:1a00::',
           'Usable Address Range: 2001:db8:1234:1a00:: to 2001:db8:1234:1a00::3:ffff (262,142)',
           'Addresses: 262,144',
           'Netmask: ffff:ffff:ffff:ffff:ffff:ffff:fffc:0 (/110 bits)'
       ) }
    @{ name = 'L4  explain 172.16.18.0/31 (edge case)'; cmd = @('explain', '172.16.18.0/31')
       want = @(
           'Base Address: 172.16.18.0',
           'Usable Address Range: 172.16.18.0 to 172.16.18.1 (2)',
           'Broadcast Address: IPv4 network has no broadcast address',
           'Addresses: 2',
           'Netmask: 255.255.255.254 (/31 bits)'
       ) }

    @{ name = 'L5  divide 10.0.0.0/16 9'; cmd = @('divide', '10.0.0.0/16', '9')
       want = @(
           '10.0.0.0/20', '10.0.16.0/20', '10.0.32.0/20', '10.0.48.0/20',
           '10.0.64.0/20', '10.0.80.0/20', '10.0.96.0/20', '10.0.112.0/20', '10.0.128.0/20'
       ) }
    @{ name = 'L5  divide 2001:db8:1111:2222:1::/80 9'; cmd = @('divide', '2001:db8:1111:2222:1::/80', '9')
       want = @(
           '2001:db8:1111:2222:1::/84', '2001:db8:1111:2222:1:1000::/84',
           '2001:db8:1111:2222:1:2000::/84', '2001:db8:1111:2222:1:3000::/84',
           '2001:db8:1111:2222:1:4000::/84', '2001:db8:1111:2222:1:5000::/84',
           '2001:db8:1111:2222:1:6000::/84', '2001:db8:1111:2222:1:7000::/84',
           '2001:db8:1111:2222:1:8000::/84'
       ) }
    @{ name = 'L5  divide 10.0.0.0/16 60 (truncated to 50)'; cmd = @('divide', '10.0.0.0/16', '60')
       want = $divideTruncated }

    @{ name = 'err count without args';        cmd = @('count');                       contains = 'error: provide a CIDR range';               wantFail = $true }
    @{ name = 'err count with a bad CIDR';     cmd = @('count', 'notacidr');           contains = 'error: invalid CIDR range: notacidr';        wantFail = $true }
    @{ name = 'err contains with a bad IP';    cmd = @('contains', '10.0.0.0/16', 'notanip'); contains = 'error: invalid IP address: notanip';   wantFail = $true }
    @{ name = 'err overlaps with one arg';     cmd = @('overlaps', '10.0.0.0/16');     contains = 'error: provide 2 CIDR ranges';               wantFail = $true }
    @{ name = 'err explain without args';      cmd = @('explain');                     contains = 'error: provide a CIDR range';               wantFail = $true }
    @{ name = 'err divide with a bad network'; cmd = @('divide', 'abc', '9');          contains = 'invalid network: abc';                      wantFail = $true }
    @{ name = 'err divide with a bad divisor'; cmd = @('divide', '10.0.0.0/16', 'x');  contains = 'invalid divisor: x';                        wantFail = $true }
)

$exe = Join-Path $RepoRoot 'cidr-check.exe'

Push-Location $RepoRoot
try {
    Write-Host ''
    Write-Host 'building ...' -ForegroundColor DarkGray
    go build -o $exe . 2>&1 | ForEach-Object { Write-Host "  $_" -ForegroundColor DarkRed }
    if ($LASTEXITCODE -ne 0) {
        Write-Host ''
        Write-Host 'BUILD FAILED. Fix the compile errors first (the compiler tells you which file/line).' -ForegroundColor Red
        exit 1
    }

    $passed = 0
    $failed = 0

    foreach ($case in $cases) {
        $cmdArgs = $case.cmd
        $raw = & $exe @cmdArgs 2>&1
        $exitCode = $LASTEXITCODE
        $got = Get-NormalizedOutput -Lines @($raw)

        $problems = @()

        if ($case.ContainsKey('want')) {
            $expected = (@($case.want) | ForEach-Object { ("$_" -replace '\s+', ' ').Trim() }) -join "`n"
            if ($got -ne $expected) {
                $problems += "output mismatch"
                $problems += "    expected : $($expected -replace "`n", ' | ')"
                $problems += "    actual   : $($got -replace "`n", ' | ')"
            }
        }

        if ($case.ContainsKey('contains') -and -not $got.Contains($case.contains)) {
            $problems += "output does not contain '$($case.contains)'"
        }

        if ($case.ContainsKey('containsAll')) {
            foreach ($needle in $case.containsAll) {
                if (-not $got.Contains($needle)) { $problems += "output does not contain '$needle'" }
            }
        }

        if ($case.ContainsKey('wantFail')) {
            if ($exitCode -eq 0) { $problems += 'expected a non-zero exit code, got 0' }
        } elseif ($exitCode -ne 0) {
            $problems += "expected exit code 0, got $exitCode"
        }

        if ($problems.Count -eq 0) {
            $passed++
            Write-Host ("  PASS  " + $case.name) -ForegroundColor Green
        } else {
            $failed++
            Write-Host ("  FAIL  " + $case.name) -ForegroundColor Red
            foreach ($problem in $problems) { Write-Host "        $problem" -ForegroundColor DarkRed }
        }
    }

    Write-Host ''
    Write-Host "  $passed passed, $failed failed" -ForegroundColor $(if ($failed -eq 0) { 'Green' } else { 'Yellow' })
    Write-Host ''
    if ($failed -gt 0) { exit 1 }
}
finally {
    Pop-Location
    Remove-Item -Force $exe -ErrorAction SilentlyContinue
}
