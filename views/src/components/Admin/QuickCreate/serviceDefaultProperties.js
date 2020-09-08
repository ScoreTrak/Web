const WINRM = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Port: {name: 'Port', type: 'field', defaultValue: '5985'},
    Command: {name: 'Command', type: 'field', defaultValue: 'whoami' },
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
    Scheme: {name: 'Scheme', type: 'select', defaultValue: 'http', options: ["http", "https"]},
    ClientType: {name: 'Client Type', type: 'select', defaultValue: 'NTLM', options: ["NTLM"]},
}

const SSH = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Port: {name: 'Port', type: 'field', defaultValue: '22'},
    Command: {name: 'Command', type: 'field', defaultValue: 'whoami'},
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
}

const SMB = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Domain: {name: 'Domain', type: 'field', },
    Port: {name: 'Port', type: 'field', defaultValue: '445'},
    TransportProtocol: {name: 'Transport Protocol', type: 'field', defaultValue: 'tcp'},
    Share: {name: 'Share', type: 'field', },
    FileName: {name: 'FileName', type: 'field', defaultValue: 'TestFile.txt'},
    Text: {name: 'Text', type: 'field', defaultValue: 'Hello World!'},
    Operation: {name: 'Operation', type: 'select', defaultValue: 'CreateAndOpen', options: ["Open", "Create", "CreateAndOpen"]},
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
}

const PING = {
    Protocol: {name: 'Protocol', type: 'select', options:["ipv4", "ipv6"], defaultValue: 'ipv4' },
    Attempts: {name: 'Attempts', type: 'field', defaultValue: '1'},
}

const LDAP = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Domain: {name: 'Domain', type: 'field', },
    Port: {name: 'Port', type: 'field', defaultValue: '389'},
    TransportProtocol: {name: 'Transport Protocol', type: 'field', defaultValue: 'tcp'},
    BaseDN: {name: 'Base DN', type: 'field', },
    ApplicationProtocol: {name: 'Application Protocol', type: 'select', defaultValue: 'ldap', options: ["ldap", "ldaps"]},
    Filter: {name: 'Filter', type: 'field', defaultValue: '(&(objectClass=organizationalPerson))'},
    Attributes: {name: 'Attributes', type: 'field', defaultValue: 'dn,cn' }
}

const IMAP = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Port: {name: 'Port', type: 'field', defaultValue: '143'},
    Scheme: {name: 'Scheme', type: 'select', defaultValue: 'imap', options: ["imap", "tls"]},
}

const HTTP = {
    Port: {name: 'Port', type: 'field', defaultValue: '80'},
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
    Scheme: {name: 'Scheme', type: 'select', defaultValue: 'http', options: ["http", "https"]},
    Path: {name: 'Path', type: 'field', },
    Subdomain: {name: 'Subdomain', type: 'field', }
}

const FTP = {
    Username: {name: 'Username', type: 'field', },
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Port: {name: 'Port', type: 'field', defaultValue: '21'},
    Text: {name: 'Text', type: 'field', },
    ReadFilename: {name: 'Read File Name', type: 'field', },
    WriteFilename: {name: 'Write File Name', type: 'field', },
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
}

const DNS = {
    Lookup: {name: 'Lookup', type: 'field', }, 
    ExpectedOutput: {name: 'Expected Output', type: 'field', },
}

const SQL = {
    Username: {name: 'Username', type: 'field'},
    Password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    Port: {name: 'Port', type: 'field', defaultValue: '3306'},
    DBType: {name: 'Database Type', type: 'select', options: ["mysql", "postgres"]},
    DBName: {name: 'Database Name', type: 'field'},
    Command: {name: 'Command', type: 'field', defaultValue: '' },
    MinExpectedRows: {name: 'Minimum Expected Rows', type: 'field'},
    MaxExpectedRows: {name: 'Maximum Expected Rows', type: 'field' },
}


export default {
    WINRM, HTTP, SSH, SMB, DNS, FTP, IMAP, LDAP, PING, SQL
}