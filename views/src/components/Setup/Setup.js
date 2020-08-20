import React from "react";
import {Route} from "react-router-dom";
import TeamService from "../../services/team/teams";
import ServiceGroupsService from "../../services/service_group/service_groups";
import ServicesService from "../../services/service/serivces";
import HostGroupsService from "../../services/host_group/host_groups";
import HostsService from "../../services/host/hosts";
import PropertyService from "../../services/property/properties";
import RoundsService from "../../services/round/round";
import UserService from "../../services/users/users";

import EditableTable from "./EditableTable/EditableTable";
function Table(props, title, isDependant, columns, disallowBulkAdd, service, owningService, fieldForLookup, owningFieldLookup=["id"], editable=true){
    return <EditableTable {...props} title={title} isDependant={isDependant} columns={columns} service={service} disallowBulkAdd={disallowBulkAdd} owningService={owningService}
                          fieldForLookup={fieldForLookup} owningFieldLookup={owningFieldLookup} editable={editable}
    />
}

export default function Setup(props) {
    return (
       <div>
           <Route exact path="/setup/teams" render={() => (
               <div>
                   {(() => {
                       const title = "Teams"
                       const isDependant = false
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Team Name', field: 'name' },
                               { title: 'Index', field: 'index', type: 'numeric' },
                               { title: 'Enabled', field: 'enabled', type: 'boolean' },
                           ]

                       return Table(props, title, isDependant, columns, false, TeamService)
                   })()}
               </div>
           )} />
           <Route exact path="/setup/host_groups" render={() => (
               <div>
                   {(() => {
                       const title = "Host Groups"
                       const isDependant = false
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Host Group Name', field: 'name' },
                               { title: 'Enabled', field: 'enabled', type: 'boolean' },
                           ]

                       return Table(props, title, isDependant, columns, false, HostGroupsService)
                   })()}
               </div>
           )} />
           <Route exact path="/setup/users" render={() => (
               <div>
                   {(() => {
                       const title = "Users"
                       const isDependant = true
                       const owningService = [TeamService]
                       const owningFieldLookup = ["name"]
                       const fieldForLookup = ["team_id"]
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Username', field: 'username' },
                               { title: 'Password', field: 'password' },
                               { title: 'Password Hash', field: 'password_hash', editable: 'never'},
                               { title: 'Team ID', field: 'team_id' },
                               { title: 'Role', field: 'role', lookup: { 'black': 'black', 'blue': 'blue' }},
                           ]

                       return Table(props, title, isDependant, columns, false, UserService, owningService, fieldForLookup, owningFieldLookup)
                   })()}
               </div>
           )} />

           <Route exact path="/setup/hosts" render={() => (
               <div>
                   {(() => {
                       const title = "Hosts"
                       const isDependant = true
                       const owningService = [TeamService, HostGroupsService]
                       const owningFieldLookup = ["name", "name"]
                       const fieldForLookup = ["team_id", "host_group_id"]
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Address', field: 'address' },
                               { title: 'Host Group ID', field: 'host_group_id' },
                               { title: 'Team ID', field: 'team_id' },
                               { title: 'Enabled', field: 'enabled', type: 'boolean' },
                               { title: 'Edit Host(Allow users to change Addresses)', field: 'edit_host', type: 'boolean' },
                           ]

                       return Table(props, title, isDependant, columns, false, HostsService, owningService, fieldForLookup, owningFieldLookup)
                   })()}
               </div>
           )} />

           <Route exact path="/setup/services" render={() => (
               <div>
                   {(() => {
                       const title = "Services"
                       const isDependant = true
                       const owningService = [ServiceGroupsService, HostsService]
                       const owningFieldLookup = ["name", "address"]
                       const fieldForLookup = ["service_group_id", "host_id"]
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Name', field: 'name', lookup: {
                                   'PING': 'PING', 'DNS':'DNS', 'FTP':'FTP', 'LDAP':'LDAP',
                                       'HTTP': 'HTTP', 'IMAP': 'IMAP', 'SMB': 'SMB', 'SSH': 'SSH',
                                       'WINRM': 'WINRM'
                               }},
                               { title: 'Display Name(Columns on Status page)', field: 'display_name' },
                               { title: 'Points(Points per successful check)', field: 'points', type: 'numeric', },
                               { title: 'Points Boost', field: 'points_boost', type: 'numeric', initialEditValue: 0},
                               { title: 'Enabled', field: 'enabled', type: 'boolean' },
                               { title: 'Service Group ID', field: 'service_group_id' },
                               { title: 'Host ID', field: 'host_id' },
                               { title: 'Round Units(Frequency)', field: 'round_units', type: 'numeric', initialEditValue: 1},
                               { title: 'Round Delay(Shift in frequency)', field: 'round_delay', type: 'numeric', initialEditValue: 0 },
                           ]

                       return Table(props, title, isDependant, columns, false, ServicesService, owningService, fieldForLookup, owningFieldLookup)
                   })()}
               </div>
           )} />
           <Route exact path="/setup/properties" render={() => (
               <div>
                   {(() => {
                       const title = "Properties"
                       const isDependant = true
                       const owningService = [ServicesService]
                       const owningFieldLookup = ["id"]
                       const fieldForLookup = ["service_id"]
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Key', field: 'key' },
                               { title: 'Value', field: 'value' },
                               { title: 'Status', field: 'status', lookup:{'View': 'View', 'Hide':'Hide', 'Edit':'Edit'}},
                               { title: 'Description', field: 'description'},
                               { title: 'Service ID', field: 'service_id'},
                           ]
                        //ToDo: Show required properties for a given service
                       return Table(props, title, isDependant, columns, false, PropertyService, owningService, fieldForLookup, owningFieldLookup)
                   })()}
               </div>
           )} />

           <Route exact path="/setup/rounds" render={() => (
               <div>
                   {(() => {
                       const title = "Rounds"
                       const isDependant = false
                       const owningService = []
                       const owningFieldLookup = []
                       const fieldForLookup = []
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Start Time', field: 'start', type:'datetime' },
                               { title: 'Note', field: 'note' },
                               { title: 'Error', field: 'err'},
                               { title: 'Finish Time', field: 'finish', type:'datetime'},
                           ]
                       return Table(props, title, isDependant, columns, false, RoundsService, owningService, fieldForLookup, owningFieldLookup, false)
                   })()}
               </div>
           )} />

           <Route exact path="/setup/service_groups" render={() => (
               <div>
                   {(() => {
                       const title = "Service Groups"
                       const isDependant = false
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Service Group Name', field: 'name', editable: 'onAdd' },
                               { title: 'Enabled', field: 'enabled', type: 'boolean'},
                               { title: 'Skip Helper(Skips autodeploy of workers)', field: 'skip_helper', type: 'boolean', editable: 'onAdd' },
                               { title: 'Label(Workers would be deployed on nodes with the following label)', field: 'label', editable: 'onAdd'},
                           ]
                       //ToDo: FigureOut Worker Reload functionality https://material-table.com/#/docs/features/actions
                       //ToDo: Allow Same for Service testing https://material-table.com/#/docs/features/actions

                       return Table(props, title, isDependant, columns, true, ServiceGroupsService)
                   })()}
               </div>
           )} />
       </div>
    );
}

//ToDo: Allow Chaining of the foreign IDs

//Todo: Add ID propagation all the way to Team Element