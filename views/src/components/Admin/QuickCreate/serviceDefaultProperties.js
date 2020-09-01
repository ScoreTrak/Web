const winrm = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    port: {name: 'Port', type: 'field', defaultValue: '5985'},
    command: {name: 'Command', type: 'field', },
    'expected_output': {name: 'Expected Output', type: 'field', },
    scheme: {name: 'Scheme', type: 'select', defaultValue: 'http', options: ["http", "https"]},
    client_type: {name: 'Client Type', type: 'select', defaultValue: 'NTLM', options: ["NTLM"]},
}

const ssh = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    port: {name: 'Port', type: 'field', defaultValue: '22'},
    command: {name: 'Command', type: 'field', defaultValue: 'whoami'},
    'expected_output': {name: 'Expected Output', type: 'field', },
}

const smb = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    domain: {name: 'Domain', type: 'field', },
    port: {name: 'Port', type: 'field', defaultValue: '445'},
    'transport_protocol': {name: 'Transport Protocol', type: 'field', defaultValue: 'tcp'},
    share: {name: 'Share', type: 'field', },
    file_name: {name: 'FileName', type: 'field', defaultValue: 'TestFile.txt'},
    text: {name: 'Text', type: 'field', defaultValue: 'Hello World!'},
    operation: {name: 'Operation', type: 'select', defaultValue: 'CreateAndOpen', options: ["Open", "Create", "CreateAndOpen"]},
    'expected_output': {name: 'Expected Output', type: 'field', },
}

const ping = {
    protocol: {name: 'Protocol', type: 'select', options:["ipv4", "ipv6"], defaultValue: 'ipv4' },
    attempts: {name: 'Attempts', type: 'field', defaultValue: '1'},
}

const ldap = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    domain: {name: 'Domain', type: 'field', },
    port: {name: 'Port', type: 'field', defaultValue: '389'},
    'transport_protocol': {name: 'Transport Protocol', type: 'field', defaultValue: 'tcp'},
    base_dn: {name: 'Base DN', type: 'field', },
    'application_protocol': {name: 'Application Protocol', type: 'select', defaultValue: 'ldap', options: ["ldap", "ldaps"]},
    filter: {name: 'Filter', type: 'field', defaultValue: '(&(objectClass=organizationalPerson))'},
    attributes: {name: 'Attributes', type: 'field', defaultValue: 'dn,cn' }
}

const imap = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    port: {name: 'Port', type: 'field', defaultValue: '143'},
    scheme: {name: 'Scheme', type: 'select', defaultValue: 'imap', options: ["imap", "tls"]},
}

const http = {
    port: {name: 'Port', type: 'field', defaultValue: '80'},
    'expected_output': {name: 'Expected Output', type: 'field', },
    scheme: {name: 'Scheme', type: 'select', defaultValue: 'http', options: ["http", "https"]},
    path: {name: 'Path', type: 'field', },
    subdomain: {name: 'Subdomain', type: 'field', }
}

const ftp = {
    username: {name: 'Username', type: 'field', },
    password: {name: 'Password', type: 'field',  defaultStatus: 'Edit'},
    port: {name: 'Port', type: 'field', defaultValue: '21'},
    text: {name: 'Text', type: 'field', },
    read_file: {name: 'Read File Name', type: 'field', },
    write_file: {name: 'Write File Name', type: 'field', },
    'expected_output': {name: 'Expected Output', type: 'field', },
}

const dns = {
    lookup: {name: 'Lookup', type: 'field', },
   'expected_output': {name: 'Expected Output', type: 'field', },
}

export default {
    winrm, http, ssh, smb, dns, ftp, imap, ldap, ping
}