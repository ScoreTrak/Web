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
import ServiceCreate from "./QuickCreate/Services";


function getSteps() {
    return ['Regular View', 'Quick Create'];
}

function getStepContent(step, props) {
    switch (step) {
        case 0:
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
                    { title: 'Enabled', field: 'enabled', type: 'boolean', initialEditValue: true},
                    { title: 'Edit Host(Allow users to change Addresses)', field: 'edit_host', type: 'boolean' },
                ]

            return Table(props, title, isDependant, columns, false, HostsService, owningService, fieldForLookup, owningFieldLookup)
        case 1:
            return <HostCreate {...props} />;
        default:
            return 'Unknown step';
    }
}


export default function HostMenu(props) {
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