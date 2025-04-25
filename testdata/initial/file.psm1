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
