import React from "react";
import {Route} from "react-router-dom";
import TeamService from "../../services/team/teams";
import ServiceGroupsService from "../../services/service_group/service_groups";
import HostGroupsService from "../../services/host_group/host_groups";
import RoundsService from "../../services/round/round";
import UserService from "../../services/users/users";
import {Table} from "./TableInterface"
import TeamMenu from "./TeamMenu";
import Paper from "@material-ui/core/Paper";
import HostMenu from "./HostMenu";
import ServiceMenu from "./ServiceMenu";
import PropertiesMenu from "./PropertiesMenu";

export default function Setup(props) {
    return (
        <Paper className={props.classesPaper} style={{minHeight: "85vh"}}>
           <Route exact path="/setup/teams" render={() => (
               <TeamMenu {...props} />
           )} />
           <Route exact path="/setup/host_groups" render={() => (
               <React.Fragment>
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
               </React.Fragment>
           )} />
           <Route exact path="/setup/users" render={() => (
               <React.Fragment>
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
               </React.Fragment>
           )} />

           <Route exact path="/setup/hosts" render={() => (
               <HostMenu {...props} />
           )} />

           <Route exact path="/setup/services" render={() => (
               <ServiceMenu {...props} />
           )} />
           <Route exact path="/setup/properties" render={() => (
               <PropertiesMenu {...props}/>
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
                       const additionalActions = [{icon: "replay", tooltip: 'redeploy workers', onFuncClick: async (event, rowData) => {
                           props.handleLoading()
                               return await ServiceGroupsService.Redeploy(rowData["id"]).then(() => {
                                  props.handleSuccess("Workers were deployed! Please make sure they are in a healthy state before enabling the service group.")
                               })

                           } }]
                       const columns=
                           [
                               { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                               { title: 'Service Group Name', field: 'name', editable: 'onAdd' },
                               { title: 'Enabled', field: 'enabled', type: 'boolean'},
                               { title: 'Skip Helper(Skips autodeploy of workers)', field: 'skip_helper', type: 'boolean' },
                               { title: 'Label(Workers would be deployed on nodes with the following label)', field: 'label', editable: 'onAdd'},
                           ]

                       return Table(props, title, isDependant, columns, true, ServiceGroupsService, [], [], [],true,additionalActions)
                   })()}
               </div>
           )} />
        </Paper>
    );
}

//ToDo: Allow Chaining of the foreign IDs
//Todo: Add ID propagation all the way to Team Element