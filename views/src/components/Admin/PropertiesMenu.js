import React from "react";
import {Table} from "./TableInterface";
import Box from "@material-ui/core/Box";
import Stepper from "@material-ui/core/Stepper";
import Step from "@material-ui/core/Step";
import StepButton from "@material-ui/core/StepButton";
import ServicesService from "../../services/service/serivces";
import PropertyService from "../../services/property/properties";
import PropertiesCreate from "./QuickCreate/Properties";


function getSteps() {
    return ['Regular View', 'Quick Create'];
}

function getStepContent(step, props) {
    switch (step) {
        case 0:
            const title = "Properties"
            const isDependant = true
            const owningService = [ServicesService]
            const owningFieldLookup = ["id"]
            const fieldForLookup = ["service_id"]
            const idFields = ["service_id", "key"]
            const columns=
                [
                    { title: 'Key', field: 'key', editable: 'onAdd'},
                    { title: 'Value', field: 'value' },
                    { title: 'Status', field: 'status', lookup:{'View': 'View', 'Hide':'Hide', 'Edit':'Edit'}},
                    { title: 'Service ID', field: 'service_id', editable: 'onAdd'},
                ]
            return Table(props, title, isDependant, columns, false, PropertyService, owningService, fieldForLookup, owningFieldLookup, true, [], idFields)
        case 1:
            return <PropertiesCreate {...props} />;
        default:
            return 'Unknown step';
    }
}


export default function PropertiesMenu(props) {
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