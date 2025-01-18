<#
 Copyright 2025 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
#>

function Get-HelloWorld {
    <#
    .SYNOPSIS
        Outputs a "Hello, World!" message.

    .DESCRIPTION
        This function is a basic example to demonstrate how to create a PowerShell module.

    .EXAMPLE
        PS> Get-HelloWorld
        Hello, World!

    .NOTES
        This is a demo function.
    #>
    Write-Output "Hello, World!"
}

# Export the function to make it available when the module is imported.
Export-ModuleMember -Function Get-HelloWorld
