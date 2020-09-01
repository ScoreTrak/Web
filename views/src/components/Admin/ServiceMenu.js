import React from "react";
import {Table} from "./TableInterface";
import TeamService from "../../services/team/teams";
import TeamCreate from "./QuickCreate/Team";
import Box from "@material-ui/core/Box";
import Stepper from "@material-ui/core/Stepper";
import Step from "@material-ui/core/Step";
import StepButton from "@material-ui/core/StepButton";
import Typography from "@material-ui/core/Typography";
import Paper from "@material-ui/core/Paper";
import HostCreate from "./QuickCreate/Host";
import UserService from "../../services/users/users";
import HostGroupsService from "../../services/host_group/host_groups";
import HostsService from "../../services/host/hosts";
import ServiceGroupsService from "../../services/service_group/service_groups";
import ServicesService from "../../services/service/serivces";


function getSteps() {
    return ['Regular View', 'Quick Create'];
}

function getStepContent(step, props) {
    switch (step) {
        case 0:
            const title = "Services"
            const isDependant = true
            const owningService = [ServiceGroupsService, HostsService]
            const owningFieldLookup = ["name", "address"]
            const fieldForLookup = ["service_group_id", "host_id"]
            const additionalActions = [{icon: "flash_on", tooltip: 'test service', onFuncClick: async (event, rowData) => {
                    props.handleLoading()
                    return await ServicesService.TestService(rowData["id"]).then((response) => {
                        let severity = "error"
                        let message = `Passed: ${response["Passed"]}, Log: ${response["Log"]}`
                        if (response["Passed"] && response["Err"] === ""){
                            severity = "success"
                        } else {
                            message += ` Error: ${response["Err"]}`
                        }
                        props.setAlert({message: message, severity:severity})
                    })
                } }]
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

            return Table(props, title, isDependant, columns, false, ServicesService, owningService, fieldForLookup, owningFieldLookup, true, additionalActions)
        case 1:
            return <HostCreate {...props} />;
        default:
            return 'Unknown step';
    }
}


export default function ServiceMenu(props) {
    const [activeStep, setActiveStep] = React.useState(0);
    const steps = getSteps();
    const handleStep = (step) => () => {
        setActiveStep(step);
    };
    return (
        <Box height="100%" width="100%" align="left" >
            <Stepper nonLinear activeStep={activeStep}>
                {steps.map((label, index) => (
                    <Step key={label}>
                        <StepButton onClick={handleStep(index)}>
                            {label}
                        </StepButton>
                    </Step>
                ))}
            </Stepper>
            <div>
                <div>
                    <div>
                        {getStepContent(activeStep, props)}
                    </div>
                </div>
            </div>
        </Box>
    );
}