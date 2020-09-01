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


function getSteps() {
    return ['Regular View', 'Quick Create'];
}

function getStepContent(step, props) {
    switch (step) {
        case 0:
            const columns=
                [
                    { title: 'ID (optional)', field: 'id', editable: 'onAdd'},
                    { title: 'Team Name', field: 'name' },
                    { title: 'Index', field: 'index', type: 'numeric' },
                    { title: 'Enabled', field: 'enabled', type: 'boolean' },
                ]
            return Table(props, "Teams", false, columns, false, TeamService)
        case 1:
            return <TeamCreate {...props} />;
        default:
            return 'Unknown step';
    }
}


export default function TeamMenu(props) {
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